package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID         uuid.UUID            `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ElectionID uuid.UUID            `json:"election_id"`
	UserID     uuid.UUID            `json:"user_id"`
	Responses  map[uuid.UUID]string `json:"responses" gorm:"type:bytea"`
	Timestamp  time.Time            `json:"timestamp"`
}
