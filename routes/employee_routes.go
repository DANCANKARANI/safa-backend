package routes

import (
	"github.com/dancankarani/safa/controllers"
	"github.com/dancankarani/safa/middleware"
	"github.com/dancankarani/safa/models"
	"github.com/gofiber/fiber/v2"
)

func SetEmployeeRoutes(app *fiber.App) {
	// Define the routes for employee management
	e := app.Group("/api/v1", middleware.JWTMiddleware)
	e.Post("/employees", controllers.CreateEmployee)
	e.Get("/employees/:id", controllers.GetEmployeeByID)
	e.Get("/employees", controllers.GetAllEmployees)
	e.Patch("/employees/:id", controllers.UpdateEmployee)
	e.Delete("/employees/:id", controllers.DeleteEmployee)

	//expenses
	e.Post("/expenses", controllers.CreateExpensesHandler)
	e.Get("/expenses", controllers.GetExpensesHandler)
	e.Get("/daily/expenses", controllers.GetExpensesByDateHandler)
	e.Get("/station/expenses/:id", controllers.GetPaginatedExpensesByStation)
	e.Get("/expenses/duration/:start_date/:end_date", controllers.GetExpensesByDurationHandler)
	e.Patch("/expenses/:id", controllers.UpdateExpensesHandler)
	e.Delete("/expenses/:id", controllers.DeleteExpensesHandler)

	//pump readings
	e.Post("/pump-readings", controllers.AddNewPumpReadings)
	e.Get("/pump-readings", controllers.GetOrderedPumpReadingsHandler)
	e.Get("/pump-readings/:pump_id", controllers.GetOpeningReadingsHandler)
	e.Get("/pump-readings/:id", controllers.GetPumpReadingsHandler)
	e.Patch("/pump-readings/:id", controllers.UpdatePumpReadingsHandler)
	e.Delete("/pump-readings/:id", controllers.DeletePumpReadingsHandler)

	e.Get("/sales/", controllers.GetAllSalesByDateHandler)

	//send email
	em := app.Group("/api/v1")
	em.Post("/send-email", controllers.SendEmail)

	//daily accounts
	e.Post("/daily-accounts", models.AddDailyAccounts)
	e.Get("/daily-accounts", models.GetDailyAccounts)
	e.Get("/daily-accounts/:id", models.GetDailyAccount)
	e.Patch("/daily-accounts/:id", models.UpdateDailyAccount)
	e.Get("/monthly-accounts", models.GetMonthlyDailyAccounts)

}