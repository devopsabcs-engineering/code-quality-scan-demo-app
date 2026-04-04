package utils

import (
	"fmt"
	"strings"
	"time"
)

// VIOLATION: Magic strings — status labels hardcoded throughout
// VIOLATION: Missing error wrapping — errors returned without context via fmt.Errorf

// FormatTaskSummary formats a task summary string.
// VIOLATION: Magic strings used throughout instead of constants
func FormatTaskSummary(id string, title string, status string, priority string) string {
	// VIOLATION: magic strings for status labels
	var statusLabel string
	if status == "pending" {
		statusLabel = "⏳ Pending"
	} else if status == "assigned" {
		statusLabel = "👤 Assigned"
	} else if status == "in-progress" {
		statusLabel = "🔄 In Progress"
	} else if status == "blocked" {
		statusLabel = "🚫 Blocked"
	} else if status == "review" {
		statusLabel = "👀 In Review"
	} else if status == "done" {
		statusLabel = "✅ Done"
	} else if status == "cancelled" {
		statusLabel = "❌ Cancelled"
	} else {
		statusLabel = "❓ Unknown"
	}

	// VIOLATION: magic strings for priority labels
	var priorityLabel string
	if priority == "critical" {
		priorityLabel = "🔴 CRITICAL"
	} else if priority == "high" {
		priorityLabel = "🟠 HIGH"
	} else if priority == "medium" {
		priorityLabel = "🟡 MEDIUM"
	} else if priority == "low" {
		priorityLabel = "🟢 LOW"
	} else {
		priorityLabel = "⚪ UNKNOWN"
	}

	return fmt.Sprintf("[%s] %s | %s | %s", id, title, statusLabel, priorityLabel)
}

// FormatProjectSummary formats a project summary string.
// VIOLATION: Code duplication — same status/priority label logic as FormatTaskSummary
func FormatProjectSummary(id string, name string, status string, priority string) string {
	// VIOLATION: duplicated magic string labels from FormatTaskSummary
	var statusLabel string
	if status == "planning" {
		statusLabel = "📋 Planning"
	} else if status == "active" {
		statusLabel = "🚀 Active"
	} else if status == "on-hold" {
		statusLabel = "⏸ On Hold"
	} else if status == "completed" {
		statusLabel = "✅ Completed"
	} else if status == "archived" {
		statusLabel = "📦 Archived"
	} else if status == "cancelled" {
		statusLabel = "❌ Cancelled"
	} else {
		statusLabel = "❓ Unknown"
	}

	// VIOLATION: duplicated priority formatting — same as FormatTaskSummary
	var priorityLabel string
	if priority == "critical" {
		priorityLabel = "🔴 CRITICAL"
	} else if priority == "high" {
		priorityLabel = "🟠 HIGH"
	} else if priority == "medium" {
		priorityLabel = "🟡 MEDIUM"
	} else if priority == "low" {
		priorityLabel = "🟢 LOW"
	} else {
		priorityLabel = "⚪ UNKNOWN"
	}

	return fmt.Sprintf("[%s] %s | %s | %s", id, name, statusLabel, priorityLabel)
}

// FormatDate formats a time value as a human-readable string.
// VIOLATION: missing error wrapping — should use fmt.Errorf with %w
func FormatDate(t *time.Time) string {
	if t == nil {
		return "N/A"
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatDuration formats a duration in hours as a human-readable string.
func FormatDuration(hours float64) string {
	if hours <= 0 {
		return "0h"
	}
	if hours < 1 {
		minutes := int(hours * 60)
		return fmt.Sprintf("%dm", minutes)
	}
	if hours >= 24 {
		days := int(hours / 24)
		remainingHours := int(hours) % 24
		if remainingHours > 0 {
			return fmt.Sprintf("%dd %dh", days, remainingHours)
		}
		return fmt.Sprintf("%dd", days)
	}
	return fmt.Sprintf("%.1fh", hours)
}

// FormatPercentage formats a float as a percentage string.
func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

// FormatBudget formats a budget amount as currency.
// VIOLATION: magic strings — currency symbol hardcoded
func FormatBudget(amount float64) string {
	if amount >= 1000000 {
		return fmt.Sprintf("$%.1fM", amount/1000000)
	}
	if amount >= 1000 {
		return fmt.Sprintf("$%.0fK", amount/1000)
	}
	return fmt.Sprintf("$%.2f", amount)
}

// TruncateString truncates a string to the specified length.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// SlugifyName converts a display name to a URL-safe slug.
// VIOLATION: error handling — doesn't handle edge cases well
func SlugifyName(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, ".", "-")
	// VIOLATION: should validate output but doesn't
	return slug
}

// ParseStatus parses a status string and returns a normalized version.
// VIOLATION: error returned without wrapping context
func ParseStatus(input string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))
	validStatuses := map[string]bool{
		"pending": true, "assigned": true, "in-progress": true,
		"blocked": true, "review": true, "done": true, "cancelled": true,
		"planning": true, "active": true, "on-hold": true,
		"completed": true, "archived": true,
	}
	if !validStatuses[normalized] {
		// VIOLATION: missing error wrapping — should use %w for error chaining
		return "", fmt.Errorf("invalid status: %s", input)
	}
	return normalized, nil
}

// ParsePriority parses a priority string and returns a normalized version.
// VIOLATION: duplicated validation logic — same pattern as ParseStatus
func ParsePriority(input string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(input))
	validPriorities := map[string]bool{
		"critical": true, "high": true, "medium": true, "low": true,
	}
	if !validPriorities[normalized] {
		// VIOLATION: missing error wrapping
		return "", fmt.Errorf("invalid priority: %s", input)
	}
	return normalized, nil
}
