package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PaginationParams holds pagination information
type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

// GetPaginationParams extracts pagination parameters from request
func GetPaginationParams(c *fiber.Ctx) PaginationParams {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// SuccessResponse returns a success response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": message,
		"data":    data,
	})
}

// CreatedResponse returns a created response
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse returns an error response
func ErrorResponse(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  false,
		"message": message,
		"data":    nil,
	})
}

// ValidationErrorResponse returns validation error response
func ValidationErrorResponse(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  false,
		"message": "validation error",
		"error":   err.Error(),
		"data":    nil,
	})
}

// PaginatedResponse returns a paginated response
func PaginatedResponse(c *fiber.Ctx, data fiber.Map, total int64, page, limit int) error {
	totalPages := (total + int64(limit) - 1) / int64(limit)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "success",
		"data":    data,
		"pagination": fiber.Map{
			"total":       total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	})
}

// OKResponse returns a simple OK response
func OKResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": message,
		"data":    data,
	})
}

// DeletedResponse returns a deleted response
func DeletedResponse(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": message,
		"data":    nil,
	})
}
