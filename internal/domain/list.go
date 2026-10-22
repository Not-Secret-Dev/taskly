package domain

import "time"

// List is an ordered collection of tasks. Mutating helpers operate on the
// underlying slice so callers keep a single working set per request.
type List []Task

// NextID returns one more than the highest task ID, or 1 for an empty list.
func (l List) NextID() uint32 {
	var max uint32
	for _, t := range l {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

// Stats summarizes the list as of now.
func (l List) Stats(now time.Time) Stats {
	s := Stats{Total: len(l)}
	for _, t := range l {
		switch {
		case t.Completed:
			s.Done++
		case !t.DueDate.IsZero() && t.DueDate.Before(now):
			s.Overdue++
		}
	}
	return s
}

// Toggle flips the completion state of the task with the given ID, returning
// the new state. found is false when no task has that ID.
func (l List) Toggle(id uint32, now time.Time) (completed, found bool) {
	for i := range l {
		if l[i].ID == id {
			l[i].Completed = !l[i].Completed
			l[i].UpdatedAt = now
			return l[i].Completed, true
		}
	}
	return false, false
}

// Without returns a copy of the list with the task with the given ID removed.
// found is false when no task has that ID.
func (l List) Without(id uint32) (List, bool) {
	for i, t := range l {
		if t.ID != id {
			continue
		}
		out := make(List, 0, len(l)-1)
		out = append(out, l[:i]...)
		out = append(out, l[i+1:]...)
		return out, true
	}
	return l, false
}
