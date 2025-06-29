package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/payshare/controller"
	"github.com/wahyujatirestu/payshare/middleware"
	"github.com/wahyujatirestu/payshare/utils/service"
)

func UserRoute(rg *gin.RouterGroup, userController *controller.UserController, jwtService service.JWTService)  {
	userGroup := rg.Group("/user")
	auhtMw := middleware.NewAuthMiddleware(jwtService)

	userGroup.Use(auhtMw.RequireToken("employee", "customer"))
	{		
		userGroup.GET("/", userController.GetAllUser)
		userGroup.GET("/:id", userController.GetUserById)
		userGroup.PUT("/:id", userController.UpdateUser)
		userGroup.DELETE("/:id", userController.DeleteUser)
	}
}