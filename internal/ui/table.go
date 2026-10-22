package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"todo-cli/internal/domain"
)

var (
	taskColumns = []int{4, 28, 10, 11, 14}
	taskHeaders = []string{"ID", "Title", "Priority", "Status", "Due"}
)

// TaskTable renders the task list as a bordered table.
func TaskTable(tasks domain.List) {
	fmt.Println()
	if len(tasks) == 0 {
		fmt.Printf("  %s✨ No tasks yet. Your list is empty!%s\n", yellow, reset)
		return
	}

	line := func(l, m, r string) {
		parts := make([]string, len(taskColumns))
		for i, w := range taskColumns {
			parts[i] = strings.Repeat("─", w+2)
		}
		fmt.Println("  " + gray + l + strings.Join(parts, m) + r + reset)
	}
	row := func(cells []string) {
		parts := make([]string, len(cells))
		for i, c := range cells {
			parts[i] = pad(c, taskColumns[i])
		}
		sep := gray + " │ " + reset
		fmt.Println("  " + gray + "│ " + reset + strings.Join(parts, sep) + gray + " │" + reset)
	}

	fmt.Printf("  %s%s YOUR TASKS%s\n", bold, "📋", reset)
	line("┌", "┬", "┐")

	hdr := make([]string, len(taskHeaders))
	for i, h := range taskHeaders {
		hdr[i] = bold + h + reset
	}
	row(hdr)
	line("├", "┼", "┤")

	now := time.Now()
	today := now.Format("2006-01-02")

	for _, t := range tasks {
		row([]string{
			dim + strconv.FormatUint(uint64(t.ID), 10) + reset,
			titleCell(t),
			priorityCell(t.Priority),
			statusCell(t.Completed),
			dueCell(t, now, today),
		})
	}
	line("└", "┴", "┘")
}

func statusCell(completed bool) string {
	if completed {
		return green + "✔ Done" + reset
	}
	return yellow + "○ Pending" + reset
}

func priorityCell(p domain.Priority) string {
	switch p {
	case domain.PriorityHigh:
		return red + "● High" + reset
	case domain.PriorityLow:
		return cyan + "● Low" + reset
	default:
		return yellow + "● Normal" + reset
	}
}

func titleCell(t domain.Task) string {
	title := truncate(t.Title, taskColumns[1])
	if t.Completed {
		title = gray + strike + title + reset
	}
	return title
}

func dueCell(t domain.Task, now time.Time, today string) string {
	if t.DueDate.IsZero() || t.DueDate.Year() <= 1 {
		return gray + "No due date" + reset
	}

	ds := t.DueDate.Format("2006-01-02")
	switch {
	case t.Completed:
		return gray + ds + reset
	case t.DueDate.Before(now) && ds != today:
		return red + ds + " ⚠" + reset
	case ds == today:
		return yellow + ds + " ⏰" + reset
	default:
		return ds
	}
}
