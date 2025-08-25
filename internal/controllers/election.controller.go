package controllers

import (
	"strconv"
	"time"

	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateElection(c *fiber.Ctx) error {
	var payload models.Election

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	newElection := models.Election{
		ID:          uuid.New(),
		Name:        payload.Name,
		Description: payload.Description,
		Status:      payload.Status,
		StartDate:   payload.StartDate,
		EndDate:     payload.EndDate,
		Questions:   payload.Questions,
	}

	result := initializers.DB.Create(&newElection)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"election": newElection}})
}

func FindElections(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var elections []models.Election
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&elections)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(elections), "elections": elections})
}

func UpdateElection(c *fiber.Ctx) error {
	electionId := c.Params("electionId")

	var payload models.Election

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var election models.Election
	result := initializers.DB.First(&election, "id = ?", electionId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No election with this Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Name != "" {
		updates["name"] = payload.Name
	}
	if payload.Description != "" {
		updates["description"] = payload.Description
	}
	if payload.Status != "" {
		updates["status"] = payload.Status
	}
	if !payload.StartDate.IsZero() {
		updates["start_date"] = payload.StartDate
	}
	if !payload.EndDate.IsZero() {
		updates["end_date"] = payload.EndDate
	}

	updates["updated_at"] = time.Now()

	initializers.DB.Model(&election).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"election": election}})
}

func FindElectionById(c *fiber.Ctx) error {
	electionId := c.Params("electionId")

	var election models.Election
	result := initializers.DB.First(&election, "id = ?", electionId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No election with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"election": election}})
}

func DeleteElection(c *fiber.Ctx) error {
	electionId := c.Params("electionId")

	result := initializers.DB.Delete(&models.Election{}, " id = ?", electionId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No election with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
