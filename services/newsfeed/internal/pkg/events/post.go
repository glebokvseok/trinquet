package events

import (
	"github.com/google/uuid"
	"time"
)

type PostEventType = string

const (
	CreatePost       PostEventType = "create_post"
	LikePost         PostEventType = "like_post"
	UnlikePost       PostEventType = "unlike_post"
	CommentPost      PostEventType = "comment_post"
	ReplyPostComment PostEventType = "reply_post_comment"
)

type CreatePostEvent struct {
	Text      string      `json:"text" validate:"omitempty"`
	MediaIDs  []uuid.UUID `json:"media_ids" validate:"omitempty,max=5"`
	CreatedOn time.Time   `json:"created_on" validate:"required"`
}

type LikePostEvent struct {
	PostID  uuid.UUID `json:"post_id" validate:"required"`
	LikedOn time.Time `json:"liked_on" validate:"required"`
}

type UnlikePostEvent struct {
	PostID    uuid.UUID `json:"post_id" validate:"required"`
	UnlikedOn time.Time `json:"unliked_on" validate:"required"`
}

type CommentPostEvent struct {
	PostID      uuid.UUID `json:"post_id" validate:"required"`
	CommentText string    `json:"text" validate:"required"`
	CommentedOn time.Time `json:"commented_on" validate:"required"`
}

type ReplyPostCommentEvent struct {
	CommentID uuid.UUID `json:"comment_id" validate:"required"`
	ReplyText string    `json:"text" validate:"required"`
	RepliedOn time.Time `json:"replied_on" validate:"required"`
}
