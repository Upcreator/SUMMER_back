package main

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Create a new news item
func createNews(db *gorm.DB, news News) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	result := db.Create(&news)
	return result.Error
}

// Get all news items
func getNews() ([]News, error) {
	var newsItems []News
	result := db.Find(&newsItems)
	return newsItems, result.Error
}

// Get a specific news item by id
func getNewsByID(id uuid.UUID) (News, error) {
	var news News
	result := db.First(&news, "id = ?", id)
	return news, result.Error
}

// Update a specific news item
func updateNews(id uuid.UUID, name *string, description *string, newsType *string) error {
	updateData := map[string]interface{}{}
	if name != nil {
		updateData["name"] = *name
	}
	if description != nil {
		updateData["description"] = *description
	}
	if newsType != nil {
		updateData["type"] = *newsType
	}
	result := db.Model(&News{}).Where("id = ?", id).Updates(updateData)
	return result.Error
}

// Delete a specific news item
func deleteNews(id uuid.UUID) error {
	result := db.Delete(&News{}, "id = ?", id)
	return result.Error
}
