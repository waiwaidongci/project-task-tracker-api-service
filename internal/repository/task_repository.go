package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"todo-api/internal/model"
)

func (r *SQLiteRepository) CreateTask(ctx context.Context, task *model.Task) error {
	tagsJSON, err := marshalTags(task.Tags)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (project_id, title, priority, status, due_date, tags)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		task.ProjectID,
		task.Title,
		task.Priority,
		task.Status,
		task.DueDate,
		tagsJSON,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = id

	created, err := r.GetTaskByID(ctx, id)
	if err != nil {
		return err
	}
	*task = *created
	return nil
}

func (r *SQLiteRepository) GetTaskByID(ctx context.Context, id int64) (*model.Task, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, title, priority, status, due_date, tags, created_at, updated_at
		 FROM tasks WHERE id = ?`,
		id,
	)
	task, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *SQLiteRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	tagsJSON, err := marshalTags(task.Tags)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx,
		`UPDATE tasks
		 SET project_id = ?, title = ?, priority = ?, status = ?, due_date = ?, tags = ?,
		     updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		task.ProjectID,
		task.Title,
		task.Priority,
		task.Status,
		task.DueDate,
		tagsJSON,
		task.ID,
	)
	if err != nil {
		return err
	}
	return requireAffected(result)
}

func (r *SQLiteRepository) UpdateTaskStatus(ctx context.Context, id int64, status string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		status,
		id,
	)
	if err != nil {
		return err
	}
	return requireAffected(result)
}

func (r *SQLiteRepository) DeleteTask(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireAffected(result)
}

func (r *SQLiteRepository) ListTasks(ctx context.Context, filter TaskFilter, page, pageSize int) ([]model.Task, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	limit, offset := buildPagination(page, pageSize)

	conditions := make([]string, 0, 4)
	args := make([]any, 0, 4)
	if filter.ProjectID > 0 {
		conditions = append(conditions, "project_id = ?")
		args = append(args, filter.ProjectID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Priority != "" {
		conditions = append(conditions, "priority = ?")
		args = append(args, filter.Priority)
	}
	if filter.DueDate != "" {
		conditions = append(conditions, "due_date = ?")
		args = append(args, filter.DueDate)
	}

	where := ""
	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM tasks" + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	query := `SELECT id, project_id, title, priority, status, due_date, tags, created_at, updated_at
		 FROM tasks` + where + ` ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer closeRows(rows)

	tasks := make([]model.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func scanTask(scanner interface{ Scan(dest ...any) error }) (model.Task, error) {
	var task model.Task
	var tagsJSON string
	err := scanner.Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Priority,
		&task.Status,
		&task.DueDate,
		&tagsJSON,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return task, err
	}
	task.Tags = unmarshalTags(tagsJSON)
	return task, nil
}

func marshalTags(tags []string) (string, error) {
	if tags == nil {
		tags = []string{}
	}
	data, err := json.Marshal(tags)
	if err != nil {
		return "", fmt.Errorf("marshal tags: %w", err)
	}
	return string(data), nil
}

func unmarshalTags(raw string) []string {
	if raw == "" || raw == "null" {
		return []string{}
	}
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil {
		return []string{}
	}
	if tags == nil {
		return []string{}
	}
	return tags
}

var _ TaskRepository = (*SQLiteRepository)(nil)
