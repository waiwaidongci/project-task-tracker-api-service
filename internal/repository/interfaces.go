package repository

import (
	"context"

	"todo-api/internal/model"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *model.Project) error
	GetProjectByID(ctx context.Context, id int64) (*model.Project, error)
	UpdateProject(ctx context.Context, project *model.Project) error
	DeleteProject(ctx context.Context, id int64) error
	ListProjects(ctx context.Context, page, pageSize int) ([]model.Project, int, error)
	ProjectStats(ctx context.Context) ([]model.ProjectStats, error)
}

type TaskFilter struct {
	ProjectID int64
	Status    string
	Priority  string
	DueDate   string
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetTaskByID(ctx context.Context, id int64) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	UpdateTaskStatus(ctx context.Context, id int64, status string) error
	DeleteTask(ctx context.Context, id int64) error
	ListTasks(ctx context.Context, filter TaskFilter, page, pageSize int) ([]model.Task, int, error)
}
