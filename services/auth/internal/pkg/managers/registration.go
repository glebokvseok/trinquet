package managers

import (
	"context"
	"github.com/google/uuid"
	txmgr "github.com/move-mates/trinquet/library/database/pkg/psql/managers"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/database/repos"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/generators"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationManager interface {
	RegisterUser(ctx context.Context, email string, password string) (accessToken string, refreshToken string, err error)
}

type registrationManager struct {
	transactionManager txmgr.TransactionManager
	userRepository     repos.UserRepository
	jwtGenerator       generators.JWTGenerator
}

func ProvideRegistrationManager(
	transactionManager txmgr.TransactionManager,
	userRepository repos.UserRepository,
	jwtGenerator generators.JWTGenerator,
) RegistrationManager {
	return &registrationManager{
		transactionManager: transactionManager,
		userRepository:     userRepository,
		jwtGenerator:       jwtGenerator,
	}
}

func (mgr *registrationManager) RegisterUser(
	ctx context.Context,
	email string,
	password string,
) (
	accessToken string,
	refreshToken string,
	returnErr error,
) {
	check, err := mgr.userRepository.CheckIfUserCredentialExists(ctx, email)
	if err != nil {
		return "", "", err
	}

	if check {
		return "", "", errors.NewUserAlreadyExistsError(email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	userId := uuid.New()

	ctx, err = mgr.transactionManager.BeginTransaction(ctx)
	if err != nil {
		return "", "", err
	}

	defer mgr.transactionManager.FinalizeTransaction(ctx, &returnErr)

	if err = mgr.userRepository.AddUserCredentialTx(ctx, userId, email, string(hash)); err != nil {
		return "", "", err
	}

	if err = mgr.transactionManager.CommitTransaction(ctx); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = mgr.jwtGenerator.GenerateTokenPair(userId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
