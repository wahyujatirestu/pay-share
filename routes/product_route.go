package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wahyujatirestu/payshare/controller"
	"github.com/wahyujatirestu/payshare/middleware"
	"github.com/wahyujatirestu/payshare/utils/service"	
)

func ProductRoute(rg *gin.RouterGroup, productController *controller.ProductController, jwtService service.JWTService)  {
	productGroup := rg.Group("/product")
	authMw := middleware.NewAuthMiddleware(jwtService)

	productGroup.GET("/", productController.GetAll)
	productGroup.GET("/:id", productController.GetById)
	
	productGroup.Use(authMw.RequireToken("employee"))
	{
		productGroup.POST("/", productController.Create)
		productGroup.PUT("/:id", productController.Update)
		productGroup.DELETE("/:id", productController.Delete)
	}
}