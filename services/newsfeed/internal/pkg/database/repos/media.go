package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/database/pkg/mongo"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/collections"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type MediaRepository interface {
	AddMedia(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID, media any) error
}

func ProvideMediaRepository(
	client mongo.Client,
	config mongo.RequestConfig,
) MediaRepository {
	return &mediaRepository{
		client: client,
		config: config,
	}
}

type mediaRepository struct {
	client mongo.Client
	config mongo.RequestConfig
}

func (r *mediaRepository) AddMedia(
	ctx context.Context,
	userID uuid.UUID,
	mediaID uuid.UUID,
	media any,
) (returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	rawMediaBson, err := bson.Marshal(media)
	if err != nil {
		return errors.WithStack(err)
	}

	var mediaBson bson.M
	err = bson.Unmarshal(rawMediaBson, &mediaBson)
	if err != nil {
		return errors.WithStack(err)
	}

	mediaBson["_id"], mediaBson["user_id"], mediaBson["created_on"] =
		mediaID.String(), userID.String(), time.Now()

	medias := r.client.Database(r.config.Database).Collection(collections.MediasCollectionName)

	_, err = medias.InsertOne(ctx, mediaBson)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
