package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Question struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Text       string         `json:"text"`
	Choices    pq.StringArray `json:"choices" gorm:"type:text[]"`
	ElectionID uuid.UUID      `json:"election_id" gorm:"type:uuid"`
}

type CreateQuestionSchema struct {
	Text    string         `json:"text" validate:"required"`
	Choices pq.StringArray `json:"choices"`
}

type UpdateQuestionSchema struct {
	Text    string         `json:"text" validate:"required"`
	Choices pq.StringArray `json:"choices"`
}
