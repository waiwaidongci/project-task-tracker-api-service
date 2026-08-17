package service

import (
	"context"
	"testing"

	"todo-api/internal/model"
)

func TestTaskListHonorsRequestedPageSize(t *testing.T) {
	repo := newServiceTestRepository(t)
	project, err := NewProjectService(repo).Create(context.Background(), "Work", "")
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	tasks := NewTaskService(repo, repo)
	for i := 0; i < 60; i++ {
		_, err := tasks.Create(context.Background(), TaskInput{
			ProjectID: project.ID,
			Title:     "task",
			Priority:  model.PriorityLow,
			Status:    model.StatusPending,
		})
		if err != nil {
			t.Fatalf("create task %d: %v", i, err)
		}
	}

	result, err := tasks.List(context.Background(), TaskListFilter{}, 1, 100)
	if err != nil {
		t.Fatalf("list tasks: %v", err)
	}
	if result.Pagination.PageSize != 100 || len(result.Data) != 60 {
		t.Fatalf("expected page_size=100 and 60 items, got page_size=%d len=%d", result.Pagination.PageSize, len(result.Data))
	}
}
