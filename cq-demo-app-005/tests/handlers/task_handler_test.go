package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/handlers"
	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/models"
	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/services"
)

// VIOLATION: Low test coverage — only GET /api/tasks is tested.
// Missing tests for:
//   - POST /api/tasks (createTask)
//   - PUT /api/tasks/{id} (updateTask)
//   - DELETE /api/tasks/{id} (deleteTask)
//   - POST /api/tasks/process (HandleProcessTask — the most complex function)
//   - All project handler endpoints
//   - All service layer functions
//   - All utility functions (validators, formatters)

func TestListTasks(t *testing.T) {
	taskService := services.NewTaskService()
	handler := handlers.NewTaskHandler(taskService)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rr := httptest.NewRecorder()

	handler.HandleTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var tasks []models.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(tasks) == 0 {
		t.Error("expected at least one task in seed data")
	}
}

func TestListTasksWithStatusFilter(t *testing.T) {
	taskService := services.NewTaskService()
	handler := handlers.NewTaskHandler(taskService)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks?status=pending", nil)
	rr := httptest.NewRecorder()

	handler.HandleTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var tasks []models.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	for _, task := range tasks {
		if task.Status != models.TaskStatusPending {
			t.Errorf("expected status 'pending', got '%s' for task %s", task.Status, task.ID)
		}
	}
}

func TestGetTaskByID(t *testing.T) {
	taskService := services.NewTaskService()
	handler := handlers.NewTaskHandler(taskService)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/task-001", nil)
	rr := httptest.NewRecorder()

	handler.HandleTaskByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var task models.Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if task.ID != "task-001" {
		t.Errorf("expected task ID 'task-001', got '%s'", task.ID)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	taskService := services.NewTaskService()
	handler := handlers.NewTaskHandler(taskService)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.HandleTaskByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	taskService := services.NewTaskService()
	handler := handlers.NewTaskHandler(taskService)

	req := httptest.NewRequest(http.MethodPatch, "/api/tasks", nil)
	rr := httptest.NewRecorder()

	handler.HandleTasks(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
}
