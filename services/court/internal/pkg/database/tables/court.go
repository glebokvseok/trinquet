package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/database/pkg/psql/geo"
	"github.com/move-mates/trinquet/services/court/internal/pkg/models"
	"time"
)

const (
	CourtTableName = "court_service.court"
)

type Court struct {
	ID          uuid.UUID          `gorm:"column:id;type:uuid;primaryKey"`
	Name        string             `gorm:"column:name"`
	SportType   models.SportType   `gorm:"column:sport_type"`
	SettingType models.SettingType `gorm:"column:setting_type"`
	SurfaceType models.SurfaceType `gorm:"column:surface_type"`
	Price       int                `gorm:"column:price"`
	Rating      float32            `gorm:"column:rating"`
	Country     string             `gorm:"column:country"`
	City        string             `gorm:"column:city"`
	Address     string             `gorm:"column:address"`
	MapLink     string             `gorm:"column:map_link"`
	Location    geo.Point          `gorm:"column:location"`
	CreatedOn   time.Time          `gorm:"column:created_on"`
	ModifiedOn  time.Time          `gorm:"column:modified_on"`
}

func (*Court) TableName() string {
	return CourtTableName
}

func (court *Court) ToCourtModel() *models.Court {
	return &models.Court{
		ID:          court.ID,
		Name:        court.Name,
		SportType:   court.SportType,
		SettingType: court.SettingType,
		SurfaceType: court.SurfaceType,
		Price:       court.Price,
		Rating:      court.Rating,
		Country:     court.Country,
		City:        court.City,
		Address:     court.Address,
		MapLink:     court.MapLink,
		Location: models.Location{
			Latitude:  court.Location.Latitude,
			Longitude: court.Location.Longitude,
		},
	}
}
