# Todo API

A two-level todo API built with Go, Gin, and SQLite.

## Requirements

- Go 1.26+
- Gin
- SQLite (modernc.org/sqlite, no CGO required)

## Run

```sh
go run ./cmd/server
```

The server listens on port `18005` by default. Set `PORT` or `DB_PATH` to override:

```sh
PORT=18006 DB_PATH=/tmp/todo.db go run ./cmd/server
```

## Project Layout

```text
cmd/server/main.go
internal/config
internal/model
internal/repository
internal/service
internal/handler
internal/router
internal/middleware
migrations
```

## Data Model

Projects contain a name and description. Tasks belong to a project and contain a title, priority, status, due date, and tags.

Status values:

- `pending`
- `in_progress`
- `completed`

Priority values:

- `low`
- `medium`
- `high`

Due dates use `YYYY-MM-DD`.

## Endpoints

### Projects

| Method | Path | Description |
| --- | --- | --- |
| POST | `/api/v1/projects` | Create a project |
| GET | `/api/v1/projects` | List projects with pagination |
| GET | `/api/v1/projects/stats` | List unfinished task counts per project |
| GET | `/api/v1/projects/:id` | Get a project |
| PUT | `/api/v1/projects/:id` | Update a project |
| DELETE | `/api/v1/projects/:id` | Delete a project and its tasks |

### Tasks

| Method | Path | Description |
| --- | --- | --- |
| POST | `/api/v1/tasks` | Create a task |
| GET | `/api/v1/tasks` | List tasks, optionally filtered |
| GET | `/api/v1/tasks/today` | List tasks due today |
| GET | `/api/v1/tasks/:id` | Get a task |
| PUT | `/api/v1/tasks/:id` | Update a task |
| PATCH | `/api/v1/tasks/:id/status` | Change task status through service validation |
| DELETE | `/api/v1/tasks/:id` | Delete a task |

Task list query filters:

- `project_id`
- `status`
- `priority`
- `due_today=true`

Pagination:

- `page`, default `1`
- `page_size`, default `20`, maximum `100`

## Request Examples

Create a project:

```sh
curl -X POST http://127.0.0.1:18005/api/v1/projects \
  -H 'Content-Type: application/json' \
  -d '{"name":"Work","description":"Work tasks"}'
```

Create a task:

```sh
curl -X POST http://127.0.0.1:18005/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"project_id":1,"title":"Finish report","priority":"high","status":"pending","due_date":"2026-08-16","tags":["work","report"]}'
```

Change status:

```sh
curl -X PATCH http://127.0.0.1:18005/api/v1/tasks/1/status \
  -H 'Content-Type: application/json' \
  -d '{"status":"in_progress"}'
```

## Error Format

All errors use the same JSON shape:

```json
{
  "error": {
    "code": "validation_error",
    "message": "task title is required"
  }
}
```

## Tests

```sh
go test ./...
```
