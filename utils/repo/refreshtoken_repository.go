package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/utils/model"
)


type RefreshTokenRepository interface {
	Save(token model.RefreshToken) error
	FindByToken(token string)(*model.RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserId(userId uuid.UUID) error
}

type refreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (rt *refreshTokenRepository) Save(token model.RefreshToken) error {
	_, err := rt.db.Exec(`INSERT INTO refresh_tokens (id, user_id, token, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)`, token.ID, token.UserId, token.Token, token.Created_At, token.Expires_At)
	return err
	
}

func (rt *refreshTokenRepository) FindByToken(token string)(*model.RefreshToken, error) {
	var t model.RefreshToken
	err := rt.db.QueryRow(`SELECT id, user_id, token, created_at, expires_at FROM refresh_tokens WHERE token=$1`, token).Scan(&t.ID, &t.UserId, &t.Token, &t.Created_At, &t.Expires_At)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (rt *refreshTokenRepository) DeleteByToken(token string) error {
	_, err := rt.db.Exec(`DELETE FROM refresh_tokens WHERE token=$1`, token)
	return err
}

func (rt *refreshTokenRepository) DeleteByUserId(userId uuid.UUID) error {
	_, err := rt.db.Exec(`DELETE FROM refresh_tokens WHERE user_id=$1`, userId)
	return err
}
