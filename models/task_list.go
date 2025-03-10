package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type TaskList struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTaskList(name string) *TaskList {
	return &TaskList{
		Id:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type TaskListOptions map[uuid.UUID]TaskList

func TaskListsToJson(options TaskListOptions) (string, error) {
	// Convert to a map with string keys
	strMap := make(map[string]TaskList)
	for key, value := range options {
		strMap[key.String()] = value
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(strMap, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func JsonToTaskLists(jsonStr string) (TaskListOptions, error) {
	// Temporary map with string keys
	strMap := make(TaskListOptions)

	err := json.Unmarshal([]byte(jsonStr), &strMap)
	if err != nil {
		return nil, err
	}

	uuidMap := make(map[uuid.UUID]TaskList)
	for key, value := range strMap {
		uuidMap[key] = value
	}

	return uuidMap, nil
}
