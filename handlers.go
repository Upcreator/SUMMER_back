package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler to create a new news item
func createNewsHandler(c echo.Context) error {
	var news News
	if err := c.Bind(&news); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	news.ID = uuid.New()        // Generate a new UUID
	news.CreatedAt = time.Now() // Set the created_at timestamp
	if err := createNews(db, news); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, news)
}

// Handler to get all news items
func getNewsHandler(c echo.Context) error {
	newsItems, err := getNews()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to fetch news"})
	}
	return c.JSON(http.StatusOK, newsItems)
}

// Handler to get a specific news item by ID
func getNewsByIDHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	newsItem, err := getNewsByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to fetch news"})
	}
	return c.JSON(http.StatusOK, newsItem)
}

// Handler to update a specific news item
func updateNewsHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	var updateData struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		Type        *string `json:"type,omitempty"`
	}
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if err := updateNews(id, updateData.Name, updateData.Description, updateData.Type); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to update news"})
	}
	return c.NoContent(http.StatusNoContent)
}

// Handler to delete a specific news item
func deleteNewsHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	if err := deleteNews(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to delete news"})
	}
	return c.NoContent(http.StatusNoContent)
}
