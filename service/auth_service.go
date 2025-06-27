package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/payshare/model"
	utilsmodel "github.com/wahyujatirestu/payshare/utils/model"
	"github.com/wahyujatirestu/payshare/utils/security"
	"github.com/wahyujatirestu/payshare/utils/service"
	"github.com/wahyujatirestu/payshare/utils/repo"
)

type AuthenticationService interface {
	Login(email string, password string)(string, string, error) 
	RefreshToken(refreshToken string) (string, error)
	Logout(refreshToken string) error
}

type authenticationService struct {
	userService UserService
	jwtService service.JWTService
	rtrepo	repository.RefreshTokenRepository
}

func NewAuthenticationService(userService userService, jwtService service.JWTService, rtrepo repository.RefreshTokenRepository) AuthenticationService {
	return &authenticationService{userService: &userService, jwtService: jwtService, rtrepo: rtrepo}
}

func (a *authenticationService) Login(identifier string, password string) (string, string, error) {
	user, err := a.userService.GetEmailUsername(identifier, password)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("Email not found")
	}

	valid, err := security.VerifyPasswordHash(user.Password, password)
	if err != nil {
		return "", "", err
	}

	if !valid {
		return "", "", errors.New("Invalid Password")
	}

	accessToken, err := a.jwtService.CreateToken(*user)
	if err != nil {
		return "", "", nil
	}

	refreshToken, err := a.jwtService.CreateRefreshToken(*user)
	if err != nil {
		return  "", "", nil
	}

	rt := utilsmodel.RefreshToken{
		ID: uuid.New(),
		UserId: user.ID,
		Token: refreshToken,
		Created_At: time.Now(),
		Expires_At: time.Now().Add(time.Hour * 24 * 7),
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

	// generate new token
	user := model.User{
		ID: claim.UserId,
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
