package main

import (
	"net/http"

	"myapp/config"
	"myapp/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize Echo
	e := echo.New()

	// Initialize Database
	config.InitDB()

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to Task Management API")
	})

	// Setup routes
	routes.SetupUserRoutes(e)
	routes.SetupTaskRoutes(e)

	// Add database check endpoint
	e.GET("/check-db", func(c echo.Context) error {
		// Try to get the underlying SQL DB object
		sqlDB, err := config.DB.DB()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status":  "error",
				"message": "Failed to get database instance: " + err.Error(),
			})
		}

		// Try to ping the database
		err = sqlDB.Ping()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status":  "error",
				"message": "Failed to ping database: " + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status":  "success",
			"message": "Database connection is working properly",
		})
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}