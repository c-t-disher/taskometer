package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"taskometer/handlers"
	s "taskometer/services"
)

func main() {
	s.InitS3Client()

	// ensure task-lists.json file is set and accurate
	_, tlErr := s.GetTaskListOptions(context.Background())
	if tlErr != nil {
		log.Fatal("Failed to get task list options:", tlErr)
	}

	r := gin.Default()

	r.POST("/tasks", handlers.CreateTask)
	r.POST("/lists", handlers.CreateTaskList)
	r.GET("/tasks/:id", handlers.GetTasks)
	r.GET("/lists", handlers.GetTaskLists)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
