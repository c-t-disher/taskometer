package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
	m "taskometer/models"
	"time"
)

const taskLists = "task-lists.json"
const taskPrefix = "tasks"

func GetTaskListOptions(ctx context.Context) (m.TaskListOptions, error) {
	body, err := GetObject(ctx, taskLists)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			err := initTaskListsOptions(ctx)
			if err != nil {
				return nil, err
			}
			body, err = GetObject(ctx, taskLists)
		} else {
			return nil, err
		}
	} else {
		defer closeBody(body)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return m.JsonToTaskLists(buf.String())
}

func initTaskListsOptions(ctx context.Context) error {
	jsonStr, err := m.TaskListsToJson(make(m.TaskListOptions))
	if err != nil {
		return err
	}

	// Upload the JSON file
	err = PutObj(ctx, taskLists, bytes.NewReader([]byte(jsonStr)))
	if err != nil {
		return fmt.Errorf("failed to init %s: %w", taskLists, err)
	}

	fmt.Println("Successfully initialized " + taskLists)
	return nil
}

func SaveTaskList(ctx context.Context, taskList m.TaskList) error {
	options, err := GetTaskListOptions(ctx)
	if err != nil {
		return err
	}
	options[taskList.Id] = taskList

	jsonStr, err := m.TaskListsToJson(options)
	if err != nil {
		return err
	}

	return PutObj(ctx, taskLists, bytes.NewReader([]byte(jsonStr)))
}

func SaveTask(ctx context.Context, task m.Task) error {
	task.UpdatedAt = time.Now()
	jsonData, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return PutObj(ctx, formatKey(task.ListId, task.Id), bytes.NewReader(jsonData))
}

func GetTask(ctx context.Context, listId uuid.UUID, taskId uuid.UUID) *m.Task {
	body, err := GetObject(ctx, formatKey(listId, taskId))
	if err != nil {
		return nil
	}
	defer closeBody(body)

	var task m.Task
	err = json.NewDecoder(body).Decode(&task)
	if err != nil {
		return nil
	}

	return &task
}

func closeBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Println("Failed to close response body:", err)
	}
}

func ListTasks(ctx context.Context, listId uuid.UUID) ([]m.Task, error) {
	resp, err := ListObjects(ctx, taskPrefix+"/"+listId.String())
	if err != nil {
		return nil, err
	}

	var tasks = make([]m.Task, 0, len(resp.Contents))
	for _, item := range resp.Contents {
		segments := strings.Split(*item.Key, "/")
		taskId := uuid.MustParse(segments[2])
		task := GetTask(ctx, listId, taskId)
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func DeleteTask(ctx context.Context, listId uuid.UUID, taskId uuid.UUID) error {
	return DeleteObj(ctx, formatKey(listId, taskId))
}

func formatKey(listId uuid.UUID, taskId uuid.UUID) string {
	return taskPrefix + "/" + listId.String() + "/" + taskId.String()
}
