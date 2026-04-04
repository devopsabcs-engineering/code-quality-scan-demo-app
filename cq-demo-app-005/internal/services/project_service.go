package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/models"
)

// ProjectService manages project operations with an in-memory store.
type ProjectService struct {
	mu       sync.RWMutex
	projects map[string]models.Project
	nextID   int
}

// NewProjectService creates a new ProjectService with seed data.
func NewProjectService() *ProjectService {
	svc := &ProjectService{
		projects: make(map[string]models.Project),
		nextID:   1,
	}
	svc.seedData()
	return svc
}

func (s *ProjectService) seedData() {
	now := time.Now()
	startDate := now.Add(-30 * 24 * time.Hour)
	endDate1 := now.Add(90 * 24 * time.Hour)
	endDate2 := now.Add(180 * 24 * time.Hour)

	seeds := []models.Project{
		{
			ID: "proj-001", Name: "Platform Modernization", Description: "Migrate legacy systems to cloud-native architecture",
			Status: models.ProjectStatusActive, Priority: models.ProjectPriorityCritical, OwnerID: "user-001",
			StartDate: &startDate, EndDate: &endDate1, CreatedAt: now, UpdatedAt: now,
			Budget: 500000, Tags: []string{"cloud", "migration"}, TeamSize: 12,
		},
		{
			ID: "proj-002", Name: "Developer Portal", Description: "Build internal developer portal with API docs",
			Status: models.ProjectStatusPlanning, Priority: models.ProjectPriorityHigh, OwnerID: "user-002",
			StartDate: &startDate, EndDate: &endDate2, CreatedAt: now, UpdatedAt: now,
			Budget: 200000, Tags: []string{"internal", "docs"}, TeamSize: 6,
		},
		{
			ID: "proj-003", Name: "Security Audit", Description: "Comprehensive security audit and remediation",
			Status: models.ProjectStatusOnHold, Priority: models.ProjectPriorityHigh, OwnerID: "user-003",
			CreatedAt: now, UpdatedAt: now,
			Budget: 150000, Tags: []string{"security", "compliance"}, TeamSize: 4,
		},
	}

	for _, p := range seeds {
		s.projects[p.ID] = p
	}
	s.nextID = 4
}

// ListProjects returns all projects.
func (s *ProjectService) ListProjects() []models.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Project, 0, len(s.projects))
	for _, p := range s.projects {
		result = append(result, p)
	}
	return result
}

// GetProject retrieves a project by ID.
func (s *ProjectService) GetProject(id string) (*models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, exists := s.projects[id]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", id)
	}
	return &project, nil
}

// CreateProject creates a new project.
func (s *ProjectService) CreateProject(req models.ProjectCreateRequest) *models.Project {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	project := models.Project{
		ID:          fmt.Sprintf("proj-%03d", s.nextID),
		Name:        req.Name,
		Description: req.Description,
		Status:      models.ProjectStatusPlanning,
		Priority:    req.Priority,
		OwnerID:     req.OwnerID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedAt:   now,
		UpdatedAt:   now,
		Budget:      req.Budget,
		Tags:        req.Tags,
		TeamSize:    req.TeamSize,
	}
	s.nextID++
	s.projects[project.ID] = project
	return &project
}

// UpdateProject updates a project by ID.
func (s *ProjectService) UpdateProject(id string, req models.ProjectUpdateRequest) (*models.Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	project, exists := s.projects[id]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", id)
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.Description != nil {
		project.Description = *req.Description
	}
	if req.Status != nil {
		project.Status = *req.Status
	}
	if req.Priority != nil {
		project.Priority = *req.Priority
	}
	if req.OwnerID != nil {
		project.OwnerID = *req.OwnerID
	}
	if req.StartDate != nil {
		project.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		project.EndDate = req.EndDate
	}
	if req.Budget != nil {
		project.Budget = *req.Budget
	}
	if req.Tags != nil {
		project.Tags = req.Tags
	}
	if req.TeamSize != nil {
		project.TeamSize = *req.TeamSize
	}
	project.UpdatedAt = time.Now()
	s.projects[id] = project
	return &project, nil
}

// DeleteProject removes a project by ID.
func (s *ProjectService) DeleteProject(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.projects[id]; !exists {
		return fmt.Errorf("project not found: %s", id)
	}
	delete(s.projects, id)
	return nil
}

// ProcessProjectStatus performs complex state-machine logic for project transitions.
// VIOLATION: CCN > 15 — duplicates the same pattern from task_service.go ProcessTask
// VIOLATION: Code duplication — nearly identical switch/case structure
func (s *ProjectService) ProcessProjectStatus(id string, action string, ownerID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	project, exists := s.projects[id]
	if !exists {
		return "", fmt.Errorf("project not found: %s", id)
	}

	var result string
	var newStatus models.ProjectStatus

	// VIOLATION: High cyclomatic complexity — duplicated switch/if pattern
	switch project.Status {
	case models.ProjectStatusPlanning:
		if action == "activate" {
			if project.Budget <= 0 {
				return "", fmt.Errorf("budget must be set before activation")
			}
			newStatus = models.ProjectStatusActive
			if project.Priority == models.ProjectPriorityCritical {
				if project.EndDate != nil && project.EndDate.Before(time.Now().Add(30*24*time.Hour)) {
					result = "URGENT: Critical project activated with tight deadline"
				} else {
					result = "Critical project activated"
				}
			} else if project.Priority == models.ProjectPriorityHigh {
				if project.EndDate != nil && project.EndDate.Before(time.Now().Add(60*24*time.Hour)) {
					result = "High priority project activated with approaching deadline"
				} else {
					result = "High priority project activated"
				}
			} else if project.Priority == models.ProjectPriorityMedium {
				result = "Medium priority project activated"
			} else {
				result = "Low priority project activated"
			}
		} else if action == "cancel" {
			newStatus = models.ProjectStatusCancelled
			result = "Planning project cancelled"
		} else {
			return "", fmt.Errorf("invalid action '%s' for planning project", action)
		}
	case models.ProjectStatusActive:
		if action == "hold" {
			newStatus = models.ProjectStatusOnHold
			if project.Priority == models.ProjectPriorityCritical {
				result = "CRITICAL: Project put on hold"
			} else if project.Priority == models.ProjectPriorityHigh {
				result = "High priority project on hold"
			} else {
				result = "Project put on hold"
			}
		} else if action == "complete" {
			newStatus = models.ProjectStatusCompleted
			result = "Project completed"
			if project.TeamSize > 10 {
				result = "Large project completed"
			}
		} else if action == "cancel" {
			newStatus = models.ProjectStatusCancelled
			result = "Active project cancelled"
		} else {
			return "", fmt.Errorf("invalid action '%s' for active project", action)
		}
	case models.ProjectStatusOnHold:
		if action == "resume" {
			newStatus = models.ProjectStatusActive
			result = "Project resumed"
		} else if action == "cancel" {
			newStatus = models.ProjectStatusCancelled
			result = "On-hold project cancelled"
		} else if action == "reassign" {
			if ownerID == "" {
				return "", fmt.Errorf("owner required for reassignment")
			}
			newStatus = models.ProjectStatusActive
			project.OwnerID = ownerID
			result = "On-hold project reassigned and resumed"
		} else {
			return "", fmt.Errorf("invalid action '%s' for on-hold project", action)
		}
	case models.ProjectStatusCompleted:
		if action == "reopen" {
			newStatus = models.ProjectStatusActive
			result = "Completed project reopened"
		} else if action == "archive" {
			newStatus = models.ProjectStatusArchived
			result = "Project archived"
		} else {
			return "", fmt.Errorf("completed projects can only be reopened or archived")
		}
	case models.ProjectStatusArchived:
		if action == "reopen" {
			newStatus = models.ProjectStatusPlanning
			result = "Archived project reopened to planning"
		} else {
			return "", fmt.Errorf("archived projects can only be reopened")
		}
	case models.ProjectStatusCancelled:
		if action == "reopen" {
			newStatus = models.ProjectStatusPlanning
			result = "Cancelled project reopened"
		} else {
			return "", fmt.Errorf("cancelled projects can only be reopened")
		}
	default:
		return "", fmt.Errorf("unknown project status: %s", project.Status)
	}

	project.Status = newStatus
	project.UpdatedAt = time.Now()
	s.projects[id] = project

	return result, nil
}

// CalculateProjectMetrics computes summary metrics for a project.
// VIOLATION: CCN > 10 — duplicated aggregation logic from task_service.go CalculateTaskMetrics
// VIOLATION: Code duplication — same metric calculation pattern
func (s *ProjectService) CalculateProjectMetrics(id string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, exists := s.projects[id]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", id)
	}

	daysActive := 0
	daysRemaining := 0
	budgetStatus := "on-track"
	riskLevel := "low"

	now := time.Now()

	if project.StartDate != nil {
		daysActive = int(now.Sub(*project.StartDate).Hours() / 24)
	}

	if project.EndDate != nil {
		daysRemaining = int(project.EndDate.Sub(now).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
			budgetStatus = "overdue"
		} else if daysRemaining < 14 {
			budgetStatus = "at-risk"
		}
	}

	if project.Priority == models.ProjectPriorityCritical {
		if project.Status == models.ProjectStatusOnHold || project.Status == models.ProjectStatusPlanning {
			riskLevel = "critical"
		} else if daysRemaining < 30 {
			riskLevel = "high"
		} else {
			riskLevel = "medium"
		}
	} else if project.Priority == models.ProjectPriorityHigh {
		if project.Status == models.ProjectStatusOnHold {
			riskLevel = "high"
		} else if daysRemaining < 14 {
			riskLevel = "high"
		} else {
			riskLevel = "medium"
		}
	} else {
		if daysRemaining < 7 {
			riskLevel = "medium"
		}
	}

	return map[string]interface{}{
		"project_id":     id,
		"project_name":   project.Name,
		"status":         string(project.Status),
		"priority":       string(project.Priority),
		"days_active":    daysActive,
		"days_remaining": daysRemaining,
		"budget":         project.Budget,
		"budget_status":  budgetStatus,
		"risk_level":     riskLevel,
		"team_size":      project.TeamSize,
	}, nil
}
