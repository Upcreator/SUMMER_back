package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName           string    `json:"first_name" validate:"required"`                              // Имя
	LastName            string    `json:"last_name" validate:"required"`                               // Фамилия
	Surname             string    `json:"surname"`                                                     // Отчество
	RegistrationAddress string    `json:"registration_address,omitempty"`                              // Адрес регистрации
	ActualAddress       string    `json:"actual_address,omitempty"`                                    // Фактический адрес
	NumberOfLand        string    `json:"number_of_land,omitempty"`                                    // Номер участка
	GovNumberOfLand     string    `json:"gov_number_of_land,omitempty"`                                // Кадастровый номер участка
	Email               string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"` // Электронная почта
	Role                string    `json:"role" validate:"required"`                                    // Роль
	Status              string    `json:"status" validate:"required"`                                  // Статус
	Password            string    `json:"password" validate:"required"`                                // Пароль
	CreatedAt           time.Time `json:"created_at"`                                                  // Время создания
	UpdatedAt           time.Time `json:"updated_at"`                                                  // Время обновления
}

type CreateUpdateUserSchema struct {
	FirstName           string `json:"first_name" validate:"required"`
	LastName            string `json:"last_name" validate:"required"`
	Surname             string `json:"surname"`
	RegistrationAddress string `json:"registration_address,omitempty"`
	ActualAddress       string `json:"actual_address,omitempty"`
	NumberOfLand        string `json:"number_of_land,omitempty"`
	GovNumberOfLand     string `json:"gov_number_of_land,omitempty"`
	Email               string `json:"email" validate:"required,email"`
	Role                string `json:"role" validate:"required"`
	Status              string `json:"status" validate:"required"`
	Password            string `json:"password" validate:"required"`
}

// Метод для хэширования пароля
func (user *CreateUpdateUserSchema) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
