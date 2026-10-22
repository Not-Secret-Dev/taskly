package domain

import (
	"testing"
	"time"
)

func testNow() time.Time {
	return time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)
}

func TestNormalizePriority(t *testing.T) {
	tests := map[string]Priority{
		"h":      PriorityHigh,
		"HIGH":   PriorityHigh,
		" high ": PriorityHigh,
		"l":      PriorityLow,
		"Low":    PriorityLow,
		"":       PriorityNormal,
		"n":      PriorityNormal,
		"urgent": PriorityNormal,
	}
	for in, want := range tests {
		if got := NormalizePriority(in); got != want {
			t.Errorf("NormalizePriority(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestListNextID(t *testing.T) {
	if got := (List{}).NextID(); got != 1 {
		t.Errorf("empty list NextID = %d, want 1", got)
	}
	list := List{{ID: 1}, {ID: 7}, {ID: 3}}
	if got := list.NextID(); got != 8 {
		t.Errorf("NextID = %d, want 8", got)
	}
}

func TestListStats(t *testing.T) {
	now := testNow()
	list := List{
		{ID: 1, Completed: true},
		{ID: 2, DueDate: now.Add(-time.Hour)},
		{ID: 3, DueDate: now.Add(time.Hour)},
		{ID: 4}, // no due date is never overdue
	}

	got := list.Stats(now)
	want := Stats{Total: 4, Done: 1, Overdue: 1}
	if got != want {
		t.Errorf("Stats = %+v, want %+v", got, want)
	}
	if got.Percent() != 25 {
		t.Errorf("Percent = %d, want 25", got.Percent())
	}
	if p := (Stats{}).Percent(); p != 0 {
		t.Errorf("empty Percent = %d, want 0", p)
	}
}

func TestListToggle(t *testing.T) {
	list := List{{ID: 1}, {ID: 2}}
	now := testNow()

	completed, found := list.Toggle(2, now)
	if !found || !completed {
		t.Fatalf("Toggle(2) = (%v, %v), want (true, true)", completed, found)
	}
	if !list[1].Completed || !list[1].UpdatedAt.Equal(now) {
		t.Errorf("task 2 not updated: %+v", list[1])
	}

	if _, found := list.Toggle(99, now); found {
		t.Error("Toggle(99) reported found")
	}
}

func TestListWithout(t *testing.T) {
	list := List{{ID: 1}, {ID: 2}, {ID: 3}}

	remaining, found := list.Without(2)
	if !found {
		t.Fatal("Without(2) reported not found")
	}
	if len(remaining) != 2 || remaining[0].ID != 1 || remaining[1].ID != 3 {
		t.Errorf("Without(2) = %+v", remaining)
	}
	if len(list) != 3 {
		t.Errorf("Without mutated the original list: %+v", list)
	}

	if _, found := list.Without(99); found {
		t.Error("Without(99) reported found")
	}
}
