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
	"github.com/google/uuid"
)

func GetVariants(ctx *gin.Context) {
	db := database.GetDB()

	// Mengambil parameter "offset" dan "limit" dari permintaan
	page, err := strconv.Atoi(ctx.DefaultQuery("offset", "0")) // Mengatur default offset ke 1 jika tidak ada yang diberikan
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

	results := []models.Variants{}

	// Membuat query untuk pencarian
	query := db.Debug().Preload("Products").Offset(page).Limit(pageSize)

	if search != "" {
		// Menambahkan kondisi pencarian ke query
		query = query.Where("variant_name LIKE ?", "%"+search+"%")
	}

	// Menghitung jumlah total data dengan atau tanpa kondisi pencarian
	var total int64
	if search != "" {
		if err := query.Model(&models.Variants{}).Count(&total).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := db.Model(&models.Variants{}).Count(&total).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": err.Error(),
			})
			return
		}
	}

	// Mengambil data varians dengan pagination dan pencarian
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

func CreateVariant(ctx *gin.Context) {
	db := database.GetDB()

	var variantReq requests.VariantRequest
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var getProduct models.Products
	if err := db.Model(&getProduct).Where("uuid = ?", variantReq.Product_id).First(&getProduct).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request (Product tidak ditemukan)",
			"message": err.Error(),
		})
		return
	}

	Product_ID := uint(getProduct.ID)

	Variant := models.Variants{
		Variant_name: variantReq.Variant_name,
		Quantity:     variantReq.Quantity,
		Product_ID:   Product_ID,
	}

	// Generate a new UUID
	newUUID := uuid.New()
	Variant.UUID = newUUID.String() // Set the generated UUID as the ID

	currentTime := time.Now() // Ambil waktu saat ini

	Variant.Created_at = currentTime // Set created_at ke waktu saat ini
	Variant.Updated_at = currentTime // Set updated_at ke waktu saat ini

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

	variantUUID := ctx.Param("variantUUID")

	// Retrieve existing book from the database
	var getVariant models.Variants
	if err := db.Model(&getVariant).Where("uuid = ?", variantUUID).First(&getVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request (Variant tidak ditemukan)",
			"message": err.Error(),
		})
		return
	}

	// Update the Variant struct with retrieved data

	var variantReq requests.VariantRequestUpdate
	if err := ctx.ShouldBind(&variantReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now() // Ambil waktu saat ini

	// Update the variant record in the database
	updateData := models.Variants{
		Variant_name: variantReq.Variant_name,
		Quantity:     variantReq.Quantity,
		Updated_at:   currentTime,
	}

	if err := db.Model(&getVariant).Where("uuid = ?", variantUUID).Updates(updateData).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": getVariant,
	})
}

func DeleteVariant(ctx *gin.Context) {
	db := database.GetDB()

	Variant := models.Variants{}

	variantUUID := ctx.Param("variantUUID")

	// Retrieve existing product from the database
	var getVariant models.Variants
	if err := db.Model(&getVariant).Where("uuid = ?", variantUUID).First(&getVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	// Delete the book record in the database
	Variant.ID = uint(getVariant.ID)

	deleteData := models.Products{
		ID: Variant.ID,
	}

	if err := db.Model(&Variant).Where("uuid = ?", variantUUID).Delete(deleteData).Error; err != nil {
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

func GetVariantById(ctx *gin.Context) {
	db := database.GetDB()

	contentType := helpers.GetContentType(ctx)
	Variant := models.Variants{}

	variantUUID := ctx.Param("variantUUID")

	if contentType == appJSON {
		ctx.ShouldBindJSON(&Variant)
	} else {
		ctx.ShouldBind(&Variant)
	}

	// Retrieve existing variant from the database
	var getVariant models.Variants
	if err := db.Model(&getVariant).Preload("Products").Where("uuid = ?", variantUUID).First(&getVariant).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    getVariant,
		"success": true,
	})

}
