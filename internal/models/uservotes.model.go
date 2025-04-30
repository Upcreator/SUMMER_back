package models

import "github.com/google/uuid"

type UserVotes struct {
	UserId uuid.UUID `gorm:"type:uuid;default:null;" json:"user_id"`
	VoteId uuid.UUID `gorm:"type:uuid;default:null;" json:"vote_id"`
}
