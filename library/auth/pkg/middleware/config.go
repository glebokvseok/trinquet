package authmw

const configSectionName = "auth.jwt"

type JWTConfig struct {
	UserAccessTokenSigningKey string `yaml:"user_access_token_signing_key"`
}
