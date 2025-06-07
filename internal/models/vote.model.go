package models

import (
	"github.com/google/uuid"
	"time"
)

type VoteOption struct {
	Id     int       `gorm:"primaryKey" json:"id"`
	VoteId uuid.UUID `json:"vote_id"`
	Label  string    `json:"label"`
	Votes  int       `json:"votes"`
}

type Vote struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID    `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Timestamp   time.Time    `json:"timestamp"`
	Options     []VoteOption `json:"options" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Ended       bool         `json:"ended" gorm:"default:false"`
}
