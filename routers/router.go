package router

import (
	"BasicTrade/controllers"
	"BasicTrade/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	adminRouter := router.Group("/admins")
	{
		adminRouter.POST("/register", controllers.AdminRegister)
		adminRouter.POST("/login", controllers.AdminLogin)
	}

	productRouter := router.Group("/products")
	{
		productRouter.GET("/", controllers.GetProducts)

		productRouter.Use(middlewares.Authentication())
		productRouter.POST("/", controllers.CreateProduct)
		productRouter.PUT("/:productUUID", middlewares.ProductAuthorization(), controllers.UpdateProduct)
	}

	return router
}
