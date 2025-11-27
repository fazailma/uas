package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"UAS/database"
	
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found")
	}

	// Connect DB
	database.ConnectPostgres()

	// Create Fiber app
	app := fiber.New()



	// Start server
	log.Println("Server running on port 8080")
	app.Listen(":8080")
}
