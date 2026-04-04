package models

import "time"

// TaskStatus represents the lifecycle state of a task.
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusAssigned   TaskStatus = "assigned"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusBlocked    TaskStatus = "blocked"
	TaskStatusReview     TaskStatus = "review"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

// TaskPriority represents the urgency of a task.
type TaskPriority string

const (
	TaskPriorityCritical TaskPriority = "critical"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityLow      TaskPriority = "low"
)

// Task represents a work item in the project management system.
type Task struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	AssigneeID  string       `json:"assignee_id"`
	ProjectID   string       `json:"project_id"`
	Deadline    *time.Time   `json:"deadline,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Tags        []string     `json:"tags,omitempty"`
	Estimate    float64      `json:"estimate_hours"`
	Actual      float64      `json:"actual_hours"`
}

// TaskCreateRequest is the payload for creating a new task.
type TaskCreateRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	AssigneeID  string       `json:"assignee_id"`
	ProjectID   string       `json:"project_id"`
	Deadline    *time.Time   `json:"deadline,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Estimate    float64      `json:"estimate_hours"`
}

// TaskUpdateRequest is the payload for updating a task.
type TaskUpdateRequest struct {
	Title       *string       `json:"title,omitempty"`
	Description *string       `json:"description,omitempty"`
	Status      *TaskStatus   `json:"status,omitempty"`
	Priority    *TaskPriority `json:"priority,omitempty"`
	AssigneeID  *string       `json:"assignee_id,omitempty"`
	Deadline    *time.Time    `json:"deadline,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	Estimate    *float64      `json:"estimate_hours,omitempty"`
	Actual      *float64      `json:"actual_hours,omitempty"`
}
