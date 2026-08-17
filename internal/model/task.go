package model

import "time"

type Task struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"project_id"`
	Title     string    `json:"title"`
	Priority  string    `json:"priority"`
	Status    string    `json:"status"`
	DueDate   string    `json:"due_date,omitempty"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
