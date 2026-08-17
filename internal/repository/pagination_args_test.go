package repository

import (
	"context"
	"testing"

	"todo-api/internal/model"
)

func TestFilteredTaskListKeepsPaginationArgsSeparate(t *testing.T) {
	ctx := context.Background()
	repo := newTestRepository(t)

	project := &model.Project{Name: "Work", Description: "work"}
	if err := repo.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	task := &model.Task{
		ProjectID: project.ID,
		Title:     "write report",
		Priority:  model.PriorityHigh,
		Status:    model.StatusPending,
		Tags:      []string{"work"},
	}
	if err := repo.CreateTask(ctx, task); err != nil {
		t.Fatalf("create task: %v", err)
	}

	items, total, err := repo.ListTasks(ctx, TaskFilter{Status: model.StatusPending}, 1, 20)
	if err != nil {
		t.Fatalf("list filtered tasks: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1 pending task, got total=%d len=%d", total, len(items))
	}
}
