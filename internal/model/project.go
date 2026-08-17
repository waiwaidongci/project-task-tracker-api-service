package model

import "time"

type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectStats struct {
	ProjectID       int64  `json:"project_id"`
	ProjectName     string `json:"project_name"`
	UnfinishedCount int    `json:"unfinished_count"`
}
