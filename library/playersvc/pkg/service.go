package playersvc

import (
	"context"
	"generated/services/player/pkg/api"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/playersvc/pkg/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PlayerService interface {
	GetAllFollowers(ctx context.Context, playerID uuid.UUID) ([]uuid.UUID, error)
	GetPlayerPreviews(ctx context.Context, playerIDs []uuid.UUID) (map[uuid.UUID]models.PlayerPreview, error)
}

func providePlayerService(
	config playerServiceConfig,
) (PlayerService, error) {
	conn, err := grpc.NewClient(config.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := api.NewPlayerServiceClient(conn)

	return &playerService{
		client: client,
		config: config,
	}, nil
}

type playerService struct {
	client api.PlayerServiceClient
	config playerServiceConfig
}

func (service *playerService) GetAllFollowers(ctx context.Context, playerID uuid.UUID) ([]uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, service.config.RequestTimeout)
	defer cancel()

	req := &api.GetAllFollowersRequest{
		UserId: playerID.String(),
	}

	resp, err := service.client.GetAllFollowers(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, errors.Errorf("response is nil")
	}

	followerIDs := make([]uuid.UUID, 0)
	for _, follower := range resp.GetFollowers() {
		followerID, err := uuid.Parse(follower)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		followerIDs = append(followerIDs, followerID)
	}

	return followerIDs, nil
}

func (service *playerService) GetPlayerPreviews(ctx context.Context, playerIDs []uuid.UUID) (map[uuid.UUID]models.PlayerPreview, error) {
	ctx, cancel := context.WithTimeout(ctx, service.config.RequestTimeout)
	defer cancel()

	ids := make([]string, len(playerIDs))
	for i, playerID := range playerIDs {
		ids[i] = playerID.String()
	}

	req := &api.GetPlayerPreviewsRequest{
		PlayerIds: ids,
	}

	resp, err := service.client.GetPlayerPreviews(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp == nil {
		return nil, errors.Errorf("response is nil")
	}

	previews := make(map[uuid.UUID]models.PlayerPreview, len(playerIDs))
	for _, preview := range resp.GetPreviews() {
		playerID, err := uuid.Parse(preview.GetPlayerId())
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var avatar *models.Avatar
		if preview.GetAvatar() != nil {
			avatarID, err := uuid.Parse(preview.GetAvatar().GetId())
			if err != nil {
				return nil, errors.WithStack(err)
			}

			avatar = &models.Avatar{
				ID:       avatarID,
				MimeType: preview.GetAvatar().GetMimeType(),
				URL:      preview.GetAvatar().GetUrl(),
				Method:   preview.GetAvatar().GetMethod(),
			}
		}

		previews[playerID] = models.PlayerPreview{
			Name:    preview.GetName(),
			Surname: preview.GetSurname(),
			Avatar:  avatar,
		}
	}

	return previews, nil
}
