package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/database/pkg/mongo"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/database/collections"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/events"
	"github.com/move-mates/trinquet/services/newsfeed/internal/pkg/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	gomongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
)

type PostRepository interface {
	CreatePost(ctx context.Context, userID uuid.UUID, postID uuid.UUID, event events.CreatePostEvent) error
	LikePost(ctx context.Context, userID uuid.UUID, event events.LikePostEvent) error
	UnlikePost(ctx context.Context, userID uuid.UUID, event events.UnlikePostEvent) error
	CommentPost(ctx context.Context, userID uuid.UUID, event events.CommentPostEvent) error
	ReplyPostComment(ctx context.Context, userID uuid.UUID, event events.ReplyPostCommentEvent) error
	GetPostsLikeCount(ctx context.Context, postIds []uuid.UUID) ([]int64, error)
	GetPostsCommentCount(ctx context.Context, postIds []uuid.UUID) ([]int64, error)
	GetPostComments(ctx context.Context, postId uuid.UUID, cursor int64, commentCount int) ([]models.Comment, error)
}

func ProvidePostRepository(
	client mongo.Client,
	config mongo.RequestConfig,
) PostRepository {
	return &postRepository{
		client: client,
		config: config,
	}
}

type postRepository struct {
	client mongo.Client
	config mongo.RequestConfig
}

func (r *postRepository) CreatePost(
	ctx context.Context,
	userID uuid.UUID,
	postID uuid.UUID,
	event events.CreatePostEvent,
) (returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	mediaIDs := make([]string, 0)
	for _, mediaID := range event.MediaIDs {
		mediaIDs = append(mediaIDs, mediaID.String())
	}

	post := collections.Post{
		ID:        postID.String(),
		AuthorID:  userID.String(),
		Text:      event.Text,
		MediaIDs:  mediaIDs,
		Timestamp: event.CreatedOn.UnixMilli(),
	}

	posts := r.client.Database(r.config.Database).Collection(collections.PostsCollection)

	_, err := posts.InsertOne(ctx, post)

	return errors.WithStack(err)
}

func (r *postRepository) LikePost(
	ctx context.Context,
	userID uuid.UUID,
	event events.LikePostEvent,
) (returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	posts := r.client.Database(r.config.Database).Collection(collections.PostsCollection)

	err := posts.FindOne(ctx, bson.M{"_id": event.PostID.String()}).Err()
	if err != nil {
		return errors.WithStack(err)
	}

	search := bson.M{
		"post_id": event.PostID.String(),
		"user_id": userID.String(),
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":          true,
			"client_modified_on": event.LikedOn,
		},
		"$setOnInsert": bson.M{
			"post_id":           event.PostID.String(),
			"user_id":           userID.String(),
			"client_created_on": event.LikedOn,
		},
	}

	likes := r.client.Database(r.config.Database).Collection(collections.PostLikesCollectionName)

	_, err = likes.UpdateOne(ctx, search, update, options.Update().SetUpsert(true))

	return errors.WithStack(err)
}

func (r *postRepository) UnlikePost(
	ctx context.Context,
	userID uuid.UUID,
	event events.UnlikePostEvent,
) (returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	search := bson.M{
		"post_id":   event.PostID.String(),
		"user_id":   userID.String(),
		"is_active": true,
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":          false,
			"client_modified_on": event.UnlikedOn,
		},
	}

	likes := r.client.Database(r.config.Database).Collection(collections.PostLikesCollectionName)

	_, err := likes.UpdateOne(ctx, search, update)

	return errors.WithStack(err)
}

func (r *postRepository) CommentPost(
	ctx context.Context,
	userID uuid.UUID,
	event events.CommentPostEvent,
) (returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	posts := r.client.Database(r.config.Database).Collection(collections.PostsCollection)

	err := posts.FindOne(ctx, bson.M{"_id": event.PostID.String()}).Err()
	if err != nil {
		return errors.WithStack(err)
	}

	comment := collections.PostComment{
		ID:        uuid.New().String(),
		PostID:    event.PostID.String(),
		AuthorID:  userID.String(),
		Text:      event.CommentText,
		IsDeleted: false,
		Timestamp: event.CommentedOn.UnixMilli(),
		Replies:   make([]collections.PostCommentReply, 0),
	}

	comments := r.client.Database(r.config.Database).Collection(collections.PostCommentsCollectionName)

	_, err = comments.InsertOne(ctx, comment)

	return errors.WithStack(err)
}

func (r *postRepository) ReplyPostComment(
	ctx context.Context,
	userID uuid.UUID,
	event events.ReplyPostCommentEvent,
) (returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	search := bson.M{
		"_id":        event.CommentID.String(),
		"is_deleted": false,
	}

	reply := collections.PostCommentReply{
		AuthorID:  userID.String(),
		Text:      event.ReplyText,
		IsDeleted: false,
		Timestamp: event.RepliedOn.UnixMilli(),
	}

	update := bson.M{
		"$push": bson.M{
			"replies": reply,
		},
	}

	comments := r.client.Database(r.config.Database).Collection(collections.PostCommentsCollectionName)

	res, err := comments.UpdateOne(ctx, search, update)
	if err != nil {
		return errors.WithStack(err)
	}

	if res.ModifiedCount == 0 {
		return errors.New("comment not found")
	}

	return nil
}

