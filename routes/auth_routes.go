package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/app/service"
)

// LoginHandler handles login endpoint
func LoginHandler(c *fiber.Ctx) error {
	var loginReq models.LoginCredential

	// Parse request body
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Initialize repositories and service
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)

	// Call login service
	response, err := authService.Login(&loginReq)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// RegisterHandler handles user registration
func RegisterHandler(c *fiber.Ctx) error {
	var regReq models.LoginCredential

	// Parse request body
	if err := c.BodyParser(&regReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Initialize repositories and service
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)

	// Call register service
	response, err := authService.Register(&regReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// LogoutHandler handles logout endpoint
func LogoutHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "logout successful",
	})
}

// RefreshTokenHandler handles refresh token endpoint
func RefreshTokenHandler(c *fiber.Ctx) error {
	var refreshReq struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&refreshReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if refreshReq.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "token is required",
		})
	}

	// TODO: Implement refresh token logic
	// For now, return a placeholder response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "token refreshed",
		"token":   "new-token-here",
	})
}

// GetProfileHandler handles get user profile endpoint
func GetProfileHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	username := c.Locals("username")
	email := c.Locals("email")
	role := c.Locals("role")
	permissions := c.Locals("permissions")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":     userID,
		"username":    username,
		"email":       email,
		"role":        role,
		"permissions": permissions,
	})
}
