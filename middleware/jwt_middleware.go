package middleware

import (
	"log"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/dancankarani/safa/utils"
	"github.com/dancankarani/safa/services"
)

func JWTMiddleware(c *fiber.Ctx) error {
    // Check for token in cookies first
    tokenString := c.Cookies("token")
    log.Println(tokenString)
    // If not found in cookies, check the Authorization header
    if tokenString == "" {
        authHeader := c.Get("Authorization")
        if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
            tokenString = strings.TrimPrefix(authHeader, "Bearer ")
        }
    }

    // If token is still not found, return unauthorized error
    if tokenString == "" {
        log.Println("missing jwt")
        return utils.NewErrorResponse(c, "unauthorized",map[string][]string{"error":{"missing jwt"}}, fiber.StatusUnauthorized)
    }

    // Validate the token
    claims, err := services.ValidateToken(tokenString)
    if err != nil {
        log.Println(err.Error())
        return utils.NewErrorResponse(c, "unauthorized",map[string][]string{"errors":{err.Error()}},fiber.StatusInternalServerError)
    }
    //get ip address and store in context
    ip := c.IP()
    c.Locals("ip_address", ip)
    // Store the userID in context
    c.Locals("user_id", claims.UserID)
    c.Locals("role",claims.Role)
    return c.Next()
}