func (r *postRepository) GetPostsLikeCount(ctx context.Context, postIds []uuid.UUID) (likeCount []int64, returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	ids := make([]string, len(postIds))
	for i, postId := range postIds {
		ids[i] = postId.String()
	}

	pipeline := gomongo.Pipeline{
		bson.D{{"$match", bson.D{
			{"post_id", bson.M{"$in": ids}},
			{"is_active", true},
		}}},

		bson.D{{"$group", bson.D{
			{"_id", "$post_id"},
			{"count", bson.M{"$sum": 1}},
		}}},
	}

	return r.getCounts(ctx, pipeline, collections.PostLikesCollectionName, ids)
}

func (r *postRepository) GetPostsCommentCount(ctx context.Context, postIds []uuid.UUID) (commentCount []int64, returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	ids := make([]string, len(postIds))
	for i, postId := range postIds {
		ids[i] = postId.String()
	}

	pipeline := gomongo.Pipeline{
		bson.D{{"$match", bson.D{
			{"post_id", bson.M{"$in": ids}},
			{"is_deleted", false},
		}}},

		bson.D{{"$project", bson.D{
			{"post_id", 1},
			{"reply_count", bson.D{{"$size", bson.D{{"$filter", bson.D{
				{"input", "$replies"},
				{"as", "reply"},
				{"cond", bson.D{{"$eq", bson.A{"$$reply.is_deleted", false}}}},
			}}}}}},
		}}},

		bson.D{{"$group", bson.D{
			{"_id", "$post_id"},
			{"count", bson.D{{"$sum", bson.D{{"$add", bson.A{1, "$reply_count"}}}}}},
		}}},
	}

	return r.getCounts(ctx, pipeline, collections.PostCommentsCollectionName, ids)
}

func (r *postRepository) getCounts(
	ctx context.Context,
	pipeline gomongo.Pipeline,
	collectionName string,
	ids []string,
) (counts []int64, returnErr error) {
	collection := r.client.Database(r.config.Database).Collection(collectionName)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer cursor.Close(ctx)

	var res struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}

	countMap := make(map[string]int64)
	for cursor.Next(ctx) {
		if err = cursor.Decode(&res); err != nil {
			return nil, errors.WithStack(err)
		}
		countMap[res.ID] = res.Count
	}

	counts = make([]int64, len(ids))
	for i, id := range ids {
		counts[i] = countMap[id]
	}

	return counts, nil
}

func (r *postRepository) GetPostComments(
	ctx context.Context,
	postId uuid.UUID,
	cursorPos int64,
	commentCount int,
) (returnComments []models.Comment, returnErr error) {
	defer mongo.PanicHandler(&returnErr)
	ctx, cancel := context.WithTimeout(ctx, r.config.RequestTimeout)
	defer cancel()

	postComments := r.client.Database(r.config.Database).Collection(collections.PostCommentsCollectionName)

	cursor, err := postComments.Aggregate(ctx, newGetPostCommentsPipeline(postId, cursorPos, commentCount))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer cursor.Close(ctx)

	var comments []collections.PostComment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, errors.WithStack(err)
	}

	returnComments = make([]models.Comment, len(comments))
	for i, comment := range comments {
		commentID, err := uuid.Parse(comment.ID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		authorID, err := uuid.Parse(comment.AuthorID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		replies := make([]models.CommentReply, len(comment.Replies))
		for j, reply := range comment.Replies {
			authorID, err := uuid.Parse(comment.AuthorID)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			replies[j] = models.CommentReply{
				Author: models.Author{
					ID: authorID,
				},
				Text:      reply.Text,
				Timestamp: reply.Timestamp,
			}
		}

		sort.Slice(replies, func(i, j int) bool {
			return replies[i].Timestamp > replies[j].Timestamp
		})

		returnComments[i] = models.Comment{
			ID: commentID,
			Author: models.Author{
				ID: authorID,
			},
			Text:      comment.Text,
			Timestamp: comment.Timestamp,
			Replies:   replies,
		}
	}

	return returnComments, nil
}

func newGetPostCommentsPipeline(postID uuid.UUID, cursor int64, postCount int) gomongo.Pipeline {
	return gomongo.Pipeline{
		bson.D{{"$sort", bson.D{{"timestamp", -1}}}},

		bson.D{{"$match", bson.D{
			{"post_id", postID.String()},
			{"timestamp", bson.D{{"$lt", cursor}}},
			{"is_deleted", false},
		}}},

		bson.D{{"$limit", postCount}},
	}
}
