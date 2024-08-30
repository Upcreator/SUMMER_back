package models

import (
	"time"

	"github.com/google/uuid"
)

type TransitionApplication struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	User      string    `json:"user"`
	Time      string    `json:"time"` // Day-Month-Year format
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
