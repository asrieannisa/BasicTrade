package middlewares

import (
	"BasicTrade/database"
	models "BasicTrade/models/entity"
	requests "BasicTrade/models/request"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func ProductAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()
		productUUID := ctx.Param("productUUID")

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		admin_ID := uint(userData["id"].(float64))

		var getProduct models.Products
		err := db.Select("admin_id").Where("uuid = ?", productUUID).First(&getProduct).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Data Not Found",
			})
			return
		}

		if getProduct.Admin_ID != admin_ID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}

		ctx.Next()
	}
}

func VariantAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()
		variantUUID := ctx.Param("variantUUID")

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		admin_ID := uint(userData["id"].(float64))

		var getVariant models.Variants
		err := db.Select("product_id").Where("uuid = ?", variantUUID).First(&getVariant).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Data Not Found",
			})
			return
		}

		var getProduct models.Products
		err2 := db.Select("admin_id").Where("id = ?", getVariant.Product_ID).First(&getProduct).Error
		if err2 != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Data Not Found",
			})
			return
		}

		if getProduct.Admin_ID != admin_ID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}

		ctx.Next()
	}
}

func CreateVariantAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()

		var variantReq requests.VariantRequest
		if err := ctx.ShouldBind(&variantReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userData := ctx.MustGet("userData").(jwt5.MapClaims)
		admin_ID := uint(userData["id"].(float64))

		var getVariant models.Variants
		err := db.Select("product_id").Where("uuid = ?", variantReq.Product_id).First(&getVariant).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Data Not Found",
			})
			return
		}

		var getProduct models.Products
		err2 := db.Select("admin_id").Where("id = ?", getVariant.Product_ID).First(&getProduct).Error
		if err2 != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   err.Error(),
				"message": "Data Not Found",
			})
			return
		}

		if getProduct.Admin_ID != admin_ID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this data",
			})
			return
		}

		ctx.Next()
	}
}
