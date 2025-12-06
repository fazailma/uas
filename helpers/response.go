package helpers

import "github.com/gofiber/fiber/v2"

// BuildSuccessResponse membangun response sukses standard
func BuildSuccessResponse(statusCode int, data interface{}) fiber.Map {
	return fiber.Map{
		"status": "success",
		"code":   statusCode,
		"data":   data,
	}
}

// BuildErrorResponse membangun response error standard
func BuildErrorResponse(statusCode int, message string) fiber.Map {
	return fiber.Map{
		"status":  "error",
		"code":    statusCode,
		"message": message,
	}
}

// BuildCreatedResponse membangun response created dengan message
func BuildCreatedResponse(message string, data interface{}) fiber.Map {
	return fiber.Map{
		"status":  "success",
		"code":    201,
		"message": message,
		"data":    data,
	}
}

// BuildOKResponse membangun response OK dengan message
func BuildOKResponse(message string, data interface{}) fiber.Map {
	return fiber.Map{
		"status":  "success",
		"code":    200,
		"message": message,
		"data":    data,
	}
}

// BuildDeletedResponse membangun response untuk deleted resource
func BuildDeletedResponse(message string) fiber.Map {
	return fiber.Map{
		"status":  "success",
		"code":    200,
		"message": message,
	}
}
