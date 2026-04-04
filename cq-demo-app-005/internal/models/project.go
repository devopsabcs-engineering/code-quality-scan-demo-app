package models

import "time"

// ProjectStatus represents the lifecycle state of a project.
type ProjectStatus string

const (
	ProjectStatusPlanning   ProjectStatus = "planning"
	ProjectStatusActive     ProjectStatus = "active"
	ProjectStatusOnHold     ProjectStatus = "on-hold"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusArchived   ProjectStatus = "archived"
	ProjectStatusCancelled  ProjectStatus = "cancelled"
)

// ProjectPriority represents the priority of a project.
type ProjectPriority string

const (
	ProjectPriorityCritical ProjectPriority = "critical"
	ProjectPriorityHigh     ProjectPriority = "high"
	ProjectPriorityMedium   ProjectPriority = "medium"
	ProjectPriorityLow      ProjectPriority = "low"
)

// Project represents a project in the management system.
type Project struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      ProjectStatus   `json:"status"`
	Priority    ProjectPriority `json:"priority"`
	OwnerID     string          `json:"owner_id"`
	StartDate   *time.Time      `json:"start_date,omitempty"`
	EndDate     *time.Time      `json:"end_date,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Budget      float64         `json:"budget"`
	Tags        []string        `json:"tags,omitempty"`
	TeamSize    int             `json:"team_size"`
}

// ProjectCreateRequest is the payload for creating a new project.
type ProjectCreateRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Priority    ProjectPriority `json:"priority"`
	OwnerID     string          `json:"owner_id"`
	StartDate   *time.Time      `json:"start_date,omitempty"`
	EndDate     *time.Time      `json:"end_date,omitempty"`
	Budget      float64         `json:"budget"`
	Tags        []string        `json:"tags,omitempty"`
	TeamSize    int             `json:"team_size"`
}

// ProjectUpdateRequest is the payload for updating a project.
type ProjectUpdateRequest struct {
	Name        *string          `json:"name,omitempty"`
	Description *string          `json:"description,omitempty"`
	Status      *ProjectStatus   `json:"status,omitempty"`
	Priority    *ProjectPriority `json:"priority,omitempty"`
	OwnerID     *string          `json:"owner_id,omitempty"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	Budget      *float64         `json:"budget,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	TeamSize    *int             `json:"team_size,omitempty"`
}
