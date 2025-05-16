package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/collections/slice"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"github.com/move-mates/trinquet/services/court/internal/pkg/database/tables"
	apierrors "github.com/move-mates/trinquet/services/court/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/court/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	similarityThreshold = 0.3
)

type CourtRepository interface {
	SearchCourt(ctx context.Context, filter models.SearchFilter) ([]*models.Court, error)
	GetCourtMedia(ctx context.Context, courtID uuid.UUID) ([]models.Media, error)
}

type courtRepository struct {
	db *gorm.DB
}

func ProvideCourtRepository(db *gorm.DB) CourtRepository {
	return &courtRepository{
		db: db,
	}
}

func (repo *courtRepository) SearchCourt(
	ctx context.Context,
	filter models.SearchFilter,
) ([]*models.Court, error) {
	type aggregatedCourt struct {
		tables.Court
		*tables.Media
	}

	var aggCourts []aggregatedCourt
	if err := buildSearchRequest(repo.db, filter).
		WithContext(ctx).
		Scan(&aggCourts).
		Error; err != nil {
		return nil, errors.WithStack(err)
	}

	courts := make([]*models.Court, len(aggCourts))
	for i, aggCourt := range aggCourts {
		courts[i] = aggCourt.ToCourtModel()
		if aggCourt.Media != nil {
			courts[i].Preview = aggCourt.ToMediaModel()
		}
	}

	return courts, nil
}

func (repo *courtRepository) GetCourtMedia(
	ctx context.Context,
	courtID uuid.UUID,
) ([]models.Media, error) {
	ok, err := repo.checkIfCourtExists(ctx, courtID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, apierrors.NewCourtNotFoundError(courtID)
	}

	var medias []tables.Media
	if err = repo.db.
		Table(tables.CourtMediaTableName).
		Where("court_id = ?", courtID).
		Order("is_preview DESC").
		Order("created_on ASC").
		Find(&medias).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	mediaModels := make([]models.Media, len(medias))
	for i, media := range medias {
		mediaModels[i].ID, mediaModels[i].MimeType = media.MediaID, media.MimeType
	}

	return mediaModels, nil
}

func (repo *courtRepository) checkIfCourtExists(
	ctx context.Context,
	id uuid.UUID,
) (ok bool, returnErr error) {
	var court = new(tables.Court)
	if err := repo.db.
		WithContext(ctx).
		Take(court, "id = ?", id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errors.WithStack(err)
	}

	return true, nil
}

func buildSearchRequest(
	db *gorm.DB,
	filter models.SearchFilter,
) *gorm.DB {
	db = db.
		Table("court_service.court AS c").
		Joins("left join court_service.court_media cm ON c.id = cm.court_id and cm.is_preview = true")

	if extensions.IsNotEmpty(filter.Name) {
		db = db.Where(
			"similarity(c.name, @name) > @threshold",
			map[string]interface{}{
				"name":      filter.Name,
				"threshold": similarityThreshold,
			},
		)
	}
	if slice.HasElements(filter.SportTypes) {
		db = db.Where("c.sport_type IN ?", filter.SportTypes)
	}
	if slice.HasElements(filter.SettingTypes) {
		db = db.Where("c.setting_type IN ?", filter.SettingTypes)
	}
	if slice.HasElements(filter.SurfaceTypes) {
		db = db.Where("c.surface_type IN ?", filter.SurfaceTypes)
	}
	if extensions.IsNotZero(filter.MinPrice) {
		db = db.Where("c.price >= ?", filter.MinPrice)
	}
	if extensions.IsNotZero(filter.MaxPrice) {
		db = db.Where("c.price <= ?", filter.MaxPrice)
	}
	if filter.MinRating > 0 {
		db = db.Where("c.rating >= ?", filter.MinRating)
	}
	if extensions.IsNotEmpty(filter.Country) {
		db = db.Where("c.country = ?", filter.Country)
	}
	if extensions.IsNotEmpty(filter.City) {
		db = db.Where("c.city = ?", filter.City)
	}

	return db
}
