package managers

import (
	"context"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/generators"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	RefreshTokenExpirationThresholdHours = 48
)

type AuthenticationManager interface {
	AuthenticateUser(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error)
	RefreshUser(curRefreshToken string) (accessToken string, newRefreshToken *string, err error)
}

type authenticationManager struct {
	userRepository repos.UserRepository
	jwtGenerator   generators.JWTGenerator
}

func ProvideAuthenticationManager(
	userRepository repos.UserRepository,
	jwtGenerator generators.JWTGenerator,
) AuthenticationManager {
	return &authenticationManager{
		userRepository: userRepository,
		jwtGenerator:   jwtGenerator,
	}
}

func (mgr *authenticationManager) AuthenticateUser(
	ctx context.Context,
	email string,
	password string,
) (
	accessToken string,
	refreshToken string,
	returnErr error,
) {
	cred, err := mgr.userRepository.GetUserCredential(ctx, email)
	if err != nil {
		return "", "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(cred.PasswordHash), []byte(password)) != nil {
		return "", "", errors.NewIncorrectPasswordError()
	}

	return mgr.jwtGenerator.GenerateTokenPair(cred.UserId)
}

func (mgr *authenticationManager) RefreshUser(
	refreshToken string,
) (
	returnAccessToken string,
	returnRefreshToken *string,
	returnErr error,
) {
	claims, err := mgr.jwtGenerator.GetRefreshTokenClaims(refreshToken)
	if err != nil {
		return "", nil, err
	}

	userId, expiresAt := claims.UserId, claims.RegisteredClaims.ExpiresAt.Time

	accessToken, err := mgr.jwtGenerator.GenerateAccessToken(userId)
	if err != nil {
		return "", nil, err
	}

	if expiresAt.Sub(time.Now()).Hours() > RefreshTokenExpirationThresholdHours {
		return accessToken, nil, nil
	}

	newRefreshToken, err := mgr.jwtGenerator.GenerateRefreshToken(userId)

	return accessToken, &newRefreshToken, nil
}
