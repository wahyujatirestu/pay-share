package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wahyujatirestu/payshare/config"
	"github.com/wahyujatirestu/payshare/model"
	modeljwt "github.com/wahyujatirestu/payshare/utils/model"
)

type JWTService interface {
	CreateToken(user model.User)(string, error)
	VerifyToken(tokenString string)(modeljwt.JWTPayloadClaim, error)
	CreateRefreshToken(user model.User)(string, error)
	VerifyRefreshToken(tokenString string)(modeljwt.JWTPayloadClaim, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func NewJWTService(cfg config.TokenConfig) JWTService {
	return  &jwtService{cfg: cfg}
}

func (j *jwtService) CreateToken(user model.User) (string, error) {
	tokenKey := j.cfg.JWTSignatureKey
	claim := modeljwt.JWTPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: j.cfg.AppName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTokenLifetime.Abs())),
		},
		UserId: user.ID,
		Role: user.Role,
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JWTSigningMethod, claim)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	} 
	return token, nil
}

func (j *jwtService) VerifyToken(ts string) (modeljwt.JWTPayloadClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(ts, modeljwt.JWTPayloadClaim{}, func(t *jwt.Token) (interface{}, error) {
		return j.cfg.JWTSignatureKey, nil
	})

	if err != nil {
		return modeljwt.JWTPayloadClaim{}, err
	}

	claim, ok := tokenParse.Claims.(*modeljwt.JWTPayloadClaim)
	if !ok {
		return modeljwt.JWTPayloadClaim{}, errors.New("invalid token")
	}

	return *claim, nil
}

func (j *jwtService) CreateRefreshToken(user model.User) (string, error) {
	tokenKey := j.cfg.JWTSignatureKey
	refreshTokenLifetime := time.Hour * 24 * 7

	claim := modeljwt.JWTPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: j.cfg.AppName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenLifetime)),
		},
		UserId: user.ID,
		Role: user.Role,
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JWTSigningMethod, claim)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (j *jwtService) VerifyRefreshToken(ts string) (modeljwt.JWTPayloadClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(ts, modeljwt.JWTPayloadClaim{}, func(t *jwt.Token) (interface{}, error) {
		return j.cfg.JWTSignatureKey, nil
	})

	if err != nil {
		return modeljwt.JWTPayloadClaim{}, err
	}

	claim, ok := tokenParse.Claims.(*modeljwt.JWTPayloadClaim)
	if !ok {
		return modeljwt.JWTPayloadClaim{}, errors.New("Invalid refresh token")
	}

	return  *claim, nil
}