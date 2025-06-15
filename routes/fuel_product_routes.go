package routes

import (
	"github.com/dancankarani/safa/controllers"
	"github.com/gofiber/fiber/v2"
)

func FuelProductRoutes(app *fiber.App) {
	f := app.Group("/fuel-products")
	f.Get("/", controllers.GetAllFuelProductsHandler)
	f.Get("/:id", controllers.GetFuelProductByIDHandler)
	f.Post("/", controllers.CreateFuelProductHandler)
	f.Put("/:id", controllers.UpdateFuelProductHandler)
	f.Delete("/:id", controllers.DeleteFuelProductHandler)
}