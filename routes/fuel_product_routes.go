package routes

import (
	"github.com/dancankarani/safa/controllers"
	"github.com/dancankarani/safa/middleware"
	"github.com/dancankarani/safa/models"
	"github.com/gofiber/fiber/v2"
)

func FuelProductRoutes(app *fiber.App) {
	f := app.Group("/fuel-products", middleware.JWTMiddleware)
	f.Get("/", controllers.GetAllFuelProductsHandler)
	f.Get("/:id", controllers.GetFuelProductByIDHandler)
	f.Post("/", controllers.CreateFuelProductHandler)
	f.Put("/:id", controllers.UpdateFuelProductHandler)
	f.Delete("/:id", controllers.DeleteFuelProductHandler)

	fp:= app.Group("/api/v1/station/fuel-price", middleware.JWTMiddleware)
	fp.Post("/", models.CreateStationFuelPrice)
	fp.Get("/:station_id", models.GetStationFuelPrices)
	fp.Patch("/:id", models.UpdateStationFuelPrice)
	fp.Delete( "/:id", models.DeleteStationFuelPrice)
}