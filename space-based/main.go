package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// LocalCache is a simple thread-safe in-memory cache.
type LocalCache struct {
	data sync.Map
}

// Get returns the value for a given key.
func (c *LocalCache) Get(key string) (string, bool) {
	val, ok := c.data.Load(key)
	if ok {
		return val.(string), true
	}
	return "", false
}

// Set stores the key/value pair in the cache.
func (c *LocalCache) Set(key, value string) {
	c.data.Store(key, value)
}

// Delete removes a key from the cache.
func (c *LocalCache) Delete(key string) {
	c.data.Delete(key)
}

// CacheUpdate is the structure for Pub/Sub messages to synchronize caches.
type CacheUpdate struct {
	Action string `json:"action"`          // "set" or "delete"
	Key    string `json:"key"`             // e.g. "project:org:repo"
	Value  string `json:"value,omitempty"` // Only applicable for "set" action
}

// bootstrapCache loads existing project data from Redis into the local cache.
func bootstrapCache(rdb *redis.Client, cache *LocalCache) {
	projects, err := rdb.HGetAll(ctx, "projects").Result()
	if err != nil {
		log.Printf("Error bootstrapping local cache: %v", err)
		return
	}
	for key, value := range projects {
		cache.Set(key, value)
		log.Printf("Bootstrapped cache key: %s", key)
	}
}

// subscribeForUpdates listens on the "cache_updates" channel to keep the local cache in sync.
func subscribeForUpdates(rdb *redis.Client, cache *LocalCache) {
	pubsub := rdb.Subscribe(ctx, "cache_updates")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Pub/Sub error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		var update CacheUpdate
		if err := json.Unmarshal([]byte(msg.Payload), &update); err != nil {
			log.Printf("Error parsing update message: %v", err)
			continue
		}

		switch update.Action {
		case "set":
			cache.Set(update.Key, update.Value)
			log.Printf("Cache updated (set): %s", update.Key)
		case "delete":
			cache.Delete(update.Key)
			log.Printf("Cache updated (delete): %s", update.Key)
		default:
			log.Printf("Unknown cache update action: %s", update.Action)
		}
	}
}

// getProjectHandler implements the GET /project/{org}/{repo} endpoint.
// It uses a read‑through cache strategy:
//  1. Check local in-memory cache.
//  2. On miss, check Redis.
//  3. On cache miss in Redis, query the GitHub API, then save to both caches and publish an update.
func getProjectHandler(rdb *redis.Client, cache *LocalCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		org := chi.URLParam(r, "org")
		repo := chi.URLParam(r, "repo")
		if org == "" || repo == "" {
			http.Error(w, "Missing organization or repository", http.StatusBadRequest)
			return
		}

		// Create a key in the form "project:{org}:{repo}"
		key := fmt.Sprintf("project:%s:%s", org, repo)

		// 1. Check local in-memory cache.
		if value, ok := cache.Get(key); ok {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(value))
			return
		}

		// 2. Check Redis (central cache).
		result, err := rdb.HGet(ctx, "projects", key).Result()
		if err == nil && result != "" {
			// Update local cache before returning.
			cache.Set(key, result)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(result))
			return
		} else if err != nil && err != redis.Nil {
			log.Printf("Redis error: %v", err)
		}

		// 3. Cache miss: Query the GitHub API.
		githubURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", org, repo)
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Get(githubURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Error querying GitHub API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading GitHub response", http.StatusInternalServerError)
			return
		}

		// Validate that body contains proper JSON.
		var js json.RawMessage
		if err := json.Unmarshal(body, &js); err != nil {
			http.Error(w, "Invalid JSON from GitHub", http.StatusInternalServerError)
			return
		}

		// Save the project data in Redis (central cache) in the "projects" hash.
		if err := rdb.HSet(ctx, "projects", key, body).Err(); err != nil {
			log.Printf("Error saving project to Redis: %v", err)
		}
		// Update local cache.
		cache.Set(key, string(body))

		// Publish an update event so that other pods can update their local caches.
		update := CacheUpdate{
			Action: "set",
			Key:    key,
			Value:  string(body),
		}
		payload, err := json.Marshal(update)
		if err == nil {
			if err := rdb.Publish(ctx, "cache_updates", string(payload)).Err(); err != nil {
				log.Printf("Error publishing cache update: %v", err)
			}
		} else {
			log.Printf("Error marshalling cache update: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

func main() {
	// Initialize the Redis client (adjust the address if necessary).
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // Use default DB.
	})
	defer rdb.Close()

	// Initialize the local cache.
	localCache := &LocalCache{}

	// 1. Bootstrap the local cache from Redis.
	bootstrapCache(rdb, localCache)
	log.Println("Local cache bootstrapped from Redis.")

	// 2. Start a background goroutine to subscribe for cache updates.
	go subscribeForUpdates(rdb, localCache)

	// 3. Set up the HTTP router.
	r := chi.NewRouter()
	r.Get("/project/{org}/{repo}", getProjectHandler(rdb, localCache))

	// Start the web server.
	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
