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

func CreateVote(c *fiber.Ctx) error {
	var payload models.Vote

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	payload.Timestamp = time.Now()

	result := initializers.DB.Create(&payload)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"vote": payload}})
}

func FindVotes(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var votes []models.Vote
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&votes)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(votes), "votes": votes})
}

func UpdateVote(c *fiber.Ctx) error {
	voteId := c.Params("voteId")

	var payload *models.Vote

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var vote models.Vote
	result := initializers.DB.First(&vote, "id = ?", voteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No vote with this Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.ElectionID != uuid.Nil {
		updates["election_id"] = payload.ElectionID
	}
	if payload.UserID != uuid.Nil {
		updates["user_id"] = payload.UserID
	}
	if payload.Responses != nil {
		updates["responses"] = payload.Responses
	}

	initializers.DB.Model(&vote).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"vote": vote}})
}

func FindVoteById(c *fiber.Ctx) error {
	voteId := c.Params("voteId")

	var vote models.Vote
	result := initializers.DB.First(&vote, "id = ?", voteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No vote with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"vote": vote}})
}

func DeleteVote(c *fiber.Ctx) error {
	voteId := c.Params("voteId")

	result := initializers.DB.Delete(&models.Vote{}, "id = ?", voteId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No vote with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
