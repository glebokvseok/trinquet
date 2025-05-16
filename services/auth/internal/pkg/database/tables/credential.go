package tables

import (
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/services/auth/internal/pkg/models"
	"time"
)

type Credential struct {
	UserId       uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
	Email        string    `gorm:"column:email"`
	PasswordHash string    `gorm:"column:password_hash"`
	CreatedOn    time.Time `gorm:"column:created_on"`
	UpdatedOn    time.Time `gorm:"column:updated_on"`
}

func (Credential) TableName() string {
	return "auth_service.credential"
}

func (credential Credential) ToCredentialModel() models.Credential {
	return models.Credential{
		UserId:       credential.UserId,
		Email:        credential.Email,
		PasswordHash: credential.PasswordHash,
	}
}
