package controllers

import (
	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/Upcreator/SUMMER_back/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func LoginUser(c *fiber.Ctx) error {
	var payload models.CreateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	var user models.User
	result := initializers.DB.First(&user, "username = ?", payload.Username)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid email or password",
		})
	}

	if !utils.ComparePassword(user.Password, payload.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid email or password",
		})
	}

	tokenLive := time.Hour * 24 * 30

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID.String(),
		"role": user.Role,
		"exp":  time.Now().Add(tokenLive).Unix(),
	})

	tokenString, err := token.SignedString([]byte(initializers.AppConfig.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	c.Cookie(utils.NewCookie("token", tokenString, "", time.Now().Add(tokenLive)))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"token":  tokenString,
	})
}

func LogoutUser(c *fiber.Ctx) error {
	c.Cookie(utils.NewCookie("token", "", "", time.Now().Add(-time.Hour)))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

func GetUser(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	var user models.User
	result := initializers.DB.First(&user, "id = ?", userId)
	if err := result.Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}
