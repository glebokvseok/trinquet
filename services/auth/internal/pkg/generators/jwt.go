package generators

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/auth/pkg/models"
	"github.com/pkg/errors"
	"time"
)

const (
	JWTConfigSectionName = "auth.jwt"
)

type JWTConfig struct {
	UserAccessTokenSigningKey  string        `yaml:"user_access_token_signing_key"`
	UserAccessTokenLifetime    time.Duration `yaml:"user_access_token_lifetime"`
	UserRefreshTokenSigningKey string        `yaml:"user_refresh_token_signing_key"`
	UserRefreshTokenLifetime   time.Duration `yaml:"user_refresh_token_lifetime"`
}

type JWTGenerator interface {
	GenerateAccessToken(userId uuid.UUID) (accessToken string, err error)
	GenerateRefreshToken(userId uuid.UUID) (accessToken string, err error)
	GenerateTokenPair(userId uuid.UUID) (accessToken string, refreshToken string, err error)
	GetRefreshTokenClaims(token string) (claims *models.CustomClaims, err error)
}

type jwtGenerator struct {
	config JWTConfig
}

func ProvideJWTGenerator(config JWTConfig) JWTGenerator {
	return &jwtGenerator{
		config: config,
	}
}

func (gen *jwtGenerator) GenerateAccessToken(userId uuid.UUID) (accessToken string, returnErr error) {
	return generateToken(userId, gen.config.UserAccessTokenSigningKey, gen.config.UserAccessTokenLifetime)
}

func (gen *jwtGenerator) GenerateRefreshToken(userId uuid.UUID) (refreshToken string, returnErr error) {
	return generateToken(userId, gen.config.UserRefreshTokenSigningKey, gen.config.UserRefreshTokenLifetime)
}

func (gen *jwtGenerator) GenerateTokenPair(userId uuid.UUID) (
	accessToken string,
	refreshToken string,
	returnErr error,
) {
	accessToken, err := gen.GenerateAccessToken(userId)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	refreshToken, err = gen.GenerateRefreshToken(userId)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	return accessToken, refreshToken, nil
}

func generateToken(
	userId uuid.UUID,
	tokenSigningKey string,
	tokenLifetime time.Duration,
) (token string, returnErr error) {
	claims := models.CustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
			Audience:  []string{"user"},
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(tokenSigningKey))

	return token, errors.WithStack(err)
}

func (gen *jwtGenerator) GetRefreshTokenClaims(token string) (claims *models.CustomClaims, returnErr error) {
	return getTokenClaims(token, gen.config.UserRefreshTokenSigningKey)
}

func getTokenClaims(token string, tokenSigningKey string) (claims *models.CustomClaims, returnErr error) {
	claims = new(models.CustomClaims)

	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) { return []byte(tokenSigningKey), nil })
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
