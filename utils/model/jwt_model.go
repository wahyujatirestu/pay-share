package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTPayloadClaim struct {
	jwt.RegisteredClaims
	UserId	uuid.UUID
	Role	string
}