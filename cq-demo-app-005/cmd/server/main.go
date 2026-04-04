package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/handlers"
	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	taskService := services.NewTaskService()
	projectService := services.NewProjectService()

	taskHandler := handlers.NewTaskHandler(taskService)
	projectHandler := handlers.NewProjectHandler(projectService)

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "cq-demo-app-005",
			"version": "1.0.0",
		})
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Task routes
	mux.HandleFunc("/api/tasks", taskHandler.HandleTasks)
	mux.HandleFunc("/api/tasks/", taskHandler.HandleTaskByID)
	mux.HandleFunc("/api/tasks/process", taskHandler.HandleProcessTask)

	// Project routes
	mux.HandleFunc("/api/projects", projectHandler.HandleProjects)
	mux.HandleFunc("/api/projects/", projectHandler.HandleProjectByID)
	mux.HandleFunc("/api/projects/status", projectHandler.HandleUpdateProjectStatus)

	log.Printf("Starting cq-demo-app-005 on port %s", port)
	addr := fmt.Sprintf(":%s", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
