package models

import "github.com/google/uuid"

type Credential struct {
	UserId       uuid.UUID
	Email        string
	PasswordHash string
}
