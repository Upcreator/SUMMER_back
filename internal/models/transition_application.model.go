package models

import (
	"github.com/google/uuid"
	"time"
)

type TransitionApplicationModel struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserId uuid.UUID `gorm:"not null" json:"user_id,omitempty"`
	Time   time.Time `gorm:"not null" json:"time"`
	Car    string    `gorm:"varchar(50);not null" json:"car"`
	Plate  string    `gorm:"varchar(10);not null" json:"plate"`
	User   User      `gorm:"foreignkey:user_id" json:"user,omitempty"`
}

type CreateTransitionApplicationSchema struct {
	UserId string `json:"user_id"`
	Time   string `json:"time"`
	Car    string `json:"car,omitempty"`
	Plate  string `json:"plate,omitempty"`
}

type UpdateTransitionApplicationSchema struct {
	User  string `json:"user,omitempty"`
	Time  string `json:"time,omitempty"`
	Car   string `json:"car,omitempty"`
	Plate string `json:"plate,omitempty"`
}
