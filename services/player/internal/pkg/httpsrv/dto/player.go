package dto

import (
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"time"
)

type Player struct {
	Username   string        `json:"username"`
	Name       string        `json:"name"`
	Surname    string        `json:"surname"`
	BirthDate  *time.Time    `json:"birth_date"`
	Gender     domain.Gender `json:"gender"`
	Height     int           `json:"height"`
	Country    string        `json:"country"`
	City       string        `json:"city"`
	MatchCount int           `json:"match_count"`
	FollowInfo FollowInfo    `json:"follow_info"`
	Avatar     *s3.Object    `json:"avatar"`
}
