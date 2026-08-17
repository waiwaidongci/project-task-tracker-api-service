package service

import (
	"context"
	"strings"

	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type ProjectService struct {
	projects repository.ProjectRepository
}

func NewProjectService(projects repository.ProjectRepository) *ProjectService {
	return &ProjectService{projects: projects}
}

func (s *ProjectService) Create(ctx context.Context, name, description string) (*model.Project, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return nil, validation("name", "project name is required")
	}
	if len(name) > 100 {
		return nil, validation("name", "project name must not exceed 100 characters")
	}
	if len(description) > 500 {
		return nil, validation("description", "project description must not exceed 500 characters")
	}

	project := &model.Project{Name: name, Description: description}
	if err := s.projects.CreateProject(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	return s.projects.GetProjectByID(ctx, id)
}

func (s *ProjectService) Update(ctx context.Context, id int64, name, description string) (*model.Project, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return nil, validation("name", "project name is required")
	}
	if len(name) > 100 {
		return nil, validation("name", "project name must not exceed 100 characters")
	}
	if len(description) > 500 {
		return nil, validation("description", "project description must not exceed 500 characters")
	}

	current, err := s.projects.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	current.Name = name
	current.Description = description
	if err := s.projects.UpdateProject(ctx, current); err != nil {
		return nil, err
	}
	return s.projects.GetProjectByID(ctx, id)
}

func (s *ProjectService) Delete(ctx context.Context, id int64) error {
	return s.projects.DeleteProject(ctx, id)
}

func (s *ProjectService) List(ctx context.Context, page, pageSize int) (model.ListResponse[model.Project], error) {
	page, pageSize = repository.NormalizePagination(page, pageSize)
	items, total, err := s.projects.ListProjects(ctx, page, pageSize)
	if err != nil {
		return model.ListResponse[model.Project]{}, err
	}
	return model.ListResponse[model.Project]{
		Data:       items,
		Pagination: model.NewPagination(page, pageSize, total),
	}, nil
}

func (s *ProjectService) Stats(ctx context.Context) ([]model.ProjectStats, error) {
	return s.projects.ProjectStats(ctx)
}
