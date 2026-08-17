package service

import (
	"context"
	"testing"

	"todo-api/internal/model"
)

func TestTaskListStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	repo := newServiceTestRepository(t)
	project, err := NewProjectService(repo).Create(context.Background(), "Work", "")
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	tasks := NewTaskService(repo, repo)
	_, err = tasks.Create(context.Background(), TaskInput{
		ProjectID: project.ID,
		Title:     "finish report",
		Priority:  model.PriorityMedium,
		Status:    model.StatusPending,
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	cancel()

	if _, err := tasks.List(ctx, TaskListFilter{}, 1, 20); err == nil {
		t.Fatal("expected canceled context to stop task list")
	}
}
