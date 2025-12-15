package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	"UAS/database"
	_ "UAS/docs"
	"UAS/routes"
)

// @title Achievement Management Backend API
// @version 1.0
// @description API for managing student achievements with role-based access control
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @basePath /api/v1
// @schemes http https

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found")
	}

	// Connect PostgreSQL
	database.ConnectPostgres()

	// Connect MongoDB
	database.ConnectMongoDB()
	defer database.DisconnectMongoDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app)

	// Swagger UI
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "UAS Backend API is running",
			"status":  "ok",
		})
	})

	// Start server
	log.Println("Server running on http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
