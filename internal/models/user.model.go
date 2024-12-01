package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName           string    `json:"first_name" validate:"required"`
	LastName            string    `json:"last_name" validate:"required"`
	Surname             string    `json:"surname"`
	RegistrationAddress string    `json:"registration_address,omitempty"`
	ActualAddress       string    `json:"actual_address,omitempty"`
	NumberOfLand        string    `json:"number_of_land,omitempty"`
	GovNumberOfLand     string    `json:"gov_number_of_land,omitempty"`
	Email               string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Role                string    `json:"role" validate:"required"`
	Status              string    `json:"status" validate:"required"`
	Password            string    `json:"password" validate:"required"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
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
