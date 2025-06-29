package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"github.com/dancankarani/safa/routes"
	"github.com/gofiber/fiber/v2/middleware/cors"
)
func RegisterEndpoint() {
    app := fiber.New()

    // Enable CORS with allowed origins (allow all origins from the same network)
    app.Use(cors.New(cors.Config{
        AllowOriginsFunc: func(origin string) bool {
            // Allow localhost and any 192.168.110.* address
            return origin == "http://localhost:3000" ||
                origin == "http://127.0.0.1:3000" ||
                (len(origin) >= 22 && origin[:22] == "http://192.168.110.")
        },
        AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
    }))
    

    // Register routes
    routes.SetAuthRoutes(app)
    routes.SetSalaryAdvanceRoutes(app)  
    routes.SetAdminRoutes(app)
    routes.SetEmployeeRoutes(app)
    routes.FuelProductRoutes(app)

    app.Listen(":8000")
}
