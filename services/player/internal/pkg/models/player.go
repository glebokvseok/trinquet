package models

import (
	"github.com/google/uuid"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/dto"
	"time"
)

type (
	Avatar s3.Object
)

type Player struct {
	Username   string
	Name       string
	Surname    string
	BirthDate  *time.Time
	Gender     domain.Gender
	Height     int
	Country    string
	City       string
	MatchCount int
	FollowInfo FollowInfo
	Avatar     *Avatar
}

type PlayerUpdate struct {
	Name      string        `json:"name" validate:"omitempty,max=32"`
	Surname   string        `json:"surname" validate:"omitempty,max=32"`
	BirthDate *time.Time    `json:"birth_date" validate:"omitempty"`
	Gender    domain.Gender `json:"gender" validate:"omitempty,min=0,max=2"`
	Height    int           `json:"height" validate:"omitempty"`
	Country   string        `json:"country" validate:"omitempty,max=32"`
	City      string        `json:"city" validate:"omitempty,max=32"`
	AvatarID  *uuid.UUID    `json:"avatar_id" validate:"omitempty"`
}

func (player *Player) ToPlayerDTO() *dto.Player {
	return &dto.Player{
		Username:   player.Username,
		Name:       player.Name,
		Surname:    player.Surname,
		BirthDate:  player.BirthDate,
		Gender:     player.Gender,
		Height:     player.Height,
		Country:    player.Country,
		City:       player.City,
		MatchCount: player.MatchCount,
		FollowInfo: player.FollowInfo.ToFollowInfoDTO(),
		Avatar:     (*s3.Object)(player.Avatar),
	}
}
