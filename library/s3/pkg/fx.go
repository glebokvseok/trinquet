package s3

import (
	"github.com/move-mates/trinquet/library/common/pkg/mfx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"s3",
		fx.Provide(
			provideS3Client,
			provideS3PresignClient,
			provideLinkGenerator,
			mfx.ProvideConfig[linkGeneratorConfig](configSectionName),
		),
	)
}
