package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/payshare/controller"
	"github.com/wahyujatirestu/payshare/middleware"
	"github.com/wahyujatirestu/payshare/utils/service"	
)

func TransactionRoute(rg *gin.RouterGroup, tsController *controller.TransactionsController, jwtService service.JWTService)  {
	tsGroup := rg.Group("/transaction")
	authMw := middleware.NewAuthMiddleware(jwtService)

	tsGroup.Use(authMw.RequireToken("customer", "employee"))
	{
		tsGroup.POST("/", tsController.Create)
		tsGroup.GET("/", tsController.GetAll)
		tsGroup.GET("/:id", tsController.GetById)
		tsGroup.PUT("/:id", tsController.Update)
		tsGroup.DELETE("/:id", tsController.Delete)
		tsGroup.POST("/midtrans/webhook", tsController.WebHook)
	}
}