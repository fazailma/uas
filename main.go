package main

import (
	"log"

	"UAS/database"
	_ "UAS/docs" // Import docs untuk Swagger (underscore karena hanya butuh side effect)
	"UAS/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	swagger "github.com/swaggo/fiber-swagger" // Ganti import swagger
)

// @title Achievement Management Backend API
// @version 1.0
// @description API untuk manajemen prestasi akademik mahasiswa

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter your JWT token in the format: Bearer {token}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
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
	app.Use(recover.New())
	app.Use(cors.New())

	// Swagger endpoint - Ganti dengan swagger handler yang benar
	app.Get("/swagger/*", swagger.WrapHandler)

	// Serve static files for uploads
	app.Static("/uploads", "./uploads")

	// Setup routes
	routes.SetupRoutes(app)

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "UAS Backend API is running",
			"status":  "ok",
		})
	})

	// Start server
	log.Println("Server running on http://localhost:8080")
	log.Println("Swagger documentation: http://localhost:8080/swagger/index.html")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
