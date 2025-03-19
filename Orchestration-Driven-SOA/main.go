// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Start all services
	go startBusinessService()
	go startEnterpriseServiceBus()
	go startProjectService()
	go startGitHubDataService()
	go startInfrastructureService()

	// Wait indefinitely
	select {}
}

// Business Service
func startBusinessService() {
	router := mux.NewRouter()
	router.HandleFunc("/project", createProject).Methods("POST")

	fmt.Println("Business Service started on :8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

func createProject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Business Service: CreateProject request received")

	// Call Orchestration Engine to handle project creation workflow
	resp, err := http.Post("http://localhost:8002/orchestrate/project",
		"application/json", r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	fmt.Fprintf(w, "Project creation orchestrated successfully")
}

// Enterprise Service Bus
func startEnterpriseServiceBus() {
	router := mux.NewRouter()

	// Orchestration Engine endpoints
	router.HandleFunc("/orchestrate/project", orchestrateProjectCreation).Methods("POST")

	// Integration Hub endpoints
	router.HandleFunc("/integrate/services", integrateServices).Methods("POST")

	fmt.Println("Enterprise Service Bus started on :8002")
	log.Fatal(http.ListenAndServe(":8002", router))
}

func orchestrateProjectCreation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Orchestration Engine: Orchestrating project creation")

	// Log the operation
	logOperation("Orchestrating project creation workflow")

	// Step 1: Initialize Project (calling separate Project Service)
	resp, err := http.Post("http://localhost:8003/project/initialize",
		"application/json", r.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to initialize project", http.StatusInternalServerError)
		return
	}
	resp.Body.Close()

	// Step 2: Acquire Data from GitHub (calling separate GitHub Data Service)
	resp, err = http.Post("http://localhost:8004/github/data",
		"application/json", r.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to acquire GitHub data", http.StatusInternalServerError)
		return
	}
	resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Project creation workflow completed")
}

func integrateServices(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Integration Hub: Integrating services")
	w.WriteHeader(http.StatusOK)
}

// Project Service
func startProjectService() {
	router := mux.NewRouter()
	router.HandleFunc("/project/initialize", initializeProject).Methods("POST")

	fmt.Println("Project Service started on :8003")
	log.Fatal(http.ListenAndServe(":8003", router))
}

func initializeProject(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Project Service: Initializing Project")
	logOperation("Project initialization")

	// Simulate project initialization
	time.Sleep(1 * time.Second)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Project initialized")
}

// GitHub Data Service
func startGitHubDataService() {
	router := mux.NewRouter()
	router.HandleFunc("/github/data", acquireDataFromGitHub).Methods("POST")

	fmt.Println("GitHub Data Service started on :8004")
	log.Fatal(http.ListenAndServe(":8004", router))
}

func acquireDataFromGitHub(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GitHub Data Service: Acquiring Data from GitHub")
	logOperation("GitHub data acquisition")

	// Simulate fetching data from GitHub
	time.Sleep(2 * time.Second)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "GitHub data acquired")
}

// Infrastructure Service
func startInfrastructureService() {
	router := mux.NewRouter()
	router.HandleFunc("/log", logEndpoint).Methods("POST")

	fmt.Println("Infrastructure Service started on :8005")
	log.Fatal(http.ListenAndServe(":8005", router))
}

func logEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Infrastructure Service: Logging")
	w.WriteHeader(http.StatusOK)
}

func logOperation(message string) {
	// Call logging service
	_, err := http.Post("http://localhost:8005/log",
		"text/plain",
		nil)

	if err != nil {
		fmt.Printf("Error logging: %v\n", err)
	}

	// For demo purposes, also print to console
	fmt.Printf("Log: %s at %v\n", message, time.Now())
}
