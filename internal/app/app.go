// Package app wires the domain rules and storage layer to the terminal UI. It
// owns the interactive loop but no rendering or persistence details.
package app

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"todo-cli/internal/domain"
	"todo-cli/internal/storage"
	"todo-cli/internal/ui"
)

const (
	// defaultDueIn is used when the user leaves the due date blank.
	defaultDueIn = 24 * time.Hour
	// endOfDay pushes a parsed date to 23:59 so it isn't overdue that morning.
	endOfDay = 23*time.Hour + 59*time.Minute
)

// App runs the interactive todo loop.
type App struct {
	repo    storage.Repository
	scanner *bufio.Scanner
}

// New returns an App backed by repo that reads commands from in.
func New(repo storage.Repository, in io.Reader) *App {
	return &App{repo: repo, scanner: bufio.NewScanner(in)}
}

// Run starts the main loop and returns when the user exits or input ends.
func (a *App) Run() {
	if !a.repo.Exists() {
		if err := a.repo.Init(); err == nil {
			ui.FileCreated(a.repo.Name())
		}
	}

	for {
		tasks, err := a.repo.Load()
		if err != nil {
			ui.Failure("Error loading tasks: " + err.Error())
			return
		}

		ui.Clear()
		ui.Banner()
		ui.StatsBar(tasks.Stats(time.Now()))
		ui.Menu()
		ui.ChoosePrompt()

		if !a.scanner.Scan() {
			return
		}

		switch strings.TrimSpace(a.scanner.Text()) {
		case "1":
			ui.TaskTable(tasks)
			ui.Pause(a.scanner)
		case "2":
			a.add(tasks)
		case "3":
			a.toggle(tasks)
		case "4":
			a.delete(tasks)
		case "5":
			ui.Goodbye()
			return
		default:
			ui.Failure("Invalid option, please choose 1-5.")
			ui.Pause(a.scanner)
		}
	}
}

// add prompts for the fields of a new task and appends it to the store.
func (a *App) add(tasks domain.List) {
	fmt.Println()
	title := ui.Prompt(a.scanner, "Task title: ")
	if title == "" {
		ui.Failure("Title cannot be empty.")
		ui.Pause(a.scanner)
		return
	}
	priority := domain.NormalizePriority(ui.Prompt(a.scanner, "Priority (High/Normal/Low) [Normal]: "))

	dueDate := time.Now().Add(defaultDueIn)
	if s := ui.Prompt(a.scanner, "Due date (YYYY-MM-DD) [tomorrow]: "); s != "" {
		d, err := time.ParseInLocation("2006-01-02", s, time.Local)
		if err != nil {
			ui.Info("Couldn't read that date, defaulting to tomorrow.")
		} else {
			dueDate = d.Add(endOfDay)
		}
	}

	now := time.Now()
	task := domain.Task{
		ID:        tasks.NextID(),
		Title:     title,
		Priority:  priority,
		CreatedAt: now,
		DueDate:   dueDate,
		UpdatedAt: now,
	}

	if err := a.repo.Append(task); err != nil {
		ui.Failure("Error adding task: " + err.Error())
	} else {
		ui.Success(fmt.Sprintf("Added \"%s\" (ID: %d)", task.Title, task.ID))
	}
	ui.Pause(a.scanner)
}

// toggle flips the completion state of the requested task and saves the list.
func (a *App) toggle(tasks domain.List) {
	ui.TaskTable(tasks)
	id, err := strconv.ParseUint(ui.Prompt(a.scanner, "Task ID to toggle: "), 10, 32)
	if err != nil {
		ui.Failure("Invalid Task ID.")
		ui.Pause(a.scanner)
		return
	}

	completed, found := tasks.Toggle(uint32(id), time.Now())
	if !found {
		ui.Failure(fmt.Sprintf("Task with ID %d not found.", id))
		ui.Pause(a.scanner)
		return
	}

	if err := a.repo.Save(tasks); err != nil {
		ui.Failure("Error updating task: " + err.Error())
		ui.Pause(a.scanner)
		return
	}

	if completed {
		ui.Success(fmt.Sprintf("Task %d marked as completed. Nice work! 🎉", id))
	} else {
		ui.Success(fmt.Sprintf("Task %d reopened.", id))
	}
	ui.Pause(a.scanner)
}

// delete removes the requested task and saves the remaining list.
func (a *App) delete(tasks domain.List) {
	ui.TaskTable(tasks)
	id, err := strconv.ParseUint(ui.Prompt(a.scanner, "Task ID to delete: "), 10, 32)
	if err != nil {
		ui.Failure("Invalid Task ID.")
		ui.Pause(a.scanner)
		return
	}

	remaining, found := tasks.Without(uint32(id))
	if !found {
		ui.Failure(fmt.Sprintf("Task with ID %d not found.", id))
		ui.Pause(a.scanner)
		return
	}

	if err := a.repo.Save(remaining); err != nil {
		ui.Failure("Error deleting task: " + err.Error())
		ui.Pause(a.scanner)
		return
	}

	ui.Success(fmt.Sprintf("Task %d deleted.", id))
	ui.Pause(a.scanner)
}
