package main

import (
	"log"

	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Set up the DB connection
	setupDB()

	// Auto-migrate the schema
	db.AutoMigrate(&models.User{}, &models.Election{}, &models.Question{}, &models.Vote{}, &models.TransitionApplication{}, &models.News{})

	// Initialize Fiber
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	}))

	// Routes
	app.Get("/api/news", getNewsHandler)           // Get all news
	app.Post("/api/news", createNewsHandler)       // Create a new news item
	app.Get("/api/news/:id", getNewsByIDHandler)   // Get a specific news item by ID
	app.Put("/api/news/:id", updateNewsHandler)    // Update a specific news item
	app.Delete("/api/news/:id", deleteNewsHandler) // Delete a specific news item

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
