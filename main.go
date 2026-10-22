package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// ───────────────────────── Styling ─────────────────────────

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	strike  = "\033[9m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[90m"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// visibleLen counts printable runes, ignoring ANSI escape codes.
func visibleLen(s string) int {
	return utf8.RuneCountInString(ansiRe.ReplaceAllString(s, ""))
}

// pad right-pads s to w visible characters (ANSI-safe).
func pad(s string, w int) string {
	n := w - visibleLen(s)
	if n < 0 {
		n = 0
	}
	return s + strings.Repeat(" ", n)
}

// center centers s within w visible characters (ANSI-safe).
func center(s string, w int) string {
	n := w - visibleLen(s)
	if n <= 0 {
		return s
	}
	left := n / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", n-left)
}

func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	return string(r[:w-1]) + "…"
}

func clearScreen() { fmt.Print("\033[H\033[2J") }

func success(msg string) { fmt.Printf("\n  %s✔ %s%s\n", green, msg, reset) }
func failure(msg string) { fmt.Printf("\n  %s✖ %s%s\n", red, msg, reset) }
func info(msg string)    { fmt.Printf("\n  %sℹ %s%s\n", cyan, msg, reset) }

func pause(scanner *bufio.Scanner) {
	fmt.Printf("\n  %sPress Enter to continue…%s", gray, reset)
	scanner.Scan()
}

