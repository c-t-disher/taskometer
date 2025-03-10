package handlers

import (
	"encoding/json"
	"net/http"
	"taskometer/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	m "taskometer/models"
)

func CreateTaskList(c *gin.Context) {
	// Parse request body as raw JSON map
	var requestBody map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&requestBody); err != nil {
		setErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	name, nameOk := requestBody["name"].(string)
	if !nameOk {
		setErrorResponse(c, http.StatusBadRequest, "Missing or invalid attribute.")
		return
	}

	// Generate a new UUID for the task list
	taskList := m.NewTaskList(name)

	// Save task list to S3
	err := services.SaveTaskList(c, *taskList)
	if err != nil {
		setErrorResponse(c, http.StatusInternalServerError, "Failed to save task list")
		return
	}

	c.JSON(http.StatusCreated, taskList)
}

func CreateTask(c *gin.Context) {
	// Parse request body as raw JSON map
	var requestBody map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&requestBody); err != nil {
		setErrorResponse(c, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Extract fields manually
	listId, listIdOk := requestBody["list_id"].(string)
	title, titleOk := requestBody["title"].(string)
	if !titleOk || !listIdOk {
		setErrorResponse(c, http.StatusBadRequest, "Missing or invalid attribute.")
		return
	}

	// Generate a new UUID for the task
	task := m.NewTask(
		uuid.MustParse(listId),
		title,
		time.Now().Add(24*time.Hour),
	)

	// Save task to S3
	err := services.SaveTask(c, *task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func GetTaskLists(c *gin.Context) {
	taskLists, err := services.GetTaskListOptions(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list task lists"})
		return
	}

	c.JSON(http.StatusOK, taskLists)
}

func GetTasks(c *gin.Context) {
	listId := uuid.MustParse(c.Param("id"))
	tasks, err := services.ListTasks(c, listId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func DeleteTask(c *gin.Context) {
	listId := uuid.MustParse(c.Param("list_id"))
	taskId := uuid.MustParse(c.Param("task_id"))
	err := services.DeleteTask(c, listId, taskId)
	if err != nil {
		setErrorResponse(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func setErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
