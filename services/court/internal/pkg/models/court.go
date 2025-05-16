package models

import (
	"github.com/google/uuid"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
)

type (
	SportType   int
	SettingType int
	SurfaceType int

	Media s3.Object
)

const (
	UnknownSportType SportType = iota
	Tennis
	Padel
)

const (
	UnknownSettingType SettingType = iota
	Indoor
	Outdoor
)

const (
	UnknownSurfaceType SurfaceType = iota
	Clay
	Glass
	Hard
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Court struct {
	ID          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	SportType   SportType   `json:"sport_type"`
	SettingType SettingType `json:"setting_type"`
	SurfaceType SurfaceType `json:"surface_type"`
	Price       int         `json:"price"`
	Rating      float32     `json:"rating"`
	Country     string      `json:"country"`
	City        string      `json:"city"`
	Address     string      `json:"address"`
	MapLink     string      `json:"map_link"`
	Location    Location    `json:"location"`
	Preview     *Media      `json:"preview"`
}

type SearchFilter struct {
	Name         string        `json:"name" validate:"omitempty"`
	SportTypes   []SportType   `json:"sport_types" validate:"omitempty,max=5"`
	SettingTypes []SettingType `json:"setting_types" validate:"omitempty,max=5"`
	SurfaceTypes []SurfaceType `json:"surface_types" validate:"omitempty,max=10"`
	MinPrice     int           `json:"min_price" validate:"omitempty,min=0,max=100000"`
	MaxPrice     int           `json:"max_price" validate:"omitempty,min=0,max=100000"`
	MinRating    float32       `json:"min_rating" validate:"omitempty,min=0,max=5"`
	Country      string        `json:"country" validate:"omitempty,max=32"`
	City         string        `json:"city" validate:"omitempty,max=32"`
}

type SearchResult struct {
	Courts []*Court `json:"courts"`
}

type CourtMedia struct {
	Medias []Media `json:"medias"`
}