func prompt(scanner *bufio.Scanner, label string) string {
	fmt.Printf("  %s▸%s %s", magenta, reset, label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func printBanner() {
	const w = 38
	fmt.Println(cyan + "╭" + strings.Repeat("─", w) + "╮" + reset)
	fmt.Println(cyan + "│" + reset + center(bold+"✦  T O D O   M A N A G E R  ✦"+reset, w) + cyan + "│" + reset)
	fmt.Println(cyan + "╰" + strings.Repeat("─", w) + "╯" + reset)
}

func progressBar(done, total, width int) string {
	if total == 0 {
		return gray + strings.Repeat("░", width) + reset
	}
	filled := done * width / total
	return green + strings.Repeat("█", filled) + reset +
	gray + strings.Repeat("░", width-filled) + reset
}

func printStats(tasks []Task) {
	done, overdue := 0, 0
	now := time.Now()
	for _, t := range tasks {
		if t.Completed {
			done++
		} else if !t.DueDate.IsZero() && t.DueDate.Before(now) {
			overdue++
		}
	}
	pct := 0
	if len(tasks) > 0 {
		pct = done * 100 / len(tasks)
	}
	fmt.Printf("  %s %s%3d%%%s  %s%d/%d done%s", progressBar(done, len(tasks), 20), bold, pct, reset, dim, done, len(tasks), reset)
	if overdue > 0 {
		fmt.Printf("  %s⚠ %d overdue%s", red, overdue, reset)
	}
	fmt.Println()
}

func printMenu() {
	fmt.Println()
	fmt.Printf("  %s1%s  📋  List tasks\n", bold+cyan, reset)
	fmt.Printf("  %s2%s  ➕  Add task\n", bold+cyan, reset)
	fmt.Printf("  %s3%s  🔁  Toggle complete\n", bold+cyan, reset)
	fmt.Printf("  %s4%s  🗑   Delete task\n", bold+cyan, reset)
	fmt.Printf("  %s5%s  👋  Exit\n", bold+cyan, reset)
	fmt.Println()
}

// ───────────────────────── Model ─────────────────────────

type Task struct {
	ID        uint32
	Title     string
	Priority  string
	Completed bool
	CreatedAt time.Time
	DueDate   time.Time
	UpdatedAt time.Time
}

func normalizePriority(p string) string {
	switch strings.ToLower(strings.TrimSpace(p)) {
		case "h", "high":
			return "High"
		case "l", "low":
			return "Low"
		default:
			return "Normal"
	}
}

// ───────────────────────── Main ─────────────────────────

func main() {
	filename := "tasks.csv"
	if !fileExists(filename) {
		writeContent(filename)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		tasks, err := loadTasks(filename)
		if err != nil {
			failure("Error loading tasks: " + err.Error())
			break
		}

		clearScreen()
		printBanner()
		printStats(tasks)
		printMenu()
		fmt.Printf("  %sChoose an option (1-5):%s ", dim, reset)

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
			case "1":
				showTasks(tasks)
				pause(scanner)

			case "2":
				fmt.Println()
				title := prompt(scanner, "Task title: ")
				if title == "" {
					failure("Title cannot be empty.")
					pause(scanner)
					break
				}
				priority := normalizePriority(prompt(scanner, "Priority (High/Normal/Low) [Normal]: "))

				dueDate := time.Now().Add(24 * time.Hour)
				if s := prompt(scanner, "Due date (YYYY-MM-DD) [tomorrow]: "); s != "" {
					d, perr := time.ParseInLocation("2006-01-02", s, time.Local)
					if perr != nil {
						info("Couldn't read that date, defaulting to tomorrow.")
					} else {
						dueDate = d.Add(23*time.Hour + 59*time.Minute)
					}
				}

				if err := addTask(filename, title, priority, dueDate); err != nil {
					failure("Error adding task: " + err.Error())
				}
				pause(scanner)

			case "3":
				showTasks(tasks)
				parsedID, err := strconv.ParseUint(prompt(scanner, "Task ID to toggle: "), 10, 32)
				if err != nil {
					failure("Invalid Task ID.")
					pause(scanner)
					break
				}
				if err := updateTask(filename, uint32(parsedID), tasks); err != nil {
					failure("Error updating task: " + err.Error())
				}
				pause(scanner)

			case "4":
				showTasks(tasks)
				parsedID, err := strconv.ParseUint(prompt(scanner, "Task ID to delete: "), 10, 32)
				if err != nil {
					failure("Invalid Task ID.")
					pause(scanner)
					break
				}
				if err := deleteTask(filename, uint32(parsedID), tasks); err != nil {
					failure("Error deleting task: " + err.Error())
				}
				pause(scanner)

			case "5":
				fmt.Printf("\n  %s✦ Goodbye! Stay productive. ✦%s\n\n", cyan, reset)
				return

			default:
				failure("Invalid option, please choose 1-5.")
				pause(scanner)
		}
	}
}

// ───────────────────────── Storage ─────────────────────────

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func loadTasks(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []Task

	for i, record := range records {
		if i == 0 || len(record) < 7 {
			continue
		}

		parsedID, _ := strconv.ParseUint(record[0], 10, 32)
		parsedCompleted, _ := strconv.ParseBool(record[3])
		createdAt, _ := time.Parse(time.RFC3339, record[4])
		dueDate, _ := time.Parse(time.RFC3339, record[5])
		updatedAt, _ := time.Parse(time.RFC3339, record[6])

		tasks = append(tasks, Task{
			ID:        uint32(parsedID),
			       Title:     record[1],
			       Priority:  record[2],
			       Completed: parsedCompleted,
			       CreatedAt: createdAt,
			       DueDate:   dueDate,
			       UpdatedAt: updatedAt,
		})
	}

	return tasks, nil
}

func writeContent(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Title", "Priority", "Completed", "CreatedAt", "DueDate", "UpdatedAt"}
	if err = writer.Write(headers); err != nil {
		return err
	}

	fmt.Printf("%sCreated %s and wrote headers.%s\n", dim, filename, reset)
	return nil
}

func addTask(filename string, title string, priority string, dueDate time.Time) error {
	existingTasks, err := loadTasks(filename)
	if err != nil {
		existingTasks = []Task{}
	}

	var maxID uint32 = 0
	for _, t := range existingTasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newID := maxID + 1

	now := time.Now()
	newTask := Task{
		ID:        newID,
		Title:     title,
		Priority:  priority,
		Completed: false,
		CreatedAt: now,
		DueDate:   dueDate,
		UpdatedAt: now,
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.FormatUint(uint64(newTask.ID), 10),
		newTask.Title,
		newTask.Priority,
		strconv.FormatBool(newTask.Completed),
		newTask.CreatedAt.Format(time.RFC3339),
		newTask.DueDate.Format(time.RFC3339),
		newTask.UpdatedAt.Format(time.RFC3339),
	}

	if err = writer.Write(record); err != nil {
		return err
	}

	success(fmt.Sprintf("Added \"%s\" (ID: %d)", newTask.Title, newTask.ID))
	return nil
}

func saveTasks(filename string, tasks []Task) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Title", "Priority", "Completed", "CreatedAt", "DueDate", "UpdatedAt"}
	if err := writer.Write(headers); err != nil {
		return err
	}
	for _, t := range tasks {
		record := []string{
			strconv.FormatUint(uint64(t.ID), 10),
			t.Title,
			t.Priority,
			strconv.FormatBool(t.Completed),
			t.CreatedAt.Format(time.RFC3339),
			t.DueDate.Format(time.RFC3339),
			t.UpdatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// ───────────────────────── Display ─────────────────────────

func showTasks(tasks []Task) error {
	fmt.Println()
	if len(tasks) == 0 {
		fmt.Printf("  %s✨ No tasks yet. Your list is empty!%s\n", yellow, reset)
		return nil
	}

	widths := []int{4, 28, 10, 11, 14}
	headers := []string{"ID", "Title", "Priority", "Status", "Due"}

	line := func(l, m, r string) {
		parts := make([]string, len(widths))
		for i, w := range widths {
			parts[i] = strings.Repeat("─", w+2)
		}
		fmt.Println("  " + gray + l + strings.Join(parts, m) + r + reset)
	}
	row := func(cells []string) {
		parts := make([]string, len(cells))
		for i, c := range cells {
			parts[i] = pad(c, widths[i])
		}
		sep := gray + " │ " + reset
		fmt.Println("  " + gray + "│ " + reset + strings.Join(parts, sep) + gray + " │" + reset)
	}

	fmt.Printf("  %s%s YOUR TASKS%s\n", bold, "📋", reset)
	line("┌", "┬", "┐")
	hdr := make([]string, len(headers))
	for i, h := range headers {
		hdr[i] = bold + h + reset
	}
	row(hdr)
	line("├", "┼", "┤")

	now := time.Now()
	today := now.Format("2006-01-02")

	for _, t := range tasks {
		// Status
		status := yellow + "○ Pending" + reset
		if t.Completed {
			status = green + "✔ Done" + reset
		}

		// Priority badge
		var prio string
		switch t.Priority {
			case "High":
				prio = red + "● High" + reset
			case "Low":
				prio = cyan + "● Low" + reset
			default:
				prio = yellow + "● Normal" + reset
		}

		// Title
		title := truncate(t.Title, widths[1])
		if t.Completed {
			title = gray + strike + title + reset
		}

		// Due date
		due := gray + "No due date" + reset
		if !t.DueDate.IsZero() && t.DueDate.Year() > 1 {
			ds := t.DueDate.Format("2006-01-02")
			switch {
				case t.Completed:
					due = gray + ds + reset
				case t.DueDate.Before(now) && ds != today:
					due = red + ds + " ⚠" + reset
				case ds == today:
					due = yellow + ds + " ⏰" + reset
				default:
					due = ds
			}
		}

		row([]string{
			dim + strconv.FormatUint(uint64(t.ID), 10) + reset,
		    title, prio, status, due,
		})
	}
	line("└", "┴", "┘")
	return nil
}

// ───────────────────────── Mutations ─────────────────────────

func updateTask(filename string, taskID uint32, tasks []Task) error {
	found := false
	nowDone := false
	for i := range tasks {
		if tasks[i].ID == taskID {
			tasks[i].Completed = !tasks[i].Completed
			tasks[i].UpdatedAt = time.Now()
			nowDone = tasks[i].Completed
			found = true
			break
		}
	}

	if !found {
		failure(fmt.Sprintf("Task with ID %d not found.", taskID))
		return nil
	}
	if err := saveTasks(filename, tasks); err != nil {
		return err
	}

	if nowDone {
		success(fmt.Sprintf("Task %d marked as completed. Nice work! 🎉", taskID))
	} else {
		success(fmt.Sprintf("Task %d reopened.", taskID))
	}
	return nil
}

func deleteTask(filename string, taskID uint32, tasks []Task) error {
	var updatedTasks []Task
	found := false
	for _, task := range tasks {
		if task.ID == taskID {
			found = true
			continue
		}
		updatedTasks = append(updatedTasks, task)
	}

	if !found {
		failure(fmt.Sprintf("Task with ID %d not found.", taskID))
		return nil
	}
	if err := saveTasks(filename, updatedTasks); err != nil {
		return err
	}

	success(fmt.Sprintf("Task %d deleted.", taskID))
	return nil
}
