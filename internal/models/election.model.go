package models

import (
	"time"

	"github.com/google/uuid"
)

type Election struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Questions   []Question `json:"questions" gorm:"foreignKey:ElectionID"`
}
