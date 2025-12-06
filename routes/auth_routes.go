package routes

import (
	"github.com/gofiber/fiber/v2"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/app/service"
	"UAS/helpers"
	"UAS/middleware"
)

func SetupAuthRoutes(app *fiber.App) {
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(userRepo)
	g := app.Group("/api/v1/auth")

	// POST /login
	g.Post("/login", func(c *fiber.Ctx) error {
		var req models.LoginCredential
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
		}
		response, err := authService.Login(&req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.BuildErrorResponse(401, err.Error()))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, response))
	})

	// POST /register
	g.Post("/register", func(c *fiber.Ctx) error {
		var req models.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
		}
		userID, err := authService.Register(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, err.Error()))
		}
		return c.Status(fiber.StatusCreated).JSON(helpers.BuildCreatedResponse("user registered successfully", fiber.Map{"user_id": userID}))
	})

	// POST /logout
	g.Post("/logout", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("logout successful", nil))
	})

	// POST /refresh
	g.Post("/refresh", func(c *fiber.Ctx) error {
		var req struct {
			Token string `json:"token"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "invalid request body"))
		}
		if req.Token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(helpers.BuildErrorResponse(400, "token is required"))
		}
		return c.Status(fiber.StatusOK).JSON(helpers.BuildOKResponse("token refreshed", fiber.Map{"token": "new-token-here"}))
	})

	protected := g.Group("", middleware.AuthMiddleware)
	// GET /profile
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(helpers.BuildSuccessResponse(200, fiber.Map{
			"user_id":     c.Locals("user_id"),
			"username":    c.Locals("username"),
			"email":       c.Locals("email"),
			"role":        c.Locals("role"),
			"permissions": c.Locals("permissions"),
		}))
	})
}
