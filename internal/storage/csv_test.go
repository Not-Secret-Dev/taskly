package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"todo-cli/internal/domain"
)

func TestCSVLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.csv")
	repo := NewCSV(path)

	if repo.Exists() {
		t.Fatal("repository should not exist before Init")
	}
	if err := repo.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if !repo.Exists() {
		t.Fatal("repository should exist after Init")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read header: %v", err)
	}
	wantHeader := "ID,Title,Priority,Completed,CreatedAt,DueDate,UpdatedAt\n"
	if string(raw) != wantHeader {
		t.Errorf("header = %q, want %q", raw, wantHeader)
	}

	tasks, err := repo.Load()
	if err != nil {
		t.Fatalf("Load empty: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("empty store loaded %d tasks", len(tasks))
	}

	created := time.Date(2026, 10, 22, 9, 30, 0, 0, time.FixedZone("test", 5*60*60+30*60))
	want := domain.Task{
		ID:        1,
		Title:     "Write tests, then relax",
		Priority:  domain.PriorityHigh,
		CreatedAt: created,
		DueDate:   created.Add(24 * time.Hour),
		UpdatedAt: created,
	}
	if err := repo.Append(want); err != nil {
		t.Fatalf("Append: %v", err)
	}

	tasks, err = repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("loaded %d tasks, want 1", len(tasks))
	}
	assertTaskEqual(t, tasks[0], want)

	// Save rewrites the whole file and must keep the header intact.
	tasks[0].Completed = true
	tasks[0].UpdatedAt = created.Add(time.Hour)
	if err := repo.Save(tasks); err != nil {
		t.Fatalf("Save: %v", err)
	}

	raw, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if lines := strings.Split(strings.TrimSpace(string(raw)), "\n"); len(lines) != 2 {
		t.Fatalf("file has %d lines, want 2:\n%s", len(lines), raw)
	}
	if !strings.HasPrefix(string(raw), wantHeader) {
		t.Errorf("file lost its header:\n%s", raw)
	}

	tasks, err = repo.Load()
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}
	if len(tasks) != 1 || !tasks[0].Completed {
		t.Errorf("reloaded tasks = %+v, want one completed task", tasks)
	}
}

func TestCSVLoadRejectsMalformedRow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.csv")
	content := strings.Join([]string{
		"ID,Title,Priority,Completed,CreatedAt,DueDate,UpdatedAt",
		"1,Valid,Normal,false,2026-10-22T09:30:00Z,2026-10-23T09:30:00Z,2026-10-22T09:30:00Z",
		"2,Missing fields",
	}, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// encoding/csv rejects rows whose field count differs from the header;
	// Load surfaces that error rather than silently dropping data.
	if _, err := NewCSV(path).Load(); err == nil {
		t.Fatal("Load accepted a row with the wrong field count")
	}
}

func assertTaskEqual(t *testing.T, got, want domain.Task) {
	t.Helper()
	if got.ID != want.ID || got.Title != want.Title || got.Priority != want.Priority || got.Completed != want.Completed {
		t.Errorf("task = %+v, want %+v", got, want)
	}
	if !got.CreatedAt.Equal(want.CreatedAt) || !got.DueDate.Equal(want.DueDate) || !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("task timestamps = %v/%v/%v, want %v/%v/%v",
			got.CreatedAt, got.DueDate, got.UpdatedAt,
			want.CreatedAt, want.DueDate, want.UpdatedAt)
	}
}
