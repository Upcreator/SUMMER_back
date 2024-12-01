package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/Upcreator/SUMMER_back/internal/initializers"
	"github.com/Upcreator/SUMMER_back/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

/*err = VerifyPassword(hashedPassword, password)
if err != nil {
	fmt.Println("Password does not match:", err)
} else {
	fmt.Println("Password matches!")
}
*/

func CreateUser(c *fiber.Ctx) error {
	var payload *models.CreateUpdateUserSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if err := payload.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Ошибка хэширования пароля: " + err.Error()})
	}

	now := time.Now()
	newUser := models.User{
		ID:                  uuid.New(),
		FirstName:           payload.FirstName,
		LastName:            payload.LastName,
		Surname:             payload.Surname,
		RegistrationAddress: payload.RegistrationAddress,
		ActualAddress:       payload.ActualAddress,
		NumberOfLand:        payload.NumberOfLand,
		GovNumberOfLand:     payload.GovNumberOfLand,
		Email:               payload.Email,
		Role:                payload.Role,
		Status:              payload.Status,
		Password:            payload.Password,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate key value violates unique") {
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

	var payload *models.CreateUpdateUserSchema

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

	if err := payload.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Ошибка хэширования пароля: " + err.Error()})
	}

	updates := make(map[string]interface{})
	if payload.FirstName != "" {
		updates["first_name"] = payload.FirstName
	}
	if payload.LastName != "" {
		updates["last_name"] = payload.LastName
	}
	if payload.Surname != "" {
		updates["surname"] = payload.Surname
	}
	if payload.RegistrationAddress != "" {
		updates["registration_address"] = payload.RegistrationAddress
	}
	if payload.ActualAddress != "" {
		updates["actual_address"] = payload.ActualAddress
	}
	if payload.NumberOfLand != "" {
		updates["number_of_land"] = payload.NumberOfLand
	}
	if payload.GovNumberOfLand != "" {
		updates["gov_number_of_land"] = payload.GovNumberOfLand
	}
	if payload.Email != "" {
		updates["email"] = payload.Email
	}
	if payload.Role != "" {
		updates["role"] = payload.Role
	}
	if payload.Status != "" {
		updates["status"] = payload.Status
	}
	if payload.Password != "" {
		updates["password"] = payload.Password
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
