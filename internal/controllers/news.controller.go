package controllers

import (
	"fmt"
	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateNews(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	content := c.FormValue("content")
	visibility := c.FormValue("visibility")

	payload := &models.CreateNewsSchema{
		Title:       title,
		Description: description,
		Content:     content,
		Visibility:  visibility,
	}
	if err := models.ValidateStruct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err})
	}

	news := models.News{
		ID:          uuid.New(),
		Title:       payload.Title,
		Description: payload.Description,
		Content:     payload.Content,
		Visibility:  payload.Visibility,
	}

	if fileHeader != nil {
		fileName := fmt.Sprintf("%s_%s", news.ID, fileHeader.Filename)
		destination := fmt.Sprintf("./uploads/%s", fileName)
		if err := c.SaveFile(fileHeader, destination); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err})
		}
		news.Preview = fmt.Sprintf("/uploads/%s", fileName)
	}
	result := initializers.DB.Create(&news)

	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exist"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "news": news})
}

func FindNews(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var news []models.News
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&news)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(news), "news": news})
}

func UpdateNews(c *fiber.Ctx) error {
	newsId := c.Params("newsId")
	var news models.News

	if err := initializers.DB.First(&news, "id = ?", newsId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Новость не найдена",
		})
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	content := c.FormValue("content")
	visibility := c.FormValue("visibility")

	payload := &models.CreateNewsSchema{
		Title:       title,
		Description: description,
		Content:     content,
		Visibility:  visibility,
	}
	if err := models.ValidateStruct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err,
		})
	}

	news.Title = payload.Title
	news.Description = payload.Description
	news.Content = payload.Content
	news.Visibility = payload.Visibility

	if fileHeader, err := c.FormFile("photo"); err == nil {
		fileName := fmt.Sprintf("%s_%s", news.ID, fileHeader.Filename)
		destination := fmt.Sprintf("./uploads/%s", fileName)

		if err := c.SaveFile(fileHeader, destination); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
		news.Preview = fmt.Sprintf("/uploads/%s", fileName)
	}

	if err := initializers.DB.Save(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// 8. Отдаём обновлённую новость
	return c.JSON(fiber.Map{
		"status": "success",
		"news":   news,
	})
}

func FindNewsById(c *fiber.Ctx) error {
	newsId := c.Params("newsId")

	var news models.News
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

	result := initializers.DB.Delete(&models.News{}, "id = ?", newsId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No news with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
