package config

import (
	"fmt"
	"log"
	"os"

	"myapp/models"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get credentials from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Validate environment variables
	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("Database configuration environment variables are not set properly")
	}

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successful!")

	// Disable foreign key checks temporarily for migrations
	DB.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// AutoMigrate models
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate User model: %v", err)
	}

	err = DB.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("Failed to migrate Task model: %v", err)
	}

	// Enable foreign key checks back
	DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	fmt.Println("Database migration successful!")
}