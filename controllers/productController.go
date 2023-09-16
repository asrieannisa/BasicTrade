package controllers

import (
	"BasicTrade/database"
	"BasicTrade/helpers"
	models "BasicTrade/models/entity"
	requests "BasicTrade/models/request"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		q := r.URL.Query()
// 		page, _ := strconv.Atoi(q.Get("offset"))
// 		if page <= 0 {
// 			page = 1
// 		}

// 		pageSize, _ := strconv.Atoi(q.Get("limit"))
// 		switch {
// 		case pageSize > 100:
// 			pageSize = 100
// 		case pageSize <= 0:
// 			pageSize = 10
// 		}

// 		offset := (page - 1) * pageSize
// 		return db.Offset(offset).Limit(pageSize)
// 	}
// }

// func GetProducts(ctx *gin.Context) {
// 	db := database.GetDB()

// 	// Menggunakan fungsi Paginate untuk menerapkan pagination
// 	db = Paginate(ctx.Request)(db)

// 	results := []models.Products{}

// 	err := db.Debug().Preload("Admin").Preload("Variants").Find(&results).Error
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Bad request",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	q := r.URL.Query()
// 	page, _ := strconv.Atoi(q.Get("offset"))

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"data": results,
// 		"pagination": gin.H{
// 			"last_page": lastPage,
// 			"limit":     pageSize,
// 			"offset":    offset,
// 			"page":      page,
// 			"total":     total,
// 		},
// 	})
// }

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
	offset := 0
	if page > 1 {
		offset = (page - 1) * pageSize
	} else if page == 1 {
		offset = 1
	}
	limit := pageSize

	results := []models.Products{}

	// Membuat query untuk pencarian
	query := db.Debug().Preload("Admin").Preload("Variants").Offset(offset).Limit(limit)

	if search != "" {
		// Menambahkan kondisi pencarian ke query
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Menghitung jumlah total data dengan kondisi pencarian
	var total int64
	if err := query.Model(&models.Products{}).Count(&total).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
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

	if currentPage > lastPage {
		currentPage = lastPage
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": results,
		"pagination": gin.H{
			"last_page": lastPage,
			"limit":     pageSize,
			"offset":    offset,
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
