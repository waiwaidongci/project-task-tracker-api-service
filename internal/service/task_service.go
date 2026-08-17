package service

import (
	"context"
	"strings"
	"time"

	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type TaskInput struct {
	ProjectID int64
	Title     string
	Priority  string
	Status    string
	DueDate   string
	Tags      []string
}

type TaskListFilter struct {
	ProjectID int64
	Status    string
	Priority  string
	DueToday  bool
}

type TaskService struct {
	tasks    repository.TaskRepository
	projects repository.ProjectRepository
	now      func() time.Time
}

func NewTaskService(tasks repository.TaskRepository, projects repository.ProjectRepository) *TaskService {
	return &TaskService{
		tasks:    tasks,
		projects: projects,
		now:      time.Now,
	}
}

func (s *TaskService) Create(ctx context.Context, input TaskInput) (*model.Task, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.DueDate = strings.TrimSpace(input.DueDate)
	if err := s.validateInput(input, ""); err != nil {
		return nil, err
	}
	if _, err := s.projects.GetProjectByID(ctx, input.ProjectID); err != nil {
		return nil, err
	}

	task := &model.Task{
		ProjectID: input.ProjectID,
		Title:     input.Title,
		Priority:  input.Priority,
		Status:    input.Status,
		DueDate:   input.DueDate,
		Tags:      normalizeTags(input.Tags),
	}
	if err := s.tasks.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	return s.tasks.GetTaskByID(ctx, id)
}

func (s *TaskService) Update(ctx context.Context, id int64, input TaskInput) (*model.Task, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.DueDate = strings.TrimSpace(input.DueDate)
	current, err := s.tasks.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.validateInput(input, ""); err != nil {
		return nil, err
	}
	if _, err := s.projects.GetProjectByID(ctx, input.ProjectID); err != nil {
		return nil, err
	}

	current.ProjectID = input.ProjectID
	current.Title = input.Title
	current.Priority = input.Priority
	current.Status = input.Status
	current.DueDate = input.DueDate
	current.Tags = normalizeTags(input.Tags)

	if err := s.tasks.UpdateTask(ctx, current); err != nil {
		return nil, err
	}
	return s.tasks.GetTaskByID(ctx, id)
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	return s.tasks.DeleteTask(ctx, id)
}

func (s *TaskService) List(ctx context.Context, filter TaskListFilter, page, pageSize int) (model.ListResponse[model.Task], error) {
	if filter.Status != "" && !model.ValidStatus(filter.Status) {
		return model.ListResponse[model.Task]{}, validation("status", "invalid status filter")
	}
	if filter.Priority != "" && !model.ValidPriority(filter.Priority) {
		return model.ListResponse[model.Task]{}, validation("priority", "invalid priority filter")
	}

	repoFilter := repository.TaskFilter{
		ProjectID: filter.ProjectID,
		Status:    filter.Status,
		Priority:  filter.Priority,
	}
	if filter.DueToday {
		repoFilter.DueDate = s.now().Format("2006-01-02")
	}

	page, pageSize = normalizePagination(page, pageSize)
	items, total, err := s.tasks.ListTasks(ctx, repoFilter, page, pageSize)
	if err != nil {
		return model.ListResponse[model.Task]{}, err
	}
	return model.ListResponse[model.Task]{
		Data:       items,
		Pagination: model.NewPagination(page, pageSize, total),
	}, nil
}

func (s *TaskService) Today(ctx context.Context, page, pageSize int) (model.ListResponse[model.Task], error) {
	return s.List(ctx, TaskListFilter{DueToday: true}, page, pageSize)
}

func (s *TaskService) ChangeStatus(ctx context.Context, id int64, status string) (*model.Task, error) {
	if !model.ValidStatus(status) {
		return nil, validation("status", "invalid task status")
	}

	current, err := s.tasks.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canTransition(current.Status, status) {
		return nil, validation("status", "invalid status transition")
	}

	if err := s.tasks.UpdateTaskStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.tasks.GetTaskByID(ctx, id)
}

func (s *TaskService) validateInput(input TaskInput, currentStatus string) error {
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		return validation("title", "task title is required")
	}
	if len(input.Title) > 200 {
		return validation("title", "task title must not exceed 200 characters")
	}
	if input.ProjectID <= 0 {
		return validation("project_id", "project_id is required")
	}
	if !model.ValidPriority(input.Priority) {
		return validation("priority", "priority must be one of low, medium, high")
	}
	if !model.ValidStatus(input.Status) {
		return validation("status", "status must be one of pending, in_progress, completed")
	}
	if input.DueDate != "" && !isValidDate(input.DueDate) {
		return validation("due_date", "due_date must use YYYY-MM-DD format")
	}
	tags := normalizeTags(input.Tags)
	if len(tags) > 10 {
		return validation("tags", "at most 10 tags are allowed")
	}
	for _, tag := range tags {
		if tag == "" || len(tag) > 30 {
			return validation("tags", "each tag must be between 1 and 30 characters")
		}
	}
	if currentStatus != "" && !canTransition(currentStatus, input.Status) {
		return validation("status", "invalid status transition")
	}
	return nil
}

func normalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, exists := seen[tag]; exists {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}
	return result
}

func isValidDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

func canTransition(from, to string) bool {
	if from == to {
		return false
	}
	transitions := map[string]map[string]bool{
		model.StatusPending: {
			model.StatusInProgress: true,
			model.StatusCompleted:  true,
		},
		model.StatusInProgress: {
			model.StatusPending:   true,
			model.StatusCompleted: true,
		},
		model.StatusCompleted: {
			model.StatusPending:    true,
			model.StatusInProgress: true,
		},
	}
	return transitions[from][to]
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
