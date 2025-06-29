package routes
import(
	"github.com/dancankarani/safa/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/dancankarani/safa/middleware"
)

func SetSalaryAdvanceRoutes(app *fiber.App) {
	// Define the routes for salary advance management
	e := app.Group("/api/v1/salary/advances",middleware.JWTMiddleware)
	//protected routes
	e.Post("/", controllers.CreateSalaryAdvanceHandler)
	e.Get("/:id", controllers.GetSalaryAdvanceByIDHandler)
	e.Patch("/:id", controllers.UpdateSalaryAdvanceHandler)
	e.Delete("/:id", controllers.DeleteSalaryAdvanceHandler)
}