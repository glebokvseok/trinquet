package models

import "github.com/google/uuid"

type (
	SearchPreview PlayerPreview
)

type BasePlayerPreview struct {
	PlayerID uuid.UUID `json:"player_id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Avatar   *Avatar   `json:"avatar"`
}

type PlayerPreview struct {
	BasePlayerPreview
	Following     bool `json:"following"`
	FollowingBack bool `json:"following_back"`
}
