package models

import (
	"github.com/google/uuid"
	"time"
)

type TaskStatus int

const (
	pending TaskStatus = iota
	complete
	cancelled
)

var taskStatusName = map[TaskStatus]string{
	pending:   "Pending",
	complete:  "Complete",
	cancelled: "Cancelled",
}

func (s TaskStatus) Name() string {
	return taskStatusName[s]
}

type Task struct {
	Id          uuid.UUID  `json:"id"`
	ListId      uuid.UUID  `json:"list_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	DueDate     time.Time  `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewTask(listId uuid.UUID, title string, dueDate time.Time) *Task {
	return &Task{
		Id:          uuid.New(),
		ListId:      listId,
		Title:       title,
		Description: "",
		Status:      pending,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
