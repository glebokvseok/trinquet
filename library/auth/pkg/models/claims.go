package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserId uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}
