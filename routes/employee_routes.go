package routes

import (
	"github.com/dancankarani/safa/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetEmployeeRoutes(app *fiber.App) {
	// Define the routes for employee management
	e := app.Group("/api/v1")
	e.Post("/employees", controllers.CreateEmployee)
	e.Get("/employees/:id", controllers.GetEmployeeByID)
	e.Get("/employees", controllers.GetAllEmployees)
	e.Patch("/employees/:id", controllers.UpdateEmployee)
	e.Delete("/employees/:id", controllers.DeleteEmployee)

	//expenses
	e.Post("/expenses", controllers.CreateExpensesHandler)
	e.Get("/expenses", controllers.GetExpensesHandler)
	e.Get("/daily/expenses", controllers.GetExpensesByDateHandler)
	e.Get("/expenses/duration/:start_date/:end_date", controllers.GetExpensesByDurationHandler)
	e.Patch("/expenses/:id", controllers.UpdateExpensesHandler)
	e.Delete("/expenses/:id", controllers.DeleteExpensesHandler)

	//pump readings
	e.Post("/pump-readings", controllers.AddNewPumpReadings)
	e.Get("/pump-readings", controllers.GetPumpReadingsHandler)
	e.Get("/pump-readings/:id", controllers.GetPumpsByStationHandler)
	e.Patch("/pump-readings/:id", controllers.UpdatePumpReadingsHandler)
	e.Delete("/pump-readings/:id", controllers.DeletePumpReadingsHandler)

	e.Get("/sales/", controllers.GetAllSalesByDateHandler)
	
}