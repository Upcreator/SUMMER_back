package controllers

import (
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
	"time"

	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/Upcreator/SUMMER_back/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateUser(c *fiber.Ctx) error {
	var payload *models.CreateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	now := time.Now()
	newUser := models.User{
		ID:        uuid.New(),
		Username:  payload.Username,
		Password:  utils.GeneratePassword(payload.Password),
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "(SQLSTATE 23505)") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Email already exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": newUser}})
}

func FindUsers(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var users []models.User
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&users)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(users), "users": users})
}

func UpdateUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var payload *models.UpdateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var user models.User
	result := initializers.DB.First(&user, "id = ?", userId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with this Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.Email != "" {
		updates["email"] = payload.Email
	}
	if payload.FullName != "" {
		updates["full_name"] = payload.FullName
	}

	updates["updated_at"] = time.Now()

	initializers.DB.Model(&user).Updates(updates)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

func FindUserById(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var user models.User
	result := initializers.DB.First(&user, "id = ?", userId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

func DeleteUser(c *fiber.Ctx) error {
	userId := c.Params("userId")

	result := initializers.DB.Delete(&models.User{}, "id = ?", userId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No user with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

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
		"sub": user.ID.String(),
		"exp": time.Now().Add(tokenLive).Unix(),
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
