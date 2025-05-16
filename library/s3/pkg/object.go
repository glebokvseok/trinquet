package s3

import "github.com/google/uuid"

type Object struct {
	ID       uuid.UUID `json:"id"`
	MimeType string    `json:"mime_type"`
	URL      string    `json:"url"`
	Method   string    `json:"method"`
}
