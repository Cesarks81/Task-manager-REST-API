package models

import "time"

var AcceptedTaskStatus = map[string]bool{
	"ongoing":   true,
	"completed": true,
	"new":       true,
}

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	CompletedAt *time.Time `json:"completedAt"`
}
