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

// ProjectHandler handles HTTP requests for project operations.
type ProjectHandler struct {
	service *services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// HandleProjects handles GET /api/projects and POST /api/projects.
func (h *ProjectHandler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listProjects(w, r)
	case http.MethodPost:
		h.createProject(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProjectByID handles GET/PUT/DELETE /api/projects/{id}.
func (h *ProjectHandler) HandleProjectByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if id == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getProject(w, r, id)
	case http.MethodPut:
		h.updateProject(w, r, id)
	case http.MethodDelete:
		h.deleteProject(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleUpdateProjectStatus handles POST /api/projects/status — HIGH CYCLOMATIC COMPLEXITY
// VIOLATION: CCN > 15 — duplicates the same nested switch/if pattern from task_handler.go
func (h *ProjectHandler) HandleUpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProjectID string `json:"project_id"`
		Action    string `json:"action"`
		OwnerID   string `json:"owner_id"`
		Comment   string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	project, err := h.service.GetProject(req.ProjectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// VIOLATION: High cyclomatic complexity — duplicates task_handler pattern
	// VIOLATION: Code duplication — same switch/if structure as HandleProcessTask
	var result string
	var newStatus models.ProjectStatus
	var notifyOwner bool
	var notifyStakeholders bool
	var escalate bool

	switch project.Status {
	case models.ProjectStatusPlanning:
		if req.Action == "activate" {
			if project.Budget <= 0 {
				http.Error(w, "Budget must be set before activation", http.StatusBadRequest)
				return
			}
			if project.Priority == models.ProjectPriorityCritical {
				notifyStakeholders = true
				escalate = true
				if project.EndDate != nil && project.EndDate.Before(time.Now().Add(30*24*time.Hour)) {
					result = "URGENT: Critical project activated with tight deadline"
					notifyOwner = true
				} else {
					result = "Critical project activated"
					notifyOwner = true
				}
			} else if project.Priority == models.ProjectPriorityHigh {
				notifyOwner = true
				if project.EndDate != nil && project.EndDate.Before(time.Now().Add(60*24*time.Hour)) {
					result = "High priority project activated with approaching deadline"
					notifyStakeholders = true
				} else {
					result = "High priority project activated"
				}
			} else if project.Priority == models.ProjectPriorityMedium {
				notifyOwner = true
				result = "Medium priority project activated"
			} else {
				notifyOwner = true
				result = "Low priority project activated"
			}
			newStatus = models.ProjectStatusActive
		} else if req.Action == "cancel" {
			newStatus = models.ProjectStatusCancelled
			result = "Project cancelled during planning"
			notifyStakeholders = true
		} else {
			http.Error(w, "Invalid action for planning project: must be 'activate' or 'cancel'", http.StatusBadRequest)
			return
		}
	case models.ProjectStatusActive:
		if req.Action == "hold" {
			newStatus = models.ProjectStatusOnHold
			if project.Priority == models.ProjectPriorityCritical {
				result = "CRITICAL: Project put on hold — immediate escalation"
				notifyStakeholders = true
				escalate = true
			} else if project.Priority == models.ProjectPriorityHigh {
				result = "High priority project on hold"
				notifyStakeholders = true
			} else {
				result = "Project put on hold"
			}
			notifyOwner = true
		} else if req.Action == "complete" {
			newStatus = models.ProjectStatusCompleted
			result = "Project completed"
			notifyOwner = true
			notifyStakeholders = true
			if project.TeamSize > 10 {
				result = "Large project completed — retrospective recommended"
			}
		} else if req.Action == "cancel" {
			newStatus = models.ProjectStatusCancelled
			result = "Active project cancelled"
			notifyOwner = true
			notifyStakeholders = true
			escalate = true
		} else {
			http.Error(w, "Invalid action for active project", http.StatusBadRequest)
			return
		}
	case models.ProjectStatusOnHold:
		if req.Action == "resume" {
			newStatus = models.ProjectStatusActive
			result = "Project resumed"
			notifyOwner = true
			if project.Priority == models.ProjectPriorityCritical || project.Priority == models.ProjectPriorityHigh {
				notifyStakeholders = true
			}
		} else if req.Action == "cancel" {
			newStatus = models.ProjectStatusCancelled
			result = "On-hold project cancelled"
			notifyStakeholders = true
		} else if req.Action == "reassign" {
			if req.OwnerID == "" {
				http.Error(w, "Owner ID required", http.StatusBadRequest)
				return
			}
			newStatus = models.ProjectStatusActive
			result = "On-hold project reassigned and resumed"
			notifyOwner = true
			notifyStakeholders = true
		} else {
			http.Error(w, "Invalid action for on-hold project", http.StatusBadRequest)
			return
		}
	case models.ProjectStatusCompleted:
		if req.Action == "reopen" {
			newStatus = models.ProjectStatusActive
			result = "Completed project reopened"
			notifyOwner = true
			notifyStakeholders = true
		} else if req.Action == "archive" {
			newStatus = models.ProjectStatusArchived
			result = "Completed project archived"
		} else {
			http.Error(w, "Completed projects can only be reopened or archived", http.StatusBadRequest)
			return
		}
	case models.ProjectStatusArchived:
		if req.Action == "reopen" {
			newStatus = models.ProjectStatusPlanning
			result = "Archived project reopened to planning"
			notifyOwner = true
			notifyStakeholders = true
		} else {
			http.Error(w, "Archived projects can only be reopened", http.StatusBadRequest)
			return
		}
	case models.ProjectStatusCancelled:
		if req.Action == "reopen" {
			newStatus = models.ProjectStatusPlanning
			result = "Cancelled project reopened"
			notifyStakeholders = true
		} else {
			http.Error(w, "Cancelled projects can only be reopened", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Unknown project status", http.StatusInternalServerError)
		return
	}

	// Apply status update
	statusVal := newStatus
	updateReq := models.ProjectUpdateRequest{Status: &statusVal}
	if req.OwnerID != "" {
		updateReq.OwnerID = &req.OwnerID
	}
	h.service.UpdateProject(req.ProjectID, updateReq)

	// VIOLATION: unused variables — computed but notification not implemented
	_ = notifyOwner
	_ = notifyStakeholders
	_ = escalate

	// Build response
	response := map[string]interface{}{
		"project_id": req.ProjectID,
		"old_status": string(project.Status),
		"new_status": string(newStatus),
		"result":     result,
		"action":     req.Action,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VIOLATION: Duplicated filtering logic — same pattern as task_handler.go listTasks
func (h *ProjectHandler) listProjects(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	projects := h.service.ListProjects()

	// VIOLATION: duplicated filtering logic — copied from task_handler.go
	var filtered []models.Project
	for _, p := range projects {
		if status != "" && string(p.Status) != status {
			continue
		}
		if priority != "" && string(p.Priority) != priority {
			continue
		}
		filtered = append(filtered, p)
	}

	// VIOLATION: ignoring error from validation — same as task_handler.go
	utils.ValidateListParams(status, priority)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func (h *ProjectHandler) getProject(w http.ResponseWriter, r *http.Request, id string) {
	project, err := h.service.GetProject(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Project not found: %s", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// VIOLATION: Duplicated validation — same title/description check as task_handler.go createTask
func (h *ProjectHandler) createProject(w http.ResponseWriter, r *http.Request) {
	var req models.ProjectCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// VIOLATION: duplicated validation — same pattern as task_handler.go createTask
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if len(req.Name) > 200 {
		http.Error(w, "Name must be 200 characters or less", http.StatusBadRequest)
		return
	}
	if req.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	project := h.service.CreateProject(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) updateProject(w http.ResponseWriter, r *http.Request, id string) {
	var req models.ProjectUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	project, err := h.service.UpdateProject(id, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Project not found: %s", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) deleteProject(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteProject(id); err != nil {
		http.Error(w, fmt.Sprintf("Project not found: %s", id), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
