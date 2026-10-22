// Package storage persists tasks. Repository is the seam that lets the
// application work against any backend (CSV today, JSON/SQLite/API later)
// without changing the UI or domain layers.
package storage

import "todo-cli/internal/domain"

// Repository stores and retrieves tasks.
type Repository interface {
	// Name identifies the store in user-facing messages.
	Name() string
	// Exists reports whether the backing store has been created yet.
	Exists() bool
	// Init creates an empty store.
	Init() error
	// Load returns every task in the store.
	Load() (domain.List, error)
	// Append writes a single task to the store.
	Append(task domain.Task) error
	// Save replaces the entire contents of the store.
	Save(tasks domain.List) error
}
