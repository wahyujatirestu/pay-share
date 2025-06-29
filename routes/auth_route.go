package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/payshare/controller"
)

func AuthRoute(rg *gin.RouterGroup, authController *controller.AuthController, userController *controller.UserController)  {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", userController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/refresh", authController.RefreshToken)
		authGroup.POST("/logout", authController.Logout)
	}
}