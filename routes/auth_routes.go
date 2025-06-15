package routes

import (
	"github.com/dancankarani/safa/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetAuthRoutes(app *fiber.App) {
	// Define the routes for authentication
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", controllers.LoginUser)
	// auth.Get("/profile", controllers.GetProfile)
	// auth.Put("/profile", controllers.UpdateProfile)
	// auth.Delete("/profile", controllers.DeleteAccount)
}