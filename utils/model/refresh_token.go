package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID 			uuid.UUID		`db:"refresh_token"`
	UserId 		uuid.UUID		`db:"user_id"`
	Token		string			`db:"token"`
	Created_At	time.Time		`db:"created_at"`
	Expires_At	time.Time		`db:"expires_at"`
}