package managers

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/common/pkg/async"
	"github.com/move-mates/trinquet/library/common/pkg/result"
	txmgr "github.com/move-mates/trinquet/library/database/pkg/psql/managers"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
	"github.com/move-mates/trinquet/services/player/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/player/internal/pkg/helpers/avatar"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"time"
)

const (
	playerPerSearchRequestCount = 20
)

type PlayerManager interface {
	GetSelf(ctx context.Context, selfID uuid.UUID) (*models.Player, []*models.RacquetProfile, error)
	GetOther(ctx context.Context, selfID uuid.UUID, playerID uuid.UUID) (*models.Player, []*models.RacquetProfile, error)
	SearchPlayers(ctx context.Context, selfID uuid.UUID, query string) ([]models.SearchPreview, error)
	CreatePlayer(ctx context.Context, playerID uuid.UUID, username string) error
	UpdatePlayer(ctx context.Context, playerID uuid.UUID, update models.PlayerUpdate) error
	GetBasePlayerPreviews(ctx context.Context, playerIDs []uuid.UUID) ([]models.BasePlayerPreview, error)
}

type playerManager struct {
	transactionManager       txmgr.TransactionManager
	playerRepository         repos.PlayerRepository
	avatarRepository         repos.AvatarRepository
	relationRepository       repos.RelationRepository
	racquetProfileRepository repos.RacquetProfileRepository
	linkGenerator            s3.LinkGenerator
}

func ProvidePlayerManager(
	transactionManager txmgr.TransactionManager,
	playerRepository repos.PlayerRepository,
	avatarRepository repos.AvatarRepository,
	relationRepository repos.RelationRepository,
	racquetProfileRepository repos.RacquetProfileRepository,
	linkGenerator s3.LinkGenerator,
) PlayerManager {
	return &playerManager{
		transactionManager:       transactionManager,
		playerRepository:         playerRepository,
		avatarRepository:         avatarRepository,
		relationRepository:       relationRepository,
		racquetProfileRepository: racquetProfileRepository,
		linkGenerator:            linkGenerator,
	}
}

func (mgr *playerManager) GetSelf(
	ctx context.Context,
	selfID uuid.UUID,
) (*models.Player, []*models.RacquetProfile, error) {
	return mgr.getPlayer(ctx, selfID, selfID) // TODO: переделать хак для получения собственной страницы
}

func (mgr *playerManager) GetOther(
	ctx context.Context,
	selfID uuid.UUID,
	playerID uuid.UUID,
) (*models.Player, []*models.RacquetProfile, error) {
	return mgr.getPlayer(ctx, selfID, playerID)
}

func (mgr *playerManager) getPlayer(
	ctx context.Context,
	selfID uuid.UUID,
	playerID uuid.UUID,
) (*models.Player, []*models.RacquetProfile, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	getFollowInfoCh := make(chan result.Result[models.FollowInfo])
	async.SafeGoRes(
		ctx,
		getFollowInfoCh,
		func(ctx context.Context) (models.FollowInfo, error) {
			return mgr.relationRepository.GetFollowInfo(ctx, selfID, playerID)
		},
	)

	getRacquetProfilesCh := make(chan result.Result[[]*models.RacquetProfile])
	async.SafeGoRes(
		ctx,
		getRacquetProfilesCh,
		func(ctx context.Context) ([]*models.RacquetProfile, error) {
			return mgr.racquetProfileRepository.GetProfiles(ctx, playerID)
		},
	)

	player, err := mgr.playerRepository.GetPlayer(ctx, playerID)
	if err != nil {
		return nil, nil, err
	}

	if player.Avatar != nil {
		player.Avatar.URL, player.Avatar.Method, err =
			mgr.linkGenerator.GenerateDownloadLink(ctx, avatar.Key(playerID, player.Avatar.ID))

		if err != nil {
			return nil, nil, err
		}
	}

	getFollowInfoRes := <-getFollowInfoCh
	if getFollowInfoRes.Err != nil {
		return nil, nil, getFollowInfoRes.Err
	}

	getProfilesRes := <-getRacquetProfilesCh
	if getProfilesRes.Err != nil {
		return nil, nil, getProfilesRes.Err
	}

	player.FollowInfo = getFollowInfoRes.Data

	return player, getProfilesRes.Data, nil
}

func (mgr *playerManager) SearchPlayers(ctx context.Context, selfID uuid.UUID, query string) ([]models.SearchPreview, error) {
	players, err := mgr.playerRepository.SearchPlayers(ctx, selfID, query, playerPerSearchRequestCount)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(players))
	for i := range players {
		ids[i] = players[i].PlayerID
		if players[i].Avatar != nil {
			players[i].Avatar.URL, players[i].Avatar.Method, err =
				mgr.linkGenerator.GenerateDownloadLink(ctx, avatar.Key(players[i].PlayerID, players[i].Avatar.ID))

			if err != nil {
				return nil, err
			}
		}
	}

	follows, err := mgr.relationRepository.GetFollows(ctx, selfID, ids)
	if err != nil {
		return nil, err
	}

	for i, follow := range follows {
		players[i].Following, players[i].FollowingBack = follow.Following, follow.FollowingBack
	}

	return players, nil
}

func (mgr *playerManager) CreatePlayer(ctx context.Context, playerID uuid.UUID, username string) (returnErr error) {
	ctx, err := mgr.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer mgr.transactionManager.FinalizeTransaction(ctx, &returnErr)

	if err = mgr.playerRepository.CreatePlayerTx(ctx, playerID, username); err != nil {
		return err
	}

	if err = mgr.racquetProfileRepository.CreateProfilesTx(ctx, playerID); err != nil {
		return err
	}

	if err = mgr.relationRepository.CreateUserNodeIfNotExists(ctx, playerID); err != nil {
		mgr.transactionManager.RollbackTransaction(ctx)

		return err
	}

	return mgr.transactionManager.CommitTransaction(ctx)
}

func (mgr *playerManager) UpdatePlayer(ctx context.Context, playerID uuid.UUID, update models.PlayerUpdate) (returnErr error) {
	ctx, err := mgr.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer mgr.transactionManager.FinalizeTransaction(ctx, &returnErr)

	currentTime := time.Now()

	if err = mgr.playerRepository.UpdatePlayerTx(ctx, playerID, update, currentTime); err != nil {
		return err
	}

	if err = mgr.avatarRepository.RemoveCurrentAvatarTx(ctx, playerID, currentTime); err != nil {
		return err
	}

	if update.AvatarID == nil {
		return mgr.transactionManager.CommitTransaction(ctx)
	}

	if err = mgr.avatarRepository.SetAvatarTx(ctx, playerID, *update.AvatarID, currentTime); err != nil {
		return err
	}

	return mgr.transactionManager.CommitTransaction(ctx)
}

func (mgr *playerManager) GetBasePlayerPreviews(ctx context.Context, playerIDs []uuid.UUID) ([]models.BasePlayerPreview, error) {
	previews, err := mgr.playerRepository.GetBasePlayerPreviews(ctx, playerIDs)
	if err != nil {
		return nil, err
	}

	for i := range previews {
		if previews[i].Avatar != nil {
			url, method, err := mgr.linkGenerator.GenerateDownloadLink(ctx, avatar.Key(previews[i].PlayerID, previews[i].Avatar.ID))
			if err != nil {
				return nil, err
			}

			previews[i].Avatar.URL, previews[i].Avatar.Method = url, method
		}
	}

	return previews, nil
}
