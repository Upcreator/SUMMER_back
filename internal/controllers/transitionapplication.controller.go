package controllers

import (
	"strconv"

	"time"

	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateTransitionApplication(c *fiber.Ctx) error {
	var payload *models.CreateTransitionApplicationSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newTransitionApplication := models.TransitionApplicationModel{
		User:     payload.User,
		Time:     payload.Time,
		Type:     payload.Type,
		CreateAt: now,
	}

	result := initializers.DB.Create(&newTransitionApplication)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"transition_application": newTransitionApplication}})
}

func FindTransitionApplications(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var transition_applications []models.TransitionApplicationModel
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&transition_applications)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(transition_applications), "transition_applications": transition_applications})
}

func UpdateTransitionApplication(c *fiber.Ctx) error {
	transitionApplicationId := c.Params("transitionApplicationId")

	var payload *models.UpdateTransitionApplicationSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var transition_application models.TransitionApplicationModel
	result := initializers.DB.First(&transition_application, "id = ?", transitionApplicationId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No transition application with this Id exits"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.User != "" {
		updates["user"] = payload.User
	}
	if payload.Time != "" {
		updates["time"] = payload.Time
	}
	if payload.Type != "" {
		updates["type"] = payload.Type
	}

	initializers.DB.Model(&transition_application).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"transition_application": transition_application}})
}

func FindTransitionApplicationById(c *fiber.Ctx) error {
	transitionApplicationId := c.Params("transitionApplicationId")

	var transition_application models.TransitionApplicationModel
	result := initializers.DB.First(&transition_application, "id = ?", transitionApplicationId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No transition_application with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"transition_application": transition_application}})
}

func DeleteTransitionApplication(c *fiber.Ctx) error {
	transitionApplicationId := c.Params("transitionApplicationId")

	result := initializers.DB.Delete(&models.TransitionApplicationModel{}, "id = ?", transitionApplicationId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No transition application with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
