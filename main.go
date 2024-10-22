package main

import (
	"log"
	"sample/db"
	"sample/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	err = db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	time.Sleep(1 * time.Second)

	// Set up Fiber app
	app := fiber.New()

	app.Use(logger.New())

	// Register routes
	routes.SetupUserRoutes(app)

	// Start server
	// Start server
	err = app.Listen("0.0.0.0:3000")

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
