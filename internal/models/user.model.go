package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"mail"`
	FullName  string    `gorm:"unique;not null" json:"fullName"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:'user'" json:"role"`
	Avatar    string    `gorm:"type:text" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserSchema struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserSchema struct {
	Email    string `gorm:"unique;not null" json:"mail"`
	FullName string `gorm:"unique;not null" json:"fullName"`
}
