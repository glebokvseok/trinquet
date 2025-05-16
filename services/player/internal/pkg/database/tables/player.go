package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/domain"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"time"
)

const (
	PlayerTableName = "player_service.player"
)

type Player struct {
	ID         uuid.UUID     `gorm:"column:id;type:uuid;primaryKey"`
	Username   string        `gorm:"column:username"`
	Name       string        `gorm:"column:name"`
	Surname    string        `gorm:"column:surname"`
	BirthDate  *time.Time    `gorm:"column:birth_date;type:date"`
	Gender     domain.Gender `gorm:"column:gender"`
	Height     int           `gorm:"column:height"`
	Country    string        `gorm:"column:country"`
	City       string        `gorm:"column:city"`
	CreatedOn  time.Time     `gorm:"column:created_on"`
	ModifiedOn time.Time     `gorm:"column:modified_on"`
}

func (*Player) TableName() string {
	return PlayerTableName
}

func (player *Player) ToPlayerModel() *models.Player {
	return &models.Player{
		Username:  player.Username,
		Name:      player.Name,
		Surname:   player.Surname,
		BirthDate: player.BirthDate,
		Gender:    player.Gender,
		Height:    player.Height,
		Country:   player.Country,
		City:      player.City,
	}
}
