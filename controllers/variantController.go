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

func CreateVariant(ctx *gin.Context) {
	db := database.GetDB()

	var variantReq requests.VariantRequest
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	contentType := helpers.GetContentType(ctx)
	Product_ID := uint(userData["id"].(float64))

	Variant := models.Variants{
		Variant_name: variantReq.Variant_name,
		Quantity:     variantReq.Quantity,
		Product_ID:   Product_ID,
	}

	// Generate a new UUID
	newUUID := uuid.New()
	Variant.UUID = newUUID.String() // Set the generated UUID as the ID

	if contentType == appJSON {
		ctx.ShouldBindJSON(&Variant)
	} else {
		ctx.ShouldBind(&Variant)
	}

	if err := db.Debug().Create(&Variant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": Variant,
	})
}

func UpdateVariant(ctx *gin.Context) {
	db := database.GetDB()

	userData := ctx.MustGet("userData").(jwt5.MapClaims)
	contentType := helpers.GetContentType(ctx)
	Variant := models.Variants{}

	variantUUID := ctx.Param("variantUUID")
	product_ID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&Variant)
	} else {
		ctx.ShouldBind(&Variant)
	}

	// Retrieve existing book from the database
	var getVariant models.Variants
	if err := db.Model(&getVariant).Where("uuid = ?", variantUUID).First(&getVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// Update the Book struct with retrieved data
	Variant.ID = uint(getVariant.ID)
	Variant.Product_ID = product_ID

	var variantReq requests.VariantRequest
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the book record in the database
	updateData := models.Variants{
		Variant_name: variantReq.Variant_name,
		Quantity:     variantReq.Quantity,
		Product_ID:   product_ID,
	}

	if err := db.Model(&Variant).Where("uuid = ?", variantUUID).Updates(updateData).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": Variant,
	})
}
