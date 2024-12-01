package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type NewsModel struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title      string    `gorm:"varchar(100);uniqueIndex;not null" json:"title,omitempty"`
	Content    string    `gorm:"not null" json:"content,omitempty"`
	Visibility bool      `gorm:"default:false;not null" json:"visibility"`
	Type       string    `gorm:"not null" json:"type,omitempty"`
	CreatedAt  time.Time `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt  time.Time `gorm:"not null" json:"updatedAt,omitempty"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type CreateNewsSchema struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Visibility bool   `json:"visibility,omitempty"`
}

type UpdateNewsSchema struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	Visibility *bool  `json:"visibility,omitempty"`
}
