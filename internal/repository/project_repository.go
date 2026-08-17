package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"todo-api/internal/model"
)

func (r *SQLiteRepository) CreateProject(ctx context.Context, project *model.Project) error {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO projects (name, description) VALUES (?, ?)",
		project.Name,
		project.Description,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	project.ID = id

	created, err := r.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}
	*project = *created
	return nil
}

func (r *SQLiteRepository) GetProjectByID(ctx context.Context, id int64) (*model.Project, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, created_at, updated_at
		 FROM projects WHERE id = ?`,
		id,
	)
	project, err := scanProject(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *SQLiteRepository) UpdateProject(ctx context.Context, project *model.Project) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE projects
		 SET name = ?, description = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		project.Name,
		project.Description,
		project.ID,
	)
	if err != nil {
		return err
	}
	return requireAffected(result)
}

func (r *SQLiteRepository) DeleteProject(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireAffected(result)
}

func (r *SQLiteRepository) ListProjects(ctx context.Context, page, pageSize int) ([]model.Project, int, error) {
	limit, offset := buildPagination(page, pageSize)
	paginationArgs := appendPaginationArgs(nil, limit, offset)

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM projects").Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, created_at, updated_at
		 FROM projects
		 ORDER BY id DESC
		 LIMIT ? OFFSET ?`,
		paginationArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer closeRows(rows)

	projects := make([]model.Project, 0)
	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *SQLiteRepository) ProjectStats(ctx context.Context) ([]model.ProjectStats, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT p.id, p.name, COUNT(t.id)
		 FROM projects p
		 LEFT JOIN tasks t ON t.project_id = p.id AND t.status != ?
		 GROUP BY p.id, p.name
		 ORDER BY p.id`,
		model.StatusCompleted,
	)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	stats := make([]model.ProjectStats, 0)
	for rows.Next() {
		var item model.ProjectStats
		if err := rows.Scan(&item.ProjectID, &item.ProjectName, &item.UnfinishedCount); err != nil {
			return nil, err
		}
		stats = append(stats, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stats, nil
}

func scanProject(scanner interface{ Scan(dest ...any) error }) (model.Project, error) {
	var project model.Project
	err := scanner.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	return project, err
}

func requireAffected(result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

var _ ProjectRepository = (*SQLiteRepository)(nil)

// Keep time imported for generated schema timestamp types in older toolchains.
var _ = time.Time{}
