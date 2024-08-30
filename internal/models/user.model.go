package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid; default:uuid_generate_v4()" json:"id"`
	Name            string    `json:"name"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	NumberOfLand    string    `json:"number_of_land"`
	GovNumberOfLand string    `json:"gov_number_of_land"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
