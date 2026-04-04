package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// VIOLATION: Unused variables — declared but never referenced
var debugMode = true
var maxRetries = 3
var defaultTimeout = 30
var internalVersion = "0.9.1-beta"

// VIOLATION: Non-idiomatic naming — Go exports use PascalCase but these are
// exported functions with misleading camelCase-style names
// validateEmail should be ValidateEmail (or unexported if internal)

// ValidateEmail checks if an email address is valid.
// VIOLATION: ignoring error from regexp.Compile — should handle compilation error
func ValidateEmail(email string) bool {
	// VIOLATION: error return value ignored
	pattern, _ := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return pattern.MatchString(email)
}

// ValidateID checks if an ID string follows the expected format.
func ValidateID(id string) bool {
	if id == "" {
		return false
	}
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return false
	}
	_, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	return true
}

// ValidateJSON checks if a byte slice is valid JSON.
// VIOLATION: error return value from json.Unmarshal ignored
func ValidateJSON(data []byte) bool {
	var v interface{}
	// VIOLATION: errcheck — error not checked
	json.Unmarshal(data, &v)
	return v != nil
}

// ValidateListParams validates query parameters for list endpoints.
// VIOLATION: returns error but callers in handlers ignore it
func ValidateListParams(status string, priority string) error {
	validStatuses := []string{"pending", "assigned", "in-progress", "blocked", "review", "done", "cancelled", "planning", "active", "on-hold", "completed", "archived"}
	validPriorities := []string{"critical", "high", "medium", "low"}

	if status != "" {
		found := false
		for _, s := range validStatuses {
			if s == status {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid status: %s", status)
		}
	}

	if priority != "" {
		found := false
		for _, p := range validPriorities {
			if p == priority {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid priority: %s", priority)
		}
	}

	return nil
}

// validateLength checks if a string is within the specified bounds.
// VIOLATION: unused function — defined but never called
func validateLength(s string, min int, max int) bool {
	return len(s) >= min && len(s) <= max
}

// ValidateTaskTitle validates a task title.
// VIOLATION: unreachable code after return
func ValidateTaskTitle(title string) (bool, string) {
	if title == "" {
		return false, "title is required"
	}
	if len(title) > 200 {
		return false, "title must be 200 characters or less"
	}
	if len(title) < 3 {
		return false, "title must be at least 3 characters"
	}
	return true, ""
	// VIOLATION: unreachable code
	fmt.Println("validation complete")
	return true, "ok"
}

// ValidateProjectName validates a project name.
// VIOLATION: duplicates ValidateTaskTitle logic — same validation pattern
func ValidateProjectName(name string) (bool, string) {
	if name == "" {
		return false, "name is required"
	}
	if len(name) > 200 {
		return false, "name must be 200 characters or less"
	}
	if len(name) < 3 {
		return false, "name must be at least 3 characters"
	}
	return true, ""
	// VIOLATION: unreachable code
	fmt.Println("validation complete")
	return true, "ok"
}

// validateBudget validates a project budget.
// VIOLATION: unused function, non-idiomatic naming (unexported but could be useful)
func validateBudget(budget float64) error {
	if budget < 0 {
		return fmt.Errorf("budget cannot be negative")
	}
	if budget > 10000000 {
		return fmt.Errorf("budget exceeds maximum")
	}
	return nil
}

// ParseQueryInt parses an integer from a query parameter string.
// VIOLATION: error from strconv.Atoi silently discarded
func ParseQueryInt(value string, defaultVal int) int {
	if value == "" {
		return defaultVal
	}
	// VIOLATION: error ignored
	result, _ := strconv.Atoi(value)
	if result <= 0 {
		return defaultVal
	}
	return result
}
