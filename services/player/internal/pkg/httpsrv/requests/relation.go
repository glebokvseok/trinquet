package requests

import (
	"github.com/google/uuid"
	"time"
)

type FollowPlayerRequest struct {
	PlayerId   uuid.UUID `json:"player_id" validate:"required"`
	FollowedOn time.Time `json:"followed_on" validate:"required"`
}

type UnfollowPlayerRequest struct {
	PlayerId     uuid.UUID `json:"player_id" validate:"required"`
	UnfollowedOn time.Time `json:"unfollowed_on" validate:"required"`
}
