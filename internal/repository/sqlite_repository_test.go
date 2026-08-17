package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"todo-api/internal/model"
	"todo-api/migrations"
)

func newTestRepository(t *testing.T) *SQLiteRepository {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := migrations.Run(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return NewSQLiteRepository(db)
}

func TestTaskRepositoryFiltersAndStats(t *testing.T) {
	ctx := context.Background()
	repo := newTestRepository(t)

	project := &model.Project{Name: "Work", Description: "work items"}
	if err := repo.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	otherProject := &model.Project{Name: "Home", Description: "home items"}
	if err := repo.CreateProject(ctx, otherProject); err != nil {
		t.Fatalf("create other project: %v", err)
	}

	tasks := []model.Task{
		{ProjectID: project.ID, Title: "finish report", Priority: model.PriorityHigh, Status: model.StatusPending, DueDate: "2026-08-16", Tags: []string{"work"}},
		{ProjectID: project.ID, Title: "review code", Priority: model.PriorityMedium, Status: model.StatusCompleted, DueDate: "2026-08-17", Tags: []string{"work"}},
		{ProjectID: otherProject.ID, Title: "buy milk", Priority: model.PriorityHigh, Status: model.StatusPending, DueDate: "2026-08-16", Tags: []string{"home"}},
	}
	for i := range tasks {
		if err := repo.CreateTask(ctx, &tasks[i]); err != nil {
			t.Fatalf("create task %d: %v", i, err)
		}
	}

	pending, total, err := repo.ListTasks(ctx, TaskFilter{Status: model.StatusPending}, 1, 20)
	if err != nil {
		t.Fatalf("list pending: %v", err)
	}
	if total != 2 || len(pending) != 2 {
		t.Fatalf("expected 2 pending tasks, got total=%d len=%d", total, len(pending))
	}

	high, total, err := repo.ListTasks(ctx, TaskFilter{Priority: model.PriorityHigh}, 1, 20)
	if err != nil {
		t.Fatalf("list high: %v", err)
	}
	if total != 2 || len(high) != 2 {
		t.Fatalf("expected 2 high priority tasks, got total=%d len=%d", total, len(high))
	}

	today, total, err := repo.ListTasks(ctx, TaskFilter{DueDate: "2026-08-16"}, 1, 20)
	if err != nil {
		t.Fatalf("list today: %v", err)
	}
	if total != 2 || len(today) != 2 {
		t.Fatalf("expected 2 tasks due today, got total=%d len=%d", total, len(today))
	}

	stats, err := repo.ProjectStats(ctx)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 project stats, got %d", len(stats))
	}
	if stats[0].ProjectID != project.ID || stats[0].UnfinishedCount != 1 {
		t.Fatalf("unexpected work stats: %+v", stats[0])
	}
	if stats[1].ProjectID != otherProject.ID || stats[1].UnfinishedCount != 1 {
		t.Fatalf("unexpected home stats: %+v", stats[1])
	}
}

func TestTaskRepositoryPagination(t *testing.T) {
	ctx := context.Background()
	repo := newTestRepository(t)

	project := &model.Project{Name: "Pagination"}
	if err := repo.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	for i := 0; i < 5; i++ {
		task := &model.Task{
			ProjectID: project.ID,
			Title:     "task",
			Priority:  model.PriorityLow,
			Status:    model.StatusPending,
		}
		if err := repo.CreateTask(ctx, task); err != nil {
			t.Fatalf("create task: %v", err)
		}
	}

	items, total, err := repo.ListTasks(ctx, TaskFilter{}, 2, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total 5, got %d", total)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items on page 2, got %d", len(items))
	}
}

func TestScanTimeFields(t *testing.T) {
	ctx := context.Background()
	repo := newTestRepository(t)
	project := &model.Project{Name: "Time"}
	if err := repo.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	if project.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be scanned")
	}
	if project.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be scanned")
	}

	createdAt := project.CreatedAt
	time.Sleep(10 * time.Millisecond)
	project.Description = "updated"
	if err := repo.UpdateProject(ctx, project); err != nil {
		t.Fatalf("update project: %v", err)
	}
	refreshed, err := repo.GetProjectByID(ctx, project.ID)
	if err != nil {
		t.Fatalf("get project: %v", err)
	}
	if refreshed.CreatedAt.IsZero() || refreshed.UpdatedAt.Before(createdAt) {
		t.Fatalf("unexpected time fields: created=%v updated=%v", refreshed.CreatedAt, refreshed.UpdatedAt)
	}
}
