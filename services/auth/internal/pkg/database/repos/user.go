package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/database/pkg/psql"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/database/tables"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/models"
	goerrors "github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	CheckIfUserCredentialExists(ctx context.Context, email string) (bool, error)
	GetUserCredential(ctx context.Context, email string) (models.Credential, error)
	AddUserCredentialTx(ctx context.Context, userId uuid.UUID, email string, passwordHash string) error
}

type userRepository struct {
	db *gorm.DB
}

func ProvideUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CheckIfUserCredentialExists(
	ctx context.Context,
	email string,
) (checkResult bool, returnErr error) {
	defer psql.PanicHandler(&returnErr)

	var credential = new(tables.Credential)
	if err := r.db.
		WithContext(ctx).
		Where("email = ?", email).
		Take(credential).
		Error; err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, goerrors.WithStack(err)
	}

	return true, nil
}

func (r *userRepository) GetUserCredential(
	ctx context.Context,
	email string,
) (returnCredential models.Credential, returnErr error) {
	defer psql.PanicHandler(&returnErr)

	var credential = new(tables.Credential)
	if err := r.db.
		WithContext(ctx).
		Take(credential, "email = ?", email).
		Error; err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return models.Credential{}, errors.NewUserNotFoundError(email)
		}

		return models.Credential{}, goerrors.WithStack(err)
	}

	return credential.ToCredentialModel(), nil
}

func (r *userRepository) AddUserCredentialTx(
	ctx context.Context,
	userId uuid.UUID,
	email string,
	passwordHash string,
) (returnErr error) {
	tx, err := getTransaction(ctx)
	if err != nil {
		return goerrors.WithStack(err)
	}

	currentTime := time.Now()
	credential := &tables.Credential{
		UserId:       userId,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedOn:    currentTime,
		UpdatedOn:    currentTime,
	}

	if err := tx.
		WithContext(ctx).
		Create(credential).
		Error; err != nil {
		return goerrors.WithStack(err)
	}

	return nil
}
