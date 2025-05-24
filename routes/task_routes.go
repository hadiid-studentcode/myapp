package routes

import (
	"myapp/handlers"

	"github.com/labstack/echo/v4"
)

// SetupTaskRoutes configures all task-related routes
func SetupTaskRoutes(e *echo.Echo) {
	e.POST("/tasks", handlers.CreateTask)
	e.GET("/tasks/:id", handlers.GetTask)
	e.GET("/tasks", handlers.GetAllTasks)
	e.PUT("/tasks/:id", handlers.UpdateTask)
	e.DELETE("/tasks/:id", handlers.DeleteTask)
	e.GET("/users/:userId/tasks", handlers.GetUserTasks)
}