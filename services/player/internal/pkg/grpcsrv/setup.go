package grpcsrv

import (
	"generated/services/player/pkg/api"
	"github.com/move-mates/trinquet/library/common/pkg/grpcsrv"
	"github.com/move-mates/trinquet/services/player/internal/pkg/managers"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ProvideGrpcServerSetupFunc(
	playerManager managers.PlayerManager,
	relationManager managers.RelationManager,
	logger *logrus.Logger,
) grpcsrv.GrpcServerSetupFunc {
	return func(server *grpc.Server) {
		api.RegisterPlayerServiceServer(
			server,
			newPlayerService(
				playerManager,
				relationManager,
				logger,
			),
		)
	}
}
