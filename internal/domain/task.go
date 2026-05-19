package domain

import "time"

type Task struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description string
	Status      TaskStatus
	Deadline    *time.Time
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type TaskStatus string

const (
	NewTaskStatus        TaskStatus = "NEW"
	InProgressTaskStatus TaskStatus = "IN_PROGRESS"
	DoneTaskStatus       TaskStatus = "DONE"
)

type TaskFilter struct {
	Status       *TaskStatus
	DeadlineFrom *time.Time
	DeadlineTo   *time.Time
	SortBy       string
	SortOrder    string
}
