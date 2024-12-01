package models

import (
	"time"

	"github.com/google/uuid"
)

type Election struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string     `json:"name" gorm:"type:varchar(255);not null" validate:"required,cyrillic,min=1,max=255"`
	Description string     `json:"description" gorm:"type:varchar(255)" validate:"omitempty,cyrillic,min=0,max=255"`
	Status      string     `json:"status" gorm:"type:varchar(20);default:'Неактивный'" validate:"required,oneof=Неактивный Активный Завершенный"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Questions   []Question `json:"questions" gorm:"foreignKey:ElectionID"`
}
