package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/payshare/dto"
	"github.com/wahyujatirestu/payshare/service"
)

type AuthController struct {
	authService service.AuthenticationService
}

func NewAuthController(authService service.AuthenticationService) *AuthController {
	return &AuthController{authService: authService}
}



func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := c.authService.Login(req.Identifier, req.Password)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	res := dto.LoginResponse{
		AccessToken:  	accessToken,
		RefreshToken: 	refreshToken,
		Message: 		"Login Successfully",
	}

	ctx.JSON(200, res)
}


func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"access token": accessToken})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := c.authService.Logout(req.RefreshToken)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to logout"})
	}

	ctx.JSON(200, gin.H{"message": "Logout Successfully"})
}