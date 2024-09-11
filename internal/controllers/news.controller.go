package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateNews(c *fiber.Ctx) error {
	var payload *models.CreateNewsSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newNews := models.NewsModel{
		Title:      payload.Title,
		Content:    payload.Content,
		Visibility: payload.Visibility,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result := initializers.DB.Create(&newNews)

	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exist"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"news": newNews}})
}

func FindNews(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var news []models.NewsModel
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&news)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(news), "news": news})
}

func UpdateNews(c *fiber.Ctx) error {
	newsId := c.Params("newsId")

	var payload *models.UpdateNewsSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var news models.NewsModel
	result := initializers.DB.First(&news, "id = ?", newsId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No news with this Id exits"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Content != "" {
		updates["content"] = payload.Content
	}
	if payload.Visibility != nil {
		updates["visibility"] = payload.Visibility
	}

	updates["updated_at"] = time.Now()

	initializers.DB.Model(&news).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"news": news}})
}

func FindNewsById(c *fiber.Ctx) error {
	newsId := c.Params("newsId")

	var news models.NewsModel
	result := initializers.DB.First(&news, "id = ?", newsId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No news with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"news": news}})
}

func DeleteNews(c *fiber.Ctx) error {
	newsId := c.Params("newsId")

	result := initializers.DB.Delete(&models.NewsModel{}, "id = ?", newsId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No news with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
