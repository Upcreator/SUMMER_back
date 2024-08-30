package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID         uuid.UUID            `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	ElectionID uuid.UUID            `json:"election_id"`
	UserID     uuid.UUID            `json:"user_id"`
	Responses  map[uuid.UUID]string `gorm:"type:jsonb" json:"responses"`
	Timestamp  time.Time            `json:"timestamp"`
}
