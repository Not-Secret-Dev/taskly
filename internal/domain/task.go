// Package domain holds the todo model and the rules around it. It has no
// knowledge of storage formats or terminal rendering.
package domain

import (
	"strings"
	"time"
)

// Priority is the urgency of a task.
type Priority string

const (
	PriorityLow    Priority = "Low"
	PriorityNormal Priority = "Normal"
	PriorityHigh   Priority = "High"
)

// Task is a single todo item.
type Task struct {
	ID        uint32
	Title     string
	Priority  Priority
	Completed bool
	CreatedAt time.Time
	DueDate   time.Time
	UpdatedAt time.Time
}

// NormalizePriority maps free-form user input to a known priority. Anything
// that isn't clearly high or low becomes normal.
func NormalizePriority(p string) Priority {
	switch strings.ToLower(strings.TrimSpace(p)) {
	case "h", "high":
		return PriorityHigh
	case "l", "low":
		return PriorityLow
	default:
		return PriorityNormal
	}
}

// Stats summarizes a collection of tasks at a point in time.
type Stats struct {
	Total   int
	Done    int
	Overdue int
}

// Percent returns the completion percentage, or 0 for an empty collection.
func (s Stats) Percent() int {
	if s.Total == 0 {
		return 0
	}
	return s.Done * 100 / s.Total
}
