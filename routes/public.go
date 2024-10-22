package routes

import (
	"sample/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes sets up routes for user-related actions
func SetupUserRoutes(app *fiber.App) {
	app.Post("/register", handlers.RegisterUser) // Route to register a new user
	app.Post("/check-in", handlers.CheckIn)      // Route for attendance check-in
}
