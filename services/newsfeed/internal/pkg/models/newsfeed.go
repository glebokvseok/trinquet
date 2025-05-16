package models

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/playersvc/pkg/models"
	s3 "github.com/move-mates/trinquet/library/s3/pkg"
)

type (
	Avatar s3.Object
)

type Newsfeed struct {
	Posts        []Post `json:"posts"`
	Cursor       int64  `json:"cursor"`
	HasMorePosts bool   `json:"has_more_posts"`
}

type Post struct {
	ID           uuid.UUID        `json:"id"`
	Author       Author           `json:"author"`
	Text         string           `json:"text"`
	Medias       []map[string]any `json:"medias"`
	IsLiked      bool             `json:"is_liked"`
	LikeCount    int64            `json:"like_count"`
	CommentCount int64            `json:"comment_count"`
	Timestamp    int64            `json:"timestamp"`
}

type Author struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Avatar  *Avatar   `json:"avatar"`
}

type CommentSection struct {
	Comments        []Comment `json:"comments"`
	Cursor          int64     `json:"cursor"`
	HasMoreComments bool      `json:"has_more_comments"`
}

type Comment struct {
	ID           uuid.UUID      `json:"id"`
	Author       Author         `json:"author"`
	SelfAuthored bool           `json:"self_authored"`
	Text         string         `json:"text"`
	Replies      []CommentReply `json:"replies"`
	Timestamp    int64          `json:"timestamp"`
}

type CommentReply struct {
	Author       Author `json:"author"`
	SelfAuthored bool   `json:"self_authored"`
	Text         string `json:"text"`
	Timestamp    int64  `json:"timestamp"`
}

func (author *Author) Complete(preview models.PlayerPreview) {
	author.Name = preview.Name
	author.Surname = preview.Surname
	author.Avatar = (*Avatar)(preview.Avatar)
}
