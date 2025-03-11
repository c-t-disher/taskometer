package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
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
	r.GET("/lists", handlers.GetTaskLists)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.LoadHTMLGlob("main/templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/", HomeView)
	r.GET("/tasks/:id", TaskListView)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}

func HomeView(c *gin.Context) {
	options, err := s.GetTaskListOptions(c)
	if err != nil {
		log.Fatal("Failed to get task list options:", err)
	}
	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"title": "The Taskometer",
		"lists": options,
	})
}

func TaskListView(c *gin.Context) {
	listId := uuid.MustParse(c.Param("id"))
	tasks, err := s.ListTasks(c, listId)
	if err != nil {
		log.Fatal("Failed to get task list options:", err)
	}
	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"listName": "The Taskometer",
		"tasks":    tasks,
	})
}
