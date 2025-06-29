package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/model"
	utilsmodel "github.com/wahyujatirestu/payshare/utils/model"
	"github.com/wahyujatirestu/payshare/utils/security"
	utils "github.com/wahyujatirestu/payshare/utils/service"
	"github.com/wahyujatirestu/payshare/utils/repo"
)

type AuthenticationService interface {
	Login(identifier string, password string)(string, string, error) 
	RefreshToken(refreshToken string) (string, error)
	Logout(refreshToken string) error
}

type authenticationService struct {
	userService UserService
	jwtService utils.JWTService
	rtrepo	repository.RefreshTokenRepository
}

func NewAuthenticationService(userService UserService, jwtService utils.JWTService, rtrepo repository.RefreshTokenRepository) AuthenticationService {
	return &authenticationService{userService: userService, jwtService: jwtService, rtrepo: rtrepo}
}

func (a *authenticationService) Login(identifier string, password string) (string, string, error) {
	user, err := a.userService.GetEmailUsername(identifier)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("user not found")
	}

	valid, err := security.VerifyPasswordHash(user.Password, password)
	if err != nil {
		return "", "", err
	}

	if !valid {
		return "", "", errors.New("invalid password")
	}

	accessToken, err := a.jwtService.CreateToken(*user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.jwtService.CreateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	rt := utilsmodel.RefreshToken{
		ID: uuid.New(),
		UserId: user.ID,
		Token: refreshToken,
		Created_At: time.Now(),
		Expires_At: time.Now().Add(7 * 24 * time.Hour),
	}

	err = a.rtrepo.Save(rt)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}



func (a *authenticationService) RefreshToken(rt string) (string, error) {
	// verifikasi token
	claim, err := a.jwtService.VerifyRefreshToken(rt)
	if err != nil {
		return "", err
	}

	// cek apakah refresh token berada di DB dan belum expired
	storedToken, err := a.rtrepo.FindByToken(rt)
	if err != nil {
		return "", errors.New("Refresh token not found")
	}

	if storedToken.Expires_At.Before(time.Now()) {
		return "", errors.New("Refresh token expired")
	}

	userId, err := uuid.Parse(string(claim.UserId))
	if err != nil {
		return "", errors.New("Invalid user ID in token")
	}

	// generate new token
	user := model.User{
		ID: userId,
		Role: claim.Role,
	}

	accessToken, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", err
	}
	return accessToken, nil

}

func (a *authenticationService) Logout(token string) error {
	err := a.rtrepo.DeleteByToken(token)
	if err != nil {
		return err
	}
	return nil
}
