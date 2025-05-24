package handlers

import (
	"myapp/config"
	"myapp/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CreateUser handles user creation
func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// GetUser retrieves a user by ID
func GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID",
		})
	}

	var user models.User
	result := config.DB.First(&user, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// GetAllUsers retrieves all users
func GetAllUsers(c echo.Context) error {
	var users []models.User
	result := config.DB.Find(&users)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to retrieve users",
		})
	}

	return c.JSON(http.StatusOK, users)
}

// UpdateUser updates a user by ID
func UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID",
		})
	}

	var existingUser models.User
	result := config.DB.First(&existingUser, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	if err := c.Bind(&existingUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request payload",
		})
	}

	config.DB.Save(&existingUser)
	return c.JSON(http.StatusOK, existingUser)
}

// DeleteUser deletes a user by ID
func DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid user ID",
		})
	}

	result := config.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete user",
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}