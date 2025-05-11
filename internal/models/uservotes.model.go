package models

import "github.com/google/uuid"

type UserVotes struct {
	UserId  uuid.UUID `gorm:"type:uuid;default:null;" json:"user_id"`
	VoteId  uuid.UUID `gorm:"type:uuid;default:null;" json:"vote_id"`
	VotedTo string    `gorm:"default:null;" json:"voted_to"`
	User    User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}
