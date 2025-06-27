package service

import "github.com/wahyujatirestu/payshare/utils/service"

type AuthenticationService interface {
	Login(email string, password string)(string, error) 
}

type authenticationService struct {
	userService UserService
	jwtService service.JWTService
}

func NewAuthenticationService(userService userService, jwtService service.JWTService) AuthenticationService {
	return &authenticationService{userService: &userService, jwtService: jwtService}
}

func (a *authenticationService) Login(email string, password string) (string, error) {
	user, err := a.userService.Login(email, password)
	if err != nil {
		return "", err
	}

	token, err := a.jwtService.CreateToken(*user)
	if err != nil {
		return "", nil
	}

	return token, nil
}