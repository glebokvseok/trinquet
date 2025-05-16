package requests

type UploadAvatarRequest struct {
	MimeType string `json:"mime_type" validate:"required,max=16"`
}
