package router

import (
	"BasicTrade/controllers"
	"BasicTrade/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	adminRouter := router.Group("/auth")
	{
		adminRouter.POST("/register", controllers.AdminRegister)
		adminRouter.POST("/login", controllers.AdminLogin)
	}

	productRouter := router.Group("/products")
	{
		productRouter.GET("/", controllers.GetProducts)
		productRouter.GET("/:productUUID", controllers.GetProductById)
		productRouter.GET("/variants/", controllers.GetVariants)
		productRouter.GET("/variants/:variantUUID", controllers.GetVariantById)
		productRouter.Use(middlewares.Authentication())

		productRouter.POST("/", controllers.CreateProduct)
		productRouter.PUT("/:productUUID", middlewares.ProductAuthorization(), controllers.UpdateProduct)
		productRouter.DELETE("/:productUUID", middlewares.ProductAuthorization(), controllers.DeleteProduct)

		productRouter.POST("/variants/", middlewares.ProductAuthorization(), controllers.CreateVariant)
		productRouter.PUT("/variants/:variantUUID", middlewares.VariantAuthorization(), controllers.UpdateVariant)
		productRouter.DELETE("/variants/:variantUUID", middlewares.VariantAuthorization(), controllers.DeleteVariant)
	}

	return router
}
