package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name            string    `json:"name"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	NumberOfLand    string    `json:"number_of_land"`
	GovNumberOfLand string    `json:"gov_number_of_land"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Question model
type Question struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ElectionID uuid.UUID `json:"election_id" gorm:"type:uuid"`
	Text       string    `json:"text"`
	Choices    []string  `json:"choices" gorm:"type:text[]"`
}

// Election model
type Election struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Questions   []Question `json:"questions" gorm:"foreignKey:ElectionID"`
}

type Vote struct {
	ID         uuid.UUID            `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	ElectionID uuid.UUID            `json:"election_id"`
	UserID     uuid.UUID            `json:"user_id"`
	Responses  map[uuid.UUID]string `gorm:"type:jsonb" json:"responses"`
	Timestamp  time.Time            `json:"timestamp"`
}

type TransitionApplication struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	User      string    `json:"user"`
	Time      string    `json:"time"` // Day-Month-Year format
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type News struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}
