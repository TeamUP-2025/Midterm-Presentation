package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes to services
	r.Mount("/auth", createProxy("http://localhost:8081"))
	r.Mount("/projects", createProxy("http://localhost:8082"))
	r.Mount("/chat", createProxy("http://localhost:8083"))
	r.Mount("/donations", createProxy("http://localhost:8084"))

	log.Println("API Gateway started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createProxy(target string) http.Handler {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return proxy
}
