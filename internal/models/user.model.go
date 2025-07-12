package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;default: null" json:"mail"`
	FullName  string    `gorm:"default: null" json:"fullName"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:'user'" json:"role"`
	Avatar    string    `gorm:"type:text" json:"avatar"`
	Region    string    `gorm:"default: null" json:"region"`
	Activated bool      `gorm:"default:false" json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserSchema struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Region   string `json:"region" validate:"required"`
}

type CreateUserAdminSchema struct {
	Username  string `json:"username" validate:"required"`
	FullName  string `json:"fullName"`
	Email     string `json:"mail"`
	Role      string `json:"role"`
	Password  string `json:"password" validate:"required"`
	Region    string `json:"region"`
	Activated bool   `json:"activated"`
}

type UpdateUserSchema struct {
	Username  string `json:"username"`
	Email     string `gorm:"unique;not null" json:"mail"`
	FullName  string `gorm:"unique;not null" json:"fullName"`
	Role      string `json:"role"`
	Region    string `json:"region"`
	Activated bool   `json:"activated"`
}
