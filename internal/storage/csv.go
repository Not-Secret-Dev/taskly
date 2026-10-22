package storage

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"todo-cli/internal/domain"
)

const filePerm = 0o644

var header = []string{"ID", "Title", "Priority", "Completed", "CreatedAt", "DueDate", "UpdatedAt"}

// CSV is a Repository backed by a plain CSV file.
type CSV struct {
	path string
}

// NewCSV returns a repository that stores tasks at path.
func NewCSV(path string) *CSV {
	return &CSV{path: path}
}

// Name returns the CSV file path.
func (c *CSV) Name() string { return c.path }

// Exists reports whether the CSV file is present.
func (c *CSV) Exists() bool {
	_, err := os.Stat(c.path)
	return err == nil
}

// Init creates the file and writes its header row.
func (c *CSV) Init() error {
	file, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(header)
}

// Load reads every task from the CSV file.
func (c *CSV) Load() (domain.List, error) {
	file, err := os.Open(c.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks domain.List
	for i, record := range records {
		if i == 0 || len(record) < len(header) {
			continue
		}

		id, _ := strconv.ParseUint(record[0], 10, 32)
		completed, _ := strconv.ParseBool(record[3])
		createdAt, _ := time.Parse(time.RFC3339, record[4])
		dueDate, _ := time.Parse(time.RFC3339, record[5])
		updatedAt, _ := time.Parse(time.RFC3339, record[6])

		tasks = append(tasks, domain.Task{
			ID:        uint32(id),
			Title:     record[1],
			Priority:  domain.Priority(record[2]),
			Completed: completed,
			CreatedAt: createdAt,
			DueDate:   dueDate,
			UpdatedAt: updatedAt,
		})
	}

	return tasks, nil
}

// Append writes one task as a new row.
func (c *CSV) Append(task domain.Task) error {
	file, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(record(task))
}

// Save replaces the file contents with header plus the given tasks.
func (c *CSV) Save(tasks domain.List) error {
	file, err := os.OpenFile(c.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}
	for _, task := range tasks {
		if err := writer.Write(record(task)); err != nil {
			return err
		}
	}

	return nil
}

// record converts a task into its CSV column order.
func record(t domain.Task) []string {
	return []string{
		strconv.FormatUint(uint64(t.ID), 10),
		t.Title,
		string(t.Priority),
		strconv.FormatBool(t.Completed),
		t.CreatedAt.Format(time.RFC3339),
		t.DueDate.Format(time.RFC3339),
		t.UpdatedAt.Format(time.RFC3339),
	}
}
