package grpcsrv

import (
	"context"
	"generated/services/player/pkg/api"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	api.UnimplementedPlayerServiceServer
	playerManager   managers.PlayerManager
	relationManager managers.RelationManager
	logger          *logrus.Logger
}

func (service *service) GetAllFollowers(ctx context.Context, req *api.GetAllFollowersRequest) (*api.GetAllFollowersResponse, error) {
	// TODO: добавить глобальные обработчики паник и ошибок + логирование
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	followers, err := service.relationManager.GetAllFollowerIDs(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetAllFollowersResponse{
		Followers: followers,
	}, nil
}

func (service *service) GetPlayerPreviews(ctx context.Context, req *api.GetPlayerPreviewsRequest) (*api.GetPlayerPreviewsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	playerIDs := make([]uuid.UUID, len(req.GetPlayerIds()))
	for i, id := range req.GetPlayerIds() {
		playerID, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		playerIDs[i] = playerID
	}

	basePreviews, err := service.playerManager.GetBasePlayerPreviews(ctx, playerIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	previews := make([]*api.PlayerPreview, len(basePreviews))
	for i, basePreview := range basePreviews {
		previews[i] = &api.PlayerPreview{
			PlayerId: basePreview.PlayerID.String(),
			Name:     basePreview.Name,
			Surname:  basePreview.Surname,
		}

		if basePreview.Avatar != nil {
			previews[i].Avatar = &api.Avatar{
				Id:       basePreview.Avatar.ID.String(),
				MimeType: basePreview.Avatar.MimeType,
				Url:      basePreview.Avatar.URL,
				Method:   basePreview.Avatar.Method,
			}
		}
	}

	return &api.GetPlayerPreviewsResponse{
		Previews: previews,
	}, nil
}

func newPlayerService(
	playerManager managers.PlayerManager,
	relationManager managers.RelationManager,
	logger *logrus.Logger,
) *service {
	return &service{
		playerManager:   playerManager,
		relationManager: relationManager,
		logger:          logger,
	}
}
