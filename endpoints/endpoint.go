package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"github.com/dancankarani/safa/routes"
	"github.com/gofiber/fiber/v2/middleware/cors"
)
func RegisterEndpoint() {
    app := fiber.New()

    // Enable CORS with default config (allows all origins)
    app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:3000",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
    }))
    
    

    // Register routes
    routes.SetEmployeeRoutes(app)
    routes.SetAuthRoutes(app)
    routes.SetSalaryAdvanceRoutes(app)
    routes.SetAdminRoutes(app)

    app.Listen(":8000")
}
