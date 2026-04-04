package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/devopsabcs-engineering/cq-demo-app-005/internal/models"
)

// TaskService manages task operations with an in-memory store.
type TaskService struct {
	mu    sync.RWMutex
	tasks map[string]models.Task
	nextID int
}

// NewTaskService creates a new TaskService with seed data.
func NewTaskService() *TaskService {
	svc := &TaskService{
		tasks:  make(map[string]models.Task),
		nextID: 1,
	}
	svc.seedData()
	return svc
}

func (s *TaskService) seedData() {
	now := time.Now()
	deadline1 := now.Add(72 * time.Hour)
	deadline2 := now.Add(168 * time.Hour)
	deadline3 := now.Add(24 * time.Hour)

	seeds := []models.Task{
		{
			ID: "task-001", Title: "Set up CI/CD pipeline", Description: "Configure GitHub Actions for automated builds",
			Status: models.TaskStatusInProgress, Priority: models.TaskPriorityHigh, AssigneeID: "user-001",
			ProjectID: "proj-001", Deadline: &deadline1, CreatedAt: now, UpdatedAt: now,
			Tags: []string{"devops", "ci-cd"}, Estimate: 8, Actual: 3,
		},
		{
			ID: "task-002", Title: "Write unit tests for auth module", Description: "Add tests for JWT validation",
			Status: models.TaskStatusPending, Priority: models.TaskPriorityCritical, AssigneeID: "",
			ProjectID: "proj-001", Deadline: &deadline3, CreatedAt: now, UpdatedAt: now,
			Tags: []string{"testing", "security"}, Estimate: 16, Actual: 0,
		},
		{
			ID: "task-003", Title: "Database schema migration", Description: "Migrate from v2 to v3 schema",
			Status: models.TaskStatusReview, Priority: models.TaskPriorityMedium, AssigneeID: "user-002",
			ProjectID: "proj-002", Deadline: &deadline2, CreatedAt: now, UpdatedAt: now,
			Tags: []string{"database", "migration"}, Estimate: 12, Actual: 10,
		},
		{
			ID: "task-004", Title: "Update README documentation", Description: "Add API examples and setup guide",
			Status: models.TaskStatusDone, Priority: models.TaskPriorityLow, AssigneeID: "user-003",
			ProjectID: "proj-001", CreatedAt: now, UpdatedAt: now,
			Tags: []string{"docs"}, Estimate: 4, Actual: 5,
		},
		{
			ID: "task-005", Title: "Performance profiling", Description: "Profile and optimize hot paths",
			Status: models.TaskStatusBlocked, Priority: models.TaskPriorityHigh, AssigneeID: "user-001",
			ProjectID: "proj-002", Deadline: &deadline2, CreatedAt: now, UpdatedAt: now,
			Tags: []string{"performance"}, Estimate: 20, Actual: 6,
		},
	}

	for _, t := range seeds {
		s.tasks[t.ID] = t
	}
	s.nextID = 6
}

// ListTasks returns all tasks.
func (s *TaskService) ListTasks() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

// GetTask retrieves a task by ID.
func (s *TaskService) GetTask(id string) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	return &task, nil
}

// CreateTask creates a new task.
func (s *TaskService) CreateTask(req models.TaskCreateRequest) *models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	task := models.Task{
		ID:          fmt.Sprintf("task-%03d", s.nextID),
		Title:       req.Title,
		Description: req.Description,
		Status:      models.TaskStatusPending,
		Priority:    req.Priority,
		AssigneeID:  req.AssigneeID,
		ProjectID:   req.ProjectID,
		Deadline:    req.Deadline,
		CreatedAt:   now,
		UpdatedAt:   now,
		Tags:        req.Tags,
		Estimate:    req.Estimate,
	}
	s.nextID++
	s.tasks[task.ID] = task
	return &task
}

// UpdateTask updates a task by ID.
func (s *TaskService) UpdateTask(id string, req models.TaskUpdateRequest) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", id)
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		task.AssigneeID = *req.AssigneeID
	}
	if req.Deadline != nil {
		task.Deadline = req.Deadline
	}
	if req.Tags != nil {
		task.Tags = req.Tags
	}
	if req.Estimate != nil {
		task.Estimate = *req.Estimate
	}
	if req.Actual != nil {
		task.Actual = *req.Actual
	}
	task.UpdatedAt = time.Now()
	s.tasks[id] = task
	return &task, nil
}

// DeleteTask removes a task by ID.
func (s *TaskService) DeleteTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("task not found: %s", id)
	}
	delete(s.tasks, id)
	return nil
}

