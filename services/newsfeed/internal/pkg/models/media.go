package models

import s3 "github.com/move-mates/trinquet/library/s3/pkg"

type (
	Media     s3.Object
	MediaType string
)

const (
	PhotoMediaType MediaType = "photo"
	VideoMediaType MediaType = "video"
)

type Photo struct {
	MediaType MediaType `json:"media_type" bson:"media_type" validate:"required"`
	MimeType  string    `json:"mime_type" bson:"mime_type" validate:"required,max=16"`
}

type Video struct {
	MediaType MediaType `json:"media_type" bson:"media_type" validate:"required"`
	MimeType  string    `json:"mime_type" bson:"mime_type" validate:"required,max=16"`
	Duration  string    `json:"duration" bson:"duration" validate:"required"`
}
