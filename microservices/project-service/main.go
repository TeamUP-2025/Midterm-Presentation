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

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     int    `json:"owner_id"`
}

var db *sql.DB

func main() {
	// Connect to database
	var err error
	connStr := "postgres://postgres:password@localhost:5432/project_service?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to project database")

	// Set up routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Project Service API"))
	})

	r.Get("/projects", getProjectsHandler)
	r.Post("/projects", createProjectHandler)
	r.Get("/projects/{id}", getProjectHandler)

	log.Println("Project Service started on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}

func getProjectsHandler(w http.ResponseWriter, r *http.Request) {
	// Mock implementation - would query DB in real implementation
	projects := []Project{
		{ID: 1, Name: "First Project", Description: "This is the first project", OwnerID: 1},
		{ID: 2, Name: "Second Project", Description: "This is the second project", OwnerID: 2},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func createProjectHandler(w http.ResponseWriter, r *http.Request) {
	// Mock implementation - would parse body and insert to DB
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Project created successfully",
		"status":  "success",
	})
}

func getProjectHandler(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	// Mock project data - would query DB in real implementation
	project := Project{
		ID:          1,
		Name:        "Example Project",
		Description: "This is an example project",
		OwnerID:     1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
