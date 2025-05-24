package routes

import (
	"myapp/handlers"

	"github.com/labstack/echo/v4"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(e *echo.Echo) {
	e.POST("/users", handlers.CreateUser)
	e.GET("/users/:id", handlers.GetUser)
	e.GET("/users", handlers.GetAllUsers)
	e.PUT("/users/:id", handlers.UpdateUser)
	e.DELETE("/users/:id", handlers.DeleteUser)
}