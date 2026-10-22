# taskly ✦

A colorful terminal todo manager written in Go. No third-party dependencies, just a CSV file.

## Features

- Priorities (High / Normal / Low) with color-coded badges
- Due dates with overdue ⚠ and due-today ⏰ flags
- Progress bar and completion stats
- Add, list, toggle, and delete tasks
- Plain CSV storage, so you can open your list in Excel
- Layered codebase that is easy to extend with new backends or views

## Requirements

- Go 1.27.1 or newer

## Run

```sh
go run .
```

Run from the repository root. On first launch taskly creates `tasks.csv` in the current directory.

## Build

```sh
go build -o taskly .
./taskly
```

## Usage

The app is a menu-driven loop:

```
1  List tasks
2  Add task
3  Toggle complete
4  Delete task
5  Exit
```

When adding a task:

| Prompt   | Accepted input                                                        |
|----------|-----------------------------------------------------------------------|
| Title    | Any non-empty text                                                    |
| Priority | `h`/`high`, `l`/`low`, or anything else for Normal (the default)      |
| Due date | `YYYY-MM-DD`, or blank for tomorrow                                   |

Explicit due dates are interpreted in the local timezone and stored at 23:59 so they aren't marked overdue that morning. A blank due date defaults to 24 hours from now. Unparseable dates fall back to tomorrow.

Task IDs are assigned as one more than the highest existing ID, so the list can grow without collisions.

## Data

Tasks live in `tasks.csv`, created with its header on first run:

```csv
ID,Title,Priority,Completed,CreatedAt,DueDate,UpdatedAt
1,Write the README,High,false,2026-10-22T23:25:00+05:00,2026-10-23T23:59:00+05:00,2026-10-22T23:25:00+05:00
```

Timestamps use RFC 3339. Titles containing commas or quotes are handled by the standard CSV writer, so the file stays Excel-friendly.

> Tip: quit taskly before editing `tasks.csv` by hand to avoid overwriting your changes.

## Test

```sh
go test ./...
```

## Project layout

```
.
├── main.go                  entry point: wires a repository to the app
├── go.mod
└── internal
    ├── app
    │   └── app.go           menu loop and use cases (add, toggle, delete)
    ├── domain
    │   ├── task.go          Task model, priorities, stats
    │   ├── list.go          list rules: next ID, stats, toggle, remove
    │   └── list_test.go
    ├── storage
    │   ├── repository.go    Repository interface (the backend seam)
    │   ├── csv.go           CSV implementation
    │   └── csv_test.go
    └── ui
        ├── theme.go         ANSI-safe padding, centering, truncation
        ├── output.go        prompts and success/failure/info messages
        ├── menu.go          banner, menu, progress bar
        └── table.go         task table rendering
```

Dependencies point inward: `main` → `app` → `domain`, with `storage` and `ui` also depending only on `domain`. The domain package knows nothing about CSV or the terminal.

### Adding a backend

`storage.Repository` is the interface every backend implements:

```go
type Repository interface {
	Name() string
	Exists() bool
	Init() error
	Load() (domain.List, error)
	Append(task domain.Task) error
	Save(tasks domain.List) error
}
```

Implement those six methods (JSON, SQLite, a remote API, …) and pass your repository to `app.New` in `main.go`. No UI or domain changes required.

## Roadmap

- Sorting and filtering
- Search
- Edit task
