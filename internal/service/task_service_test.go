package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"todo-api/internal/model"
	"todo-api/internal/repository"
	"todo-api/migrations"
)

func newServiceTestRepository(t *testing.T) *repository.SQLiteRepository {
	t.Helper()
	db, err := repository.Open(filepath.Join(t.TempDir(), "service.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := migrations.Run(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return repository.NewSQLiteRepository(db)
}

func TestTaskServiceStatusChangeValidation(t *testing.T) {
	ctx := context.Background()
	repo := newServiceTestRepository(t)
	project, err := NewProjectService(repo).Create(ctx, "Work", "")
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	tasks := NewTaskService(repo, repo)

	task, err := tasks.Create(ctx, TaskInput{
		ProjectID: project.ID,
		Title:     "write tests",
		Priority:  model.PriorityMedium,
		Status:    model.StatusPending,
		Tags:      []string{"test", "test"},
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if len(task.Tags) != 1 {
		t.Fatalf("expected deduplicated tags, got %#v", task.Tags)
	}

	completed, err := tasks.ChangeStatus(ctx, task.ID, model.StatusCompleted)
	if err != nil {
		t.Fatalf("change status: %v", err)
	}
	if completed.Status != model.StatusCompleted {
		t.Fatalf("expected completed status, got %s", completed.Status)
	}

	if _, err := tasks.ChangeStatus(ctx, task.ID, model.StatusCompleted); err == nil {
		t.Fatal("expected same-status transition to fail")
	} else {
		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
	}

	if _, err := tasks.ChangeStatus(ctx, task.ID, "archived"); err == nil {
		t.Fatal("expected invalid status to fail")
	}
}

func TestTaskServiceTodayFilter(t *testing.T) {
	ctx := context.Background()
	repo := newServiceTestRepository(t)
	project, err := NewProjectService(repo).Create(ctx, "Work", "")
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	tasks := NewTaskService(repo, repo)
	tasks.now = func() time.Time { return time.Date(2026, 8, 16, 10, 0, 0, 0, time.Local) }

	for _, dueDate := range []string{"2026-08-16", "2026-08-17"} {
		_, err := tasks.Create(ctx, TaskInput{
			ProjectID: project.ID,
			Title:     "task due " + dueDate,
			Priority:  model.PriorityLow,
			Status:    model.StatusPending,
			DueDate:   dueDate,
		})
		if err != nil {
			t.Fatalf("create task %s: %v", dueDate, err)
		}
	}

	result, err := tasks.Today(ctx, 1, 20)
	if err != nil {
		t.Fatalf("today: %v", err)
	}
	if result.Pagination.Total != 1 || len(result.Data) != 1 {
		t.Fatalf("expected 1 task due today, got total=%d len=%d", result.Pagination.Total, len(result.Data))
	}
	if result.Data[0].DueDate != "2026-08-16" {
		t.Fatalf("unexpected task: %+v", result.Data[0])
	}
}

func TestTaskServiceCreateRequiresProject(t *testing.T) {
	ctx := context.Background()
	repo := newServiceTestRepository(t)
	tasks := NewTaskService(repo, repo)

	_, err := tasks.Create(ctx, TaskInput{
		ProjectID: 9999,
		Title:     "orphan",
		Priority:  model.PriorityLow,
		Status:    model.StatusPending,
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected repository not found, got %v", err)
	}
}
