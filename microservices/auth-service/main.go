package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var db *sql.DB

func main() {
	// Connect to database
	var err error
	connStr := "postgres://postgres:password@localhost:5432/auth_service?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to auth database")

	// Set up routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth Service API"))
	})

	r.Post("/login", loginHandler)
	r.Post("/register", registerHandler)
	r.Get("/users/{id}", getUserHandler)

	log.Println("Auth Service started on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Mock implementation
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":  "mock-jwt-token",
		"status": "success",
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Mock implementation - in real code, parse request body and add user to DB
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"status":  "success",
	})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	// In a real implementation, fetch user from database
	// For this example, return mock data
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
