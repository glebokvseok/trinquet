package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/move-mates/trinquet/library/common/pkg/collections/slice"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"github.com/move-mates/trinquet/library/database/pkg/psql"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/tables"
	apierrors "github.com/move-mates/trinquet/services/player/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

const (
	similarityThreshold = 0.3
)

const (
	getPlayerQuery = `
		SELECT * FROM player_service.player p
		LEFT JOIN player_service.avatar a
		ON p.id = a.player_id AND a.is_used = TRUE
		WHERE p.id = ?;
	`

	searchPlayerQuery = `
		SELECT p.id, p.name, p.surname, a.avatar_id, a.mime_type,
			GREATEST(
				similarity(p.username, @query),
				similarity(p.name, @query),
				similarity(p.surname, @query)
			) AS similarity
		FROM player_service.player p
		LEFT JOIN player_service.avatar a
		ON p.id = a.player_id AND a.is_used = TRUE
		WHERE GREATEST(
			similarity(p.username, @query),
			similarity(p.name, @query),
			similarity(p.surname, @query)
		) > @threshold
		AND p.id != @id
		ORDER BY similarity DESC
		LIMIT @limit
	`

	getBasePlayerPreviewsQuery = `
		SELECT p.id, p.name, p.surname, a.avatar_id, a.mime_type
		FROM player_service.player p
		LEFT JOIN player_service.avatar a 
		ON p.id = a.player_id AND a.is_used = TRUE
		WHERE p.id IN ?;
    `
)

type PlayerRepository interface {
	GetPlayer(ctx context.Context, id uuid.UUID) (*models.Player, error)
	CreatePlayerTx(ctx context.Context, id uuid.UUID, username string) error
	UpdatePlayerTx(ctx context.Context, id uuid.UUID, update models.PlayerUpdate, currentTime time.Time) error
	SearchPlayers(ctx context.Context, selfID uuid.UUID, query string, playerCount int) ([]models.SearchPreview, error)
	CheckIfPlayerExists(ctx context.Context, id uuid.UUID) (bool, error)
	GetBasePlayerPreviews(ctx context.Context, ids []uuid.UUID) ([]models.BasePlayerPreview, error)
}

type playerRepository struct {
	db *gorm.DB
}

func ProvidePlayerRepository(
	db *gorm.DB,
) PlayerRepository {
	return &playerRepository{
		db: db,
	}
}

func (repo *playerRepository) GetPlayer(
	ctx context.Context,
	id uuid.UUID,
) (*models.Player, error) {
	type aggregatedPlayer struct {
		tables.Player
		*tables.Avatar
	}

	aggPlayer := new(aggregatedPlayer)
	if err := repo.db.
		WithContext(ctx).
		Raw(getPlayerQuery, id).
		Take(&aggPlayer).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewPlayerNotFoundError(id)
		}

		return nil, errors.WithStack(err)
	}

	player := aggPlayer.ToPlayerModel()
	if aggPlayer.Avatar != nil {
		player.Avatar = aggPlayer.ToAvatarModel()
	}

	return player, nil
}

func (repo *playerRepository) CreatePlayerTx(
	ctx context.Context,
	id uuid.UUID,
	username string,
) error {
	tx, err := getTransaction(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	player := &tables.Player{
		ID:       id,
		Username: username,
	}

	if err = tx.
		WithContext(ctx).
		Select("id", "username").
		Create(player).
		Error; err != nil {
		if pgErr := (*pgconn.PgError)(nil); errors.As(err, &pgErr) && pgErr.Code == psql.UniqueConstraintViolationErrorCode {
			switch pgErr.ConstraintName {
			case "player_pkey":
				return apierrors.NewPlayerAlreadyExistsError(id)
			case "player_username_key":
				return apierrors.NewUsernameIsTakenError(username)
			default:
				return errors.WithStack(err)
			}
		}

		return errors.WithStack(err)
	}

	return nil
}

func (repo *playerRepository) UpdatePlayerTx(
	ctx context.Context,
	id uuid.UUID,
	update models.PlayerUpdate,
	currentTime time.Time,
) error {
	tx, err := getTransaction(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	res := tx.
		WithContext(ctx).
		Table(tables.PlayerTableName).
		Where("id = ?", id).
		Updates(map[string]any{
			"name":        extensions.ToNullIfEmpty(update.Name),
			"surname":     extensions.ToNullIfEmpty(update.Surname),
			"birth_date":  update.BirthDate,
			"gender":      update.Gender,
			"height":      extensions.ToNullIfZero(update.Height),
			"country":     extensions.ToNullIfEmpty(update.Country),
			"city":        extensions.ToNullIfEmpty(update.City),
			"modified_on": currentTime,
		})

	if res.Error != nil {
		return errors.WithStack(res.Error)
	}

	if res.RowsAffected != 1 {
		return apierrors.NewPlayerNotFoundError(id)
	}

	return nil
}

func (repo *playerRepository) SearchPlayers(
	ctx context.Context,
	selfID uuid.UUID,
	query string,
	playerCount int,
) ([]models.SearchPreview, error) {
	params := map[string]interface{}{
		"id":        selfID,
		"query":     query,
		"threshold": similarityThreshold,
		"limit":     playerCount,
	}

	var previews []preview
	if err := repo.db.
		WithContext(ctx).
		Raw(searchPlayerQuery, params).
		Scan(&previews).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	mapFunc := func(preview preview) models.SearchPreview {
		var avatar *models.Avatar
		if preview.AvatarID != nil {
			avatar = &models.Avatar{
				ID:       *preview.AvatarID,
				MimeType: preview.MimeType,
			}
		}

		return models.SearchPreview{
			BasePlayerPreview: models.BasePlayerPreview{
				PlayerID: preview.PlayerID,
				Name:     preview.Name,
				Surname:  preview.Surname,
				Avatar:   avatar,
			},
		}
	}

	return slice.Map(previews, mapFunc), nil
}

func (repo *playerRepository) CheckIfPlayerExists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	var player = new(tables.Player)
	if err := repo.db.
		WithContext(ctx).
		Take(player, "id = ?", id).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errors.WithStack(err)
	}

	return true, nil
}

func (repo *playerRepository) GetBasePlayerPreviews(
	ctx context.Context,
	ids []uuid.UUID,
) ([]models.BasePlayerPreview, error) {
	var previews []*preview
	if err := repo.db.
		WithContext(ctx).
		Raw(getBasePlayerPreviewsQuery, ids).
		Scan(&previews).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	basePreviewsMap := make(map[uuid.UUID]models.BasePlayerPreview, len(previews))
	for _, preview := range previews {
		var avatar *models.Avatar
		if preview.AvatarID != nil {
			avatar = &models.Avatar{
				ID:       *preview.AvatarID,
				MimeType: preview.MimeType,
			}
		}

		basePreviewsMap[preview.PlayerID] = models.BasePlayerPreview{
			PlayerID: preview.PlayerID,
			Name:     preview.Name,
			Surname:  preview.Surname,
			Avatar:   avatar,
		}
	}

	basePreviews := make([]models.BasePlayerPreview, len(ids))
	for i, id := range ids {
		basePreviews[i] = basePreviewsMap[id]
	}

	return basePreviews, nil
}

type preview struct {
	PlayerID uuid.UUID  `gorm:"column:id;type:uuid"`
	Name     string     `gorm:"column:name"`
	Surname  string     `gorm:"column:surname"`
	AvatarID *uuid.UUID `gorm:"column:avatar_id;type:uuid"`
	MimeType string     `gorm:"column:mime_type"`
}
