package models

import (
	"github.com/google/uuid"
)

type TransitionApplicationModel struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	User  string    `gorm:"not null" json:"user,omitempty"`
	Time  string    `gorm:"not null" json:"time"`
	Car   string    `gorm:"varchar(50);not null" json:"car"`
	Plate string    `gorm:"varchar(10);not null" json:"plate"`
}

type CreateTransitionApplicationSchema struct {
	User  string `json:"user" validate:"required"`
	Time  string `json:"time" validate:"required"`
	Car   string `json:"car,omitempty"`
	Plate string `json:"plate,omitempty"`
}

type UpdateTransitionApplicationSchema struct {
	User  string `json:"user,omitempty"`
	Time  string `json:"time,omitempty"`
	Car   string `json:"car,omitempty"`
	Plate string `json:"plate,omitempty"`
}
