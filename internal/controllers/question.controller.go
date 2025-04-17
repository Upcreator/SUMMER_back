package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateQuestion(c *fiber.Ctx) error {
	var payload *models.CreateQuestionSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	newQuestion := models.Question{
		Text:    payload.Text,
		Choices: pq.StringArray(payload.Choices),
	}

	fmt.Printf("Choices: %+v\n", newQuestion.Choices)
	result := initializers.DB.Create(&newQuestion)

	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exist"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"question": newQuestion}})
}

func FindQuestion(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var questions []models.Question
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&questions)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(questions), "questions": questions})
}

func UpdateQuestion(c *fiber.Ctx) error {
	questionId := c.Params("questionId")

	var payload *models.UpdateQuestionSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var question models.Question
	result := initializers.DB.First(&question, "id = ?", questionId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No question with this Id exits"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Text != "" {
		updates["text"] = payload.Text
	}
	if len(payload.Choices) != 0 {
		updates["choices"] = payload.Choices
	}

	initializers.DB.Model(&question).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"question": question}})
}

func FindQuestionById(c *fiber.Ctx) error {
	questionId := c.Params("questionId")

	var question models.Question
	result := initializers.DB.First(&question, "id = ?", questionId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No question with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"question": question}})
}

func DeleteQuestion(c *fiber.Ctx) error {
	questionId := c.Params("questionId")

	result := initializers.DB.Delete(&models.Question{}, "id = ?", questionId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No question with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
