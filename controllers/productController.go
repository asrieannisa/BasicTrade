package controllers

import (
	"BasicTrade/database"
	"BasicTrade/helpers"
	models "BasicTrade/models/entity"
	requests "BasicTrade/models/request"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetProducts(ctx *gin.Context) {
	db := database.GetDB()

	results := []models.Products{}

	err := db.Debug().Preload("Admin").Preload("Variants").Find(&results).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": results,
	})

}

func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()

	var productReq requests.ProductRequest
	if err := ctx.ShouldBind(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the filename without extension
	fileName := helpers.RemoveExtension(productReq.Image.Filename)

	uploadResult, err := helpers.UploadFile(productReq.Image, fileName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	admin_ID := uint(userData["id"].(float64))

	Product := models.Products{
		Name:     productReq.Name,
		ImageURL: uploadResult,
	}

	// Generate a new UUID
	newUUID := uuid.New()
	Product.UUID = newUUID.String() // Set the generated UUID as the ID

	Product.Admin_ID = admin_ID

	err = db.Debug().Create(&Product).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": Product,
	})
}

func UpdateProduct(ctx *gin.Context) {
	db := database.GetDB()

	var productReq requests.ProductRequestUpdate
	if err := ctx.ShouldBind(&productReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productUUID := ctx.Param("productUUID")

	var getProduct models.Products
	if err := db.Model(&getProduct).Where("uuid = ?", productUUID).First(&getProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	admin_ID := uint(userData["id"].(float64))

	Product := models.Products{}
	Product.ID = uint(getProduct.ID)
	Product.Admin_ID = admin_ID
	Product.ImageURL = getProduct.ImageURL
	Product.Name = getProduct.Name

	// Update the product record in the database
	updateData := models.Products{
		Name: productReq.Name,
	}

	if err := db.Model(&Product).Where("uuid = ?", productUUID).Updates(updateData).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": Product,
	})
}

func DeleteProduct(ctx *gin.Context) {
	db := database.GetDB()

	productUUID := ctx.Param("productUUID")

	// Retrieve existing product from the database
	var product models.Products
	if err := db.Where("uuid = ?", productUUID).First(&product).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// Delete the product record in the database
	if err := db.Delete(&product).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"success": true,
	})
}

func GetProductById(ctx *gin.Context) {
	db := database.GetDB()

	productUUID := ctx.Param("productUUID")

	// Retrieve existing product from the database
	var getProduct models.Products
	if err := db.Model(&getProduct).Where("UUID = ?", productUUID).First(&getProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": getProduct,
	})

}