// ProcessTask performs complex state-machine logic for task transitions.
// VIOLATION: CCN > 15 — deeply nested switch/case with if/else for priority, deadline, assignment
func (s *TaskService) ProcessTask(id string, action string, assigneeID string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return "", fmt.Errorf("task not found: %s", id)
	}

	var result string
	var newStatus models.TaskStatus

	// VIOLATION: High cyclomatic complexity — deeply nested switch/case with if/else chains
	// This duplicates logic from task_handler.go HandleProcessTask
	switch task.Status {
	case models.TaskStatusPending:
		if action == "assign" {
			if assigneeID == "" {
				return "", fmt.Errorf("assignee required")
			}
			newStatus = models.TaskStatusAssigned
			task.AssigneeID = assigneeID
			if task.Priority == models.TaskPriorityCritical {
				if task.Deadline != nil && task.Deadline.Before(time.Now().Add(24*time.Hour)) {
					result = "URGENT: Critical task with imminent deadline assigned"
				} else {
					result = "Critical task assigned"
				}
			} else if task.Priority == models.TaskPriorityHigh {
				if task.Deadline != nil && task.Deadline.Before(time.Now().Add(48*time.Hour)) {
					result = "High priority task with approaching deadline assigned"
				} else {
					result = "High priority task assigned"
				}
			} else if task.Priority == models.TaskPriorityMedium {
				result = "Medium priority task assigned"
			} else {
				result = "Low priority task assigned"
			}
		} else if action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "Pending task cancelled"
		} else {
			return "", fmt.Errorf("invalid action '%s' for pending task", action)
		}
	case models.TaskStatusAssigned:
		if action == "start" {
			newStatus = models.TaskStatusInProgress
			if task.Priority == models.TaskPriorityCritical {
				result = "Critical task started"
			} else if task.Priority == models.TaskPriorityHigh {
				result = "High priority task started"
			} else {
				result = "Task started"
			}
		} else if action == "reassign" {
			if assigneeID == "" {
				return "", fmt.Errorf("assignee required for reassignment")
			}
			newStatus = models.TaskStatusAssigned
			task.AssigneeID = assigneeID
			result = "Task reassigned"
		} else if action == "block" {
			newStatus = models.TaskStatusBlocked
			result = "Task blocked"
		} else {
			return "", fmt.Errorf("invalid action '%s' for assigned task", action)
		}
	case models.TaskStatusInProgress:
		if action == "review" {
			newStatus = models.TaskStatusReview
			result = "Task submitted for review"
			if task.Estimate > 0 && task.Actual > task.Estimate*1.5 {
				result = result + " — over budget"
			}
		} else if action == "block" {
			newStatus = models.TaskStatusBlocked
			result = "Task blocked"
			if task.Priority == models.TaskPriorityCritical {
				result = "CRITICAL: Task blocked"
			}
		} else if action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "In-progress task cancelled"
		} else {
			return "", fmt.Errorf("invalid action '%s' for in-progress task", action)
		}
	case models.TaskStatusBlocked:
		if action == "unblock" {
			newStatus = models.TaskStatusInProgress
			result = "Task unblocked"
		} else if action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "Blocked task cancelled"
		} else if action == "reassign" {
			if assigneeID == "" {
				return "", fmt.Errorf("assignee required")
			}
			newStatus = models.TaskStatusAssigned
			task.AssigneeID = assigneeID
			result = "Blocked task reassigned"
		} else {
			return "", fmt.Errorf("invalid action '%s' for blocked task", action)
		}
	case models.TaskStatusReview:
		if action == "approve" {
			newStatus = models.TaskStatusDone
			result = "Task approved"
			if task.Priority == models.TaskPriorityCritical {
				result = "Critical task completed"
			}
		} else if action == "reject" {
			newStatus = models.TaskStatusInProgress
			result = "Task rejected — returned to assignee"
		} else if action == "cancel" {
			newStatus = models.TaskStatusCancelled
			result = "Task cancelled during review"
		} else {
			return "", fmt.Errorf("invalid action '%s' for task in review", action)
		}
	case models.TaskStatusDone:
		if action == "reopen" {
			newStatus = models.TaskStatusInProgress
			result = "Completed task reopened"
		} else {
			return "", fmt.Errorf("completed tasks can only be reopened")
		}
	case models.TaskStatusCancelled:
		if action == "reopen" {
			newStatus = models.TaskStatusPending
			result = "Cancelled task reopened"
		} else {
			return "", fmt.Errorf("cancelled tasks can only be reopened")
		}
	default:
		return "", fmt.Errorf("unknown task status: %s", task.Status)
	}

	task.Status = newStatus
	task.UpdatedAt = time.Now()
	s.tasks[id] = task

	return result, nil
}

// CalculateTaskMetrics computes summary metrics for all tasks in a project.
// VIOLATION: CCN > 10 — complex aggregation with nested conditionals
func (s *TaskService) CalculateTaskMetrics(projectID string) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalTasks := 0
	completedTasks := 0
	blockedTasks := 0
	overdueTasks := 0
	totalEstimate := 0.0
	totalActual := 0.0
	criticalOpen := 0
	highOpen := 0

	now := time.Now()

	for _, task := range s.tasks {
		if task.ProjectID != projectID {
			continue
		}
		totalTasks++

		if task.Status == models.TaskStatusDone {
			completedTasks++
		} else if task.Status == models.TaskStatusBlocked {
			blockedTasks++
			if task.Priority == models.TaskPriorityCritical {
				criticalOpen++
			} else if task.Priority == models.TaskPriorityHigh {
				highOpen++
			}
		} else if task.Status != models.TaskStatusCancelled {
			if task.Priority == models.TaskPriorityCritical {
				criticalOpen++
			} else if task.Priority == models.TaskPriorityHigh {
				highOpen++
			}
		}

		if task.Deadline != nil && task.Deadline.Before(now) && task.Status != models.TaskStatusDone && task.Status != models.TaskStatusCancelled {
			overdueTasks++
		}

		totalEstimate += task.Estimate
		totalActual += task.Actual
	}

	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	budgetVariance := 0.0
	if totalEstimate > 0 {
		budgetVariance = ((totalActual - totalEstimate) / totalEstimate) * 100
	}

	return map[string]interface{}{
		"total_tasks":      totalTasks,
		"completed_tasks":  completedTasks,
		"blocked_tasks":    blockedTasks,
		"overdue_tasks":    overdueTasks,
		"critical_open":    criticalOpen,
		"high_open":        highOpen,
		"completion_rate":  completionRate,
		"total_estimate":   totalEstimate,
		"total_actual":     totalActual,
		"budget_variance":  budgetVariance,
	}
}
