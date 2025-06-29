package routes

import (
	"github.com/dancankarani/safa/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetAuthRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", controllers.LoginUser)
}