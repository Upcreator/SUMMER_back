// model.go
package main

import (
	"time"

	"github.com/google/uuid"
)

// User roles
const (
	RoleChar      = "char"
	RoleSecretary = "secretary"
	RoleSecurity  = "security"
	RoleMember    = "member"
)

// User statuses
const (
	UserStatusNew      = "new"
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
)

// Election statuses
const (
	ElectionStatusUpcoming  = "upcoming"
	ElectionStatusOngoing   = "ongoing"
	ElectionStatusCompleted = "completed"
)

// TransitionApplication types
const (
	TransitionApplicationCar   = "car"
	TransitionApplicationTruck = "truck"
)

// User model
type User struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	NumberOfLand    string    `json:"number_of_land"`
	GovNumberOfLand string    `json:"gov_number_of_land"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
}

// Question model
type Question struct {
	ID      uuid.UUID `json:"id"`
	Text    string    `json:"text"`
	Choices []string  `json:"choices"`
}

// Election model
type Election struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Questions   []Question `json:"questions"`
}

// Vote model
type Vote struct {
	ID         uuid.UUID            `json:"id"`
	ElectionID uuid.UUID            `json:"election_id"`
	UserID     uuid.UUID            `json:"user_id"`
	Responses  map[uuid.UUID]string `json:"responses"`
	Timestamp  time.Time            `json:"timestamp"`
}

// TransitionApplication model
type TransitionApplication struct {
	ID        uuid.UUID `json:"id"`
	User      string    `json:"user"`
	Time      string    `json:"time"` // Day-Month-Year format
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// News model
type News struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}
