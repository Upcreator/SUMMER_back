package models

import (
	"time"

	"github.com/google/uuid"
)

type TransitionApplicationModel struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	User     string    `gorm:"not null" json:"user,omitempty"`
	Time     string    `gorm:"not null" json:"time"`
	Type     string    `gorm:"varchar(50);not null" json:"type"`
	CreateAt time.Time `gorm:"not null" json:"createAt"`
}

type CreateTransitionApplicationSchema struct {
	User string `json:"user" validate:"required"`
	Time string `json:"time" validate:"required"`
	Type string `json:"type,omitempty"`
}

type UpdateTransitionApplicationSchema struct {
	User string `json:"user,omitempty"`
	Time string `json:"time,omitempty"`
	Type string `json:"type,omitempty"`
}
