package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/models"
	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/services"
	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/utils"
)

// TaskHandler handles HTTP requests for task operations.
type TaskHandler struct {
	service *services.TaskService
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// HandleTasks handles GET /api/tasks and POST /api/tasks.
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleTaskByID handles GET/PUT/DELETE /api/tasks/{id}.
func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProcessTask handles POST /api/tasks/process — HIGH CYCLOMATIC COMPLEXITY
// VIOLATION: CCN > 15 — deeply nested switch/case with nested if/else
func (h *TaskHandler) HandleProcessTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TaskID     string `json:"task_id"`
		Action     string `json:"action"`
		AssigneeID string `json:"assignee_id"`
		Comment    string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(req.TaskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// VIOLATION: High cyclomatic complexity — deeply nested switch + if/else chains
	// This function has CCN > 15 and should be refactored into smaller functions
	var result string
	var newStatus models.TaskStatus
	var notifyAssignee bool
	var notifyManager bool
	var escalate bool

	switch task.Status {
	case models.TaskStatusPending:
		if req.Action == "assign" {
			if req.AssigneeID == "" {
				http.Error(w, "Assignee ID required for assignment", http.StatusBadRequest)
				return
			}
			if task.Priority == models.TaskPriorityCritical {
				notifyManager = true
				escalate = true
				if task.Deadline != nil && task.Deadline.Before(time.Now().Add(24*time.Hour)) {
					result = "URGENT: Critical task assigned with imminent deadline"
					notifyAssignee = true
				} else {
					result = "Critical task assigned"
					notifyAssignee = true
				}
			} else if task.Priority == models.TaskPriorityHigh {
				notifyAssignee = true
				if task.Deadline != nil && task.Deadline.Before(time.Now().Add(48*time.Hour)) {
					result = "High priority task assigned with approaching deadline"
					notifyManager = true
				} else {
					result = "High priority task assigned"
				}
			} else if task.Priority == models.TaskPriorityMedium {
				notifyAssignee = true
				result = "Medium priority task assigned"
			} else {
				notifyAssignee = true
				result = "Low priority task assigned"
			}
			newStatus = models.TaskStatusAssigned
		} else if req.Action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "Task cancelled from pending state"
			notifyManager = true
		} else {
			http.Error(w, "Invalid action for pending task: must be 'assign' or 'cancel'", http.StatusBadRequest)
			return
		}
	case models.TaskStatusAssigned:
		if req.Action == "start" {
			newStatus = models.TaskStatusInProgress
			if task.Priority == models.TaskPriorityCritical {
				result = "Critical task started — daily standup reporting enabled"
				notifyManager = true
			} else if task.Priority == models.TaskPriorityHigh {
				result = "High priority task started"
				notifyManager = true
			} else {
				result = "Task work started"
			}
		} else if req.Action == "reassign" {
			if req.AssigneeID == "" {
				http.Error(w, "Assignee ID required for reassignment", http.StatusBadRequest)
				return
			}
			newStatus = models.TaskStatusAssigned
			result = "Task reassigned"
			notifyAssignee = true
		} else if req.Action == "block" {
			newStatus = models.TaskStatusBlocked
			result = "Task blocked"
			notifyManager = true
			if task.Priority == models.TaskPriorityCritical || task.Priority == models.TaskPriorityHigh {
				escalate = true
			}
		} else {
			http.Error(w, "Invalid action for assigned task", http.StatusBadRequest)
			return
		}
	case models.TaskStatusInProgress:
		if req.Action == "review" {
			newStatus = models.TaskStatusReview
			result = "Task submitted for review"
			notifyManager = true
			if task.Estimate > 0 && task.Actual > task.Estimate*1.5 {
				result = result + " — WARNING: actual hours exceed estimate by 50%"
				escalate = true
			}
		} else if req.Action == "block" {
			newStatus = models.TaskStatusBlocked
			result = "Task blocked during progress"
			notifyManager = true
			if task.Priority == models.TaskPriorityCritical {
				escalate = true
				result = "CRITICAL: Task blocked — immediate escalation"
			}
		} else if req.Action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "Task cancelled during progress"
			notifyManager = true
			notifyAssignee = true
		} else {
			http.Error(w, "Invalid action for in-progress task", http.StatusBadRequest)
			return
		}
	case models.TaskStatusBlocked:
		if req.Action == "unblock" {
			newStatus = models.TaskStatusInProgress
			result = "Task unblocked — resuming work"
			notifyAssignee = true
			if task.Priority == models.TaskPriorityCritical || task.Priority == models.TaskPriorityHigh {
				notifyManager = true
			}
		} else if req.Action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "Blocked task cancelled"
			notifyManager = true
		} else if req.Action == "reassign" {
			if req.AssigneeID == "" {
				http.Error(w, "Assignee ID required", http.StatusBadRequest)
				return
			}
			newStatus = models.TaskStatusAssigned
			result = "Blocked task reassigned to different team member"
			notifyAssignee = true
			notifyManager = true
		} else {
			http.Error(w, "Invalid action for blocked task", http.StatusBadRequest)
			return
		}
	case models.TaskStatusReview:
		if req.Action == "approve" {
			newStatus = models.TaskStatusDone
			result = "Task approved and completed"
			notifyAssignee = true
			if task.Priority == models.TaskPriorityCritical {
				notifyManager = true
				result = "Critical task completed successfully"
			}
		} else if req.Action == "reject" {
			newStatus = models.TaskStatusInProgress
			result = "Task rejected — returned to assignee"
			notifyAssignee = true
			if req.Comment == "" {
				result = result + " — WARNING: no rejection reason provided"
			}
		} else if req.Action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "Task cancelled during review"
			notifyManager = true
			notifyAssignee = true
		} else {
			http.Error(w, "Invalid action for task in review", http.StatusBadRequest)
			return
		}
	case models.TaskStatusDone:
		if req.Action == "reopen" {
			newStatus = models.TaskStatusInProgress
			result = "Completed task reopened"
			notifyAssignee = true
			notifyManager = true
		} else {
			http.Error(w, "Completed tasks can only be reopened", http.StatusBadRequest)
			return
		}
	case models.TaskStatusCancelled:
		if req.Action == "reopen" {
			newStatus = models.TaskStatusPending
			result = "Cancelled task reopened"
			notifyManager = true
		} else {
			http.Error(w, "Cancelled tasks can only be reopened", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Unknown task status", http.StatusInternalServerError)
		return
	}

	// Apply status update
	statusVal := newStatus
	updateReq := models.TaskUpdateRequest{Status: &statusVal}
	if req.AssigneeID != "" {
		updateReq.AssigneeID = &req.AssigneeID
	}
	h.service.UpdateTask(req.TaskID, updateReq)

	// VIOLATION: unused variables — notifyAssignee, notifyManager, escalate are computed
	// but the notification system is not implemented (just logged for now)
	_ = notifyAssignee
	_ = notifyManager
	_ = escalate

	// Build response
	response := map[string]interface{}{
		"task_id":    req.TaskID,
		"old_status": string(task.Status),
		"new_status": string(newStatus),
		"result":     result,
		"action":     req.Action,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	tasks := h.service.ListTasks()

	// VIOLATION: duplicated filtering logic — same pattern exists in project_handler.go
	var filtered []models.Task
	for _, t := range tasks {
		if status != "" && string(t.Status) != status {
			continue
		}
		if priority != "" && string(t.Priority) != priority {
			continue
		}
		filtered = append(filtered, t)
	}

	// VIOLATION: ignoring error from validation
	utils.ValidateListParams(status, priority)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.service.GetTask(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Task not found: %s", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req models.TaskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// VIOLATION: duplicated validation — same pattern in project_handler.go
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if len(req.Title) > 200 {
		http.Error(w, "Title must be 200 characters or less", http.StatusBadRequest)
		return
	}
	if req.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	task := h.service.CreateTask(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var req models.TaskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.service.UpdateTask(id, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Task not found: %s", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteTask(id); err != nil {
		http.Error(w, fmt.Sprintf("Task not found: %s", id), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
