package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/database/pkg/mongo"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/collections"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	gomongo "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type NewsfeedRepository interface {
	AddPostToNewsfeeds(ctx context.Context, userIDs []uuid.UUID, postID uuid.UUID, postCreatedOn time.Time) (int64, error)
	GetNewsfeedPosts(ctx context.Context, userID uuid.UUID, cursor int64, postCount int) (posts []models.Post, returnErr error)
}

func ProvideNewsfeedRepository(
	client mongo.Client,
	config mongo.RequestConfig,
) NewsfeedRepository {
	return &newsfeedRepository{
		client: client,
		config: config,
	}
}

type newsfeedRepository struct {
	client mongo.Client
	config mongo.RequestConfig
}

func (r *newsfeedRepository) AddPostToNewsfeeds(
	ctx context.Context,
	userIDs []uuid.UUID,
	postID uuid.UUID,
	postCreatedOn time.Time,
) (addedCount int64, returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	bulkInsertOps := make([]gomongo.WriteModel, len(userIDs))
	for i := range userIDs {
		pst := collections.NewsfeedPost{
			UserID:    userIDs[i].String(),
			PostID:    postID.String(),
			Timestamp: postCreatedOn.UnixMilli(),
		}

		bulkInsertOps[i] = gomongo.NewInsertOneModel().SetDocument(pst)
	}

	newsfeedPosts := r.client.Database(r.config.Database).Collection(collections.NewsfeedPostsCollectionName)

	res, err := newsfeedPosts.BulkWrite(ctx, bulkInsertOps)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return res.InsertedCount, nil
}

func (r *newsfeedRepository) GetNewsfeedPosts(
	ctx context.Context,
	userID uuid.UUID,
	cursorPos int64,
	postCount int,
) (posts []models.Post, returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	newsfeedPosts := r.client.Database(r.config.Database).Collection(collections.NewsfeedPostsCollectionName)

	cursor, err := newsfeedPosts.Aggregate(ctx, newGetNewsfeedPostsPipeline(userID, cursorPos, postCount))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer cursor.Close(ctx)

	var aggPosts []collections.AggregatedPost
	if err = cursor.All(ctx, &aggPosts); err != nil {
		return nil, errors.WithStack(err)
	}

	posts = make([]models.Post, len(aggPosts))
	for i, aggPost := range aggPosts {
		postID, err := uuid.Parse(aggPost.ID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		authorID, err := uuid.Parse(aggPost.AuthorID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		posts[i] = models.Post{
			ID: postID,
			Author: models.Author{
				ID: authorID,
			},
			Text:      aggPost.Text,
			IsLiked:   aggPost.IsLiked,
			Medias:    aggPost.Medias,
			Timestamp: aggPost.Timestamp,
		}
	}

	return posts, nil
}

func newGetNewsfeedPostsPipeline(userID uuid.UUID, cursorPosition int64, postCount int) gomongo.Pipeline {
	return gomongo.Pipeline{
		bson.D{{"$sort", bson.D{{"timestamp", -1}}}},

		bson.D{{"$match", bson.D{
			{"user_id", userID.String()},
			{"timestamp", bson.D{{"$lt", cursorPosition}}},
			{"is_deleted", false},
		}}},

		bson.D{{"$limit", postCount}},

		bson.D{{"$lookup", bson.D{
			{"from", collections.PostsCollection},
			{"localField", "post_id"},
			{"foreignField", "_id"},
			{"as", "post"},
		}}},

		bson.D{{"$unwind", "$post"}},

		bson.D{{"$lookup", bson.D{
			{"from", collections.MediasCollectionName},
			{"localField", "post.media_ids"},
			{"foreignField", "_id"},
			{"as", "medias"},
			{"pipeline", bson.A{
				bson.D{{"$project", bson.D{
					{"id", "$_id"},
					{"_id", 0},
					{"media_type", 1},
					{"mime_type", 1},
					{"duration", 1},
				}}},
			}},
		}}},

		bson.D{{"$lookup", bson.D{
			{"from", collections.PostLikesCollectionName},
			{"let", bson.D{
				{"postId", "$post._id"},
			}},
			{"pipeline", bson.A{
				bson.D{{"$match", bson.D{
					{"$expr", bson.D{
						{"$and", bson.A{
							bson.D{{"$eq", bson.A{"$post_id", "$$postId"}}},
							bson.D{{"$eq", bson.A{"$user_id", userID.String()}}},
							bson.D{{"$eq", bson.A{"$is_active", true}}},
						}},
					}},
				}}},
			}},
			{"as", "like"},
		}}},

		bson.D{{"$addFields", bson.D{
			{"is_liked", bson.D{
				{"$gt", bson.A{
					bson.D{{"$size", "$like"}},
					0,
				}},
			}},
		}}},

		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"id", "$post._id"},
			{"author_id", "$post.author_id"},
			{"text", "$post.text"},
			{"timestamp", "$post.timestamp"},
			{"is_liked", "$is_liked"},
			{"medias", "$medias"},
		}}},
	}
}
