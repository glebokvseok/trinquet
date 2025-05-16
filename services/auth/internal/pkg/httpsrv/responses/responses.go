package responses

type AuthenticateUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshUserResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}

type RegisterUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
