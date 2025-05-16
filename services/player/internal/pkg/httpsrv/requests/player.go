package requests

type CreatePlayerRequest struct {
	Username string `json:"username" validate:"required,username"`
}
