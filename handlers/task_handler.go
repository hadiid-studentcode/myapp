package handlers

import (
	"net/http"
	"strconv"

	"myapp/config"
	"myapp/models"

	"github.com/labstack/echo/v4"
)

// CreateTask creates a new task
func CreateTask(c echo.Context) error {
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	// Validate required fields
	if task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Title is required",
		})
	}

	if task.UserID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "User ID is required",
		})
	}

	// Check if user exists
	var user models.User
	if err := config.DB.First(&user, task.UserID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID - user not found",
		})
	}

	// Set default status if not provided
	if task.Status == "" {
		task.Status = "pending"
	}

	if result := config.DB.Create(&task); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create task",
			"error":   result.Error.Error(),
		})
	}

	// Fetch the created task with user information
	config.DB.Preload("User").First(&task, task.ID)
	return c.JSON(http.StatusCreated, task)
}

// GetTask retrieves a task by ID
func GetTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid task ID",
		})
	}

	var task models.Task
	result := config.DB.Preload("User").First(&task, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Task not found",
		})
	}

	return c.JSON(http.StatusOK, task)
}

// GetAllTasks retrieves all tasks with optional filters
func GetAllTasks(c echo.Context) error {
	var tasks []models.Task
	db := config.DB.Preload("User")

	// Handle optional status filter
	if status := c.QueryParam("status"); status != "" {
		db = db.Where("status = ?", status)
	}

	// Handle optional user filter
	if userID := c.QueryParam("user_id"); userID != "" {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Find(&tasks)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to retrieve tasks",
			"error":   result.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, tasks)
}

// UpdateTask updates a task by ID
func UpdateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid task ID",
		})
	}

	var existingTask models.Task
	if err := config.DB.First(&existingTask, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Task not found",
		})
	}

	// Store the old values
	oldUserID := existingTask.UserID

	// Bind new data
	if err := c.Bind(&existingTask); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	// Validate title
	if existingTask.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Title cannot be empty",
		})
	}

	// If user ID has changed, verify the new user exists
	if existingTask.UserID != oldUserID {
		var user models.User
		if err := config.DB.First(&user, existingTask.UserID).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid user ID - user not found",
			})
		}
	}

	// Update task
	if err := config.DB.Save(&existingTask).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to update task",
			"error":   err.Error(),
		})
	}

	// Fetch updated task with user information
	config.DB.Preload("User").First(&existingTask, id)
	return c.JSON(http.StatusOK, existingTask)
}

// DeleteTask deletes a task by ID
func DeleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid task ID",
		})
	}

	// Check if task exists
	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "Task not found",
		})
	}

	// Delete the task
	if err := config.DB.Delete(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete task",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Task deleted successfully",
	})
}

// GetUserTasks retrieves all tasks for a specific user with optional status filter
func GetUserTasks(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID",
		})
	}

	// Verify user exists
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	var tasks []models.Task
	db := config.DB.Where("user_id = ?", userID).Preload("User")

	// Handle optional status filter
	if status := c.QueryParam("status"); status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to retrieve tasks",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tasks)
}