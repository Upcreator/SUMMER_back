package controllers

import (
	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type IncomingVote struct {
	Title   string              `json:"title"`
	Options []models.VoteOption `json:"options"`
}

func CreateVote(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	var incoming IncomingVote
	if err := c.BodyParser(&incoming); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	uid, _ := uuid.Parse(userId)

	vote := models.Vote{
		UserID:    uid,
		Title:     incoming.Title,
		Timestamp: time.Now(),
	}
	result := initializers.DB.Create(&vote)

	for i := range incoming.Options {
		incoming.Options[i].VoteId = vote.ID
	}

	initializers.DB.Create(incoming.Options)

	if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"vote": vote}})
}

func FindVotes(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var votes []models.Vote
	results := initializers.DB.Limit(intLimit).Offset(offset).Order("timestamp desc").Model(&models.Vote{}).Preload("Options").Find(&votes)
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
	//if payload.ElectionID != uuid.Nil {
	//	updates["election_id"] = payload.ElectionID
	//}
	//if payload.UserID != uuid.Nil {
	//	updates["user_id"] = payload.UserID
	//}
	//if payload.Responses != nil {
	//	updates["responses"] = payload.Responses
	//}

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

type VoteSchema struct {
	Id int `json:"id"`
}

func UserVote(c *fiber.Ctx) error {
	userId, _ := uuid.Parse(c.Locals("userId").(string))
	voteId := c.Params("voteId")

	var payload VoteSchema
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var userVote models.UserVotes
	result := initializers.DB.Where("user_id = ?", userId.String()).Where("vote_id = ?", voteId).First(&userVote)

	if err := result.Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Already voted"})
	}

	var voteOption models.VoteOption

	result = initializers.DB.Where("vote_id = ?", voteId).First(&voteOption, "id = ?", payload.Id)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No vote with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	initializers.DB.Exec("UPDATE vote_options SET votes = votes + 1 WHERE id = ?", voteOption.Id)
	initializers.DB.Create(&models.UserVotes{
		UserId: userId,
		VoteId: voteOption.VoteId,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "voted successfully"})
}
