package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Set up the DB connection
	setupDB()

	// Auto-migrate the schema
	db.AutoMigrate(&User{}, &Election{}, &Question{}, &Vote{}, &TransitionApplication{}, &News{})

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Routes
	e.GET("/api/news", getNewsHandler)           // Get all news
	e.POST("/api/news", createNewsHandler)       // Create a new news item
	e.GET("/api/news/:id", getNewsByIDHandler)   // Get a specific news item by ID
	e.PUT("/api/news/:id", updateNewsHandler)    // Update a specific news item
	e.DELETE("/api/news/:id", deleteNewsHandler) // Delete a specific news item

	// Start the server
	log.Fatal(e.Start(":8080"))
}
