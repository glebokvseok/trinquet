package models

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/httpsrv/dto"
)

type (
	FollowSortType string
	FollowPreview  PlayerPreview
)

const (
	DateTimeAsc  FollowSortType = "datetime_asc"
	DateTimeDesc FollowSortType = "datetime_desc"
)

type Follow struct {
	ID            uuid.UUID
	Following     bool
	FollowingBack bool
}

type FollowSort struct {
	Type   FollowSortType
	Cursor int64
}

type FollowSortResult struct {
	Follows []Follow
	Cursor  int64
}

type FollowSection struct {
	Follows        []FollowPreview `json:"follows"`
	Cursor         int64           `json:"cursor"`
	HasMoreFollows bool            `json:"has_more_follows"`
}

type FollowInfo struct {
	FollowersCount int64
	FollowingCount int64
	Following      bool
	FollowingBack  bool
}

func (info FollowInfo) ToFollowInfoDTO() dto.FollowInfo {
	return dto.FollowInfo{
		FollowersCount: info.FollowersCount,
		FollowingCount: info.FollowingCount,
		Following:      info.Following,
		FollowingBack:  info.FollowingBack,
	}
}
