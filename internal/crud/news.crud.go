package crud

import (
	"errors"

	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Create a new news item
func CreateNews(db *gorm.DB, news models.News) error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	result := db.Create(&news)
	return result.Error
}

// Get all news items
func GetNews(db *gorm.DB) ([]models.News, error) {
	var newsItems []models.News
	result := db.Find(&newsItems)
	return newsItems, result.Error
}

// Get a specific news item by id
func GetNewsByID(db *gorm.DB, id uuid.UUID) (models.News, error) {
	var news models.News
	result := db.First(&news, "id = ?", id)
	return news, result.Error
}

// Update a specific news item
func UpdateNews(db *gorm.DB, id uuid.UUID, name *string, description *string, newsType *string) error {
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
	result := db.Model(&models.News{}).Where("id = ?", id).Updates(updateData)
	return result.Error
}

// Delete a specific news item
func DeleteNews(db *gorm.DB, id uuid.UUID) error {
	result := db.Delete(&models.News{}, "id = ?", id)
	return result.Error
}
