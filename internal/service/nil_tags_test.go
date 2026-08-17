package service

import (
	"context"
	"testing"

	"todo-api/internal/model"
)

func TestCreateTaskWithTagsDoesNotPanic(t *testing.T) {
	repo := newServiceTestRepository(t)
	project, err := NewProjectService(repo).Create(context.Background(), "Work", "")
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	tasks := NewTaskService(repo, repo)

	_, err = tasks.Create(context.Background(), TaskInput{
		ProjectID: project.ID,
		Title:     "write report",
		Priority:  model.PriorityMedium,
		Status:    model.StatusPending,
		Tags:      []string{"work", "report"},
	})
	if err != nil {
		t.Fatalf("create task with tags: %v", err)
	}
}
