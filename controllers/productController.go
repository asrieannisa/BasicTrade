package controllers

import (
	"BasicTrade/database"
	"BasicTrade/helpers"
	models "BasicTrade/models/entity"
	requests "BasicTrade/models/request"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetProducts(ctx *gin.Context) {
	db := database.GetDB()

	// Mengambil parameter "offset" dan "limit" dari permintaan
	page, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// Mengambil parameter pencarian "search" dari permintaan
	search := ctx.DefaultQuery("search", "")

	// Menghitung offset dan limit berdasarkan parameter pagination

	if page > 0 {
		page = ((page - 1) * pageSize) + 1
	}

	results := []models.Products{}

	// Membuat query untuk pencarian
	query := db.Debug().Preload("Admin").Preload("Variants").Offset(page).Limit(pageSize)

	if search != "" {
		// Menambahkan kondisi pencarian ke query
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Menghitung jumlah total data dengan kondisi pencarian
	var total int64
	if search != "" {
		if err := query.Model(&models.Products{}).Count(&total).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := db.Model(&models.Products{}).Count(&total).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
			return
		}
	}

	// Mengambil data produk dengan pagination dan pencarian
	err = query.Find(&results).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// Menghitung jumlah halaman terakhir
	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))

	// Menghitung halaman saat ini
	currentPage := int((page / pageSize) + 1)

	ctx.JSON(http.StatusOK, gin.H{
		"data": results,
		"pagination": gin.H{
			"last_page": lastPage,
			"limit":     pageSize,
			"offset":    page,
			"page":      currentPage,
			"total":     total,
		},
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

	currentTime := time.Now() // Ambil waktu saat ini

	Product.Created_at = currentTime // Set created_at ke waktu saat ini
	Product.Updated_at = currentTime // Set updated_at ke waktu saat ini

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

	currentTime := time.Now() // Ambil waktu saat ini

	// Update the product record in the database
	updateData := models.Products{
		Name:       productReq.Name,
		Updated_at: currentTime,
	}

	if err := db.Model(&getProduct).Where("uuid = ?", productUUID).Updates(updateData).Error; err != nil {
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
	if err := db.Model(&getProduct).Preload("Admin").Preload("Variants").Where("UUID = ?", productUUID).First(&getProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    getProduct,
		"success": true,
	})

}
