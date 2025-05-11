package controllers

import (
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateTransitionApplication(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	userRole := c.Locals("userRole")

	var payload *models.CreateTransitionApplicationSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var transition models.TransitionApplicationModel

	if payload.Car != "" {
		transition.Car = payload.Car
	}
	if payload.Time != "" {
		t, err := time.Parse("2006-01-02T15:04", payload.Time)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "invalid time format: " + err.Error(),
			})
		}

		transition.Time = t
	}
	if payload.Plate != "" {
		transition.Plate = payload.Plate
	}

	if userRole != "admin" {
		userUuid, _ := uuid.Parse(userId)
		transition.UserId = userUuid
	} else {
		if payload.UserId != "" {
			userUuid, _ := uuid.Parse(payload.UserId)
			transition.UserId = userUuid
		}
	}

	result := initializers.DB.Create(&transition)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": transition})
}

func FindTransitionApplications(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	userRole := c.Locals("userRole").(string)

	var transition_applications []models.TransitionApplicationModel
	query := initializers.DB
	if userRole != "admin" {
		query = query.Where("user_id = ?", userId)
	}
	results := query.Preload("User").Find(&transition_applications)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(transition_applications), "transition_applications": transition_applications})
}

func UpdateTransitionApplication(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	userUuid, _ := uuid.Parse(userId)
	userRole := c.Locals("userRole").(string)
	transitionApplicationId := c.Params("transitionApplicationId")

	var payload *models.UpdateTransitionApplicationSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var transitionApplication models.TransitionApplicationModel
	result := initializers.DB.First(&transitionApplication, "id = ?", transitionApplicationId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No transition application with this Id exits"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if userRole != "admin" && transitionApplication.UserId != userUuid {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No transition application with this Id exits"})
	}
	updates := make(map[string]interface{})
	if payload.Car != "" {
		updates["car"] = payload.Car
	}
	if payload.Time != "" {
		t, err := time.Parse("2006-01-02T15:04", payload.Time)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "fail",
				"message": "invalid time format: " + err.Error(),
			})
		}

		updates["time"] = t
	}
	if payload.Plate != "" {
		updates["plate"] = payload.Plate
	}
	if payload.UserId != "" && userRole == "admin" {
		updates["user_id"] = payload.UserId
	}

	initializers.DB.Model(&transitionApplication).Updates(updates)

	var newModel models.TransitionApplicationModel
	initializers.DB.Preload("User").Find(&newModel, "id = ?", transitionApplicationId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"transition_application": newModel}})
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
	userRole := c.Locals("userRole").(string)
	userId := c.Locals("userId").(string)
	transitionApplicationId := c.Params("transitionApplicationId")
	fmt.Println(transitionApplicationId)
	transitionApplicationUuid, _ := uuid.Parse(transitionApplicationId)

	var result *gorm.DB

	if userRole == "admin" {
		result = initializers.DB.Delete(&models.TransitionApplicationModel{}, "id = ?", transitionApplicationUuid)
	} else {
		userUuid, _ := uuid.Parse(userId)
		result = initializers.DB.Delete(&models.TransitionApplicationModel{}, "id = ? AND user_id = ?", transitionApplicationUuid, userUuid)
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No transition application with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "deleted"})
}
