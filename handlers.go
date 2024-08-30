package main

import (
	"time"

	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler to create a new news item
func createNewsHandler(c *fiber.Ctx) error {
	var news models.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	news.ID = uuid.New()        // Generate UUID
	news.CreatedAt = time.Now() // Set current timestamp
	if err := createNews(db, news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(news)
}

// Handler to get all news items
func getNewsHandler(c *fiber.Ctx) error {
	newsItems, err := getNews()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch news"})
	}
	return c.Status(fiber.StatusOK).JSON(newsItems)
}

// Handler to get a specific news item by ID
func getNewsByIDHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	newsItem, err := getNewsByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to fetch news"})
	}
	return c.Status(fiber.StatusOK).JSON(newsItem)
}

// Handler to update a specific news item
func updateNewsHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	var updateData struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		Type        *string `json:"type,omitempty"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}
	if err := updateNews(id, updateData.Name, updateData.Description, updateData.Type); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to update news"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Handler to delete a specific news item
func deleteNewsHandler(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	if err := deleteNews(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to delete news"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
