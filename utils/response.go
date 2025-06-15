package utils

import "github.com/gofiber/fiber/v2"

func SuccessResponse(c *fiber.Ctx,message string, data interface{}) error {
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   data,
		"message": message,
	})
}

func NewErrorResponse(c *fiber.Ctx, message string, errors map[string][]string, statusCode int) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "fail",
		"message": message,
		"error": errors,
	})
}

func NotFoundResponse(c *fiber.Ctx, message string) error {
	err_str := map[string][]string{
		"error": {message},
	}
	return NewErrorResponse(c,"not found",err_str , fiber.StatusNotFound)
}
func BadRequestResponse(c *fiber.Ctx, message string) error {
	err_str := map[string][]string{
		"error": {message},	
	}
	return NewErrorResponse(c,"bad request",err_str , fiber.StatusBadRequest)
}

// SendMessage sends a simple success message in JSON format
// It can be used for simple acknowledgments or confirmations.

func SendMessage(c *fiber.Ctx, message string) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": message,
	})
}