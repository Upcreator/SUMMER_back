package models

import "github.com/google/uuid"

type Question struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ElectionID uuid.UUID `json:"election_id" gorm:"type:uuid"`
	Text       string    `json:"text"`
	Choices    []string  `json:"choices" gorm:"type:text[]"`
}
