package request

// type VariantRequestCreate struct {
// 	Variant_name string `form:"variant_name" binding:"required"`
// 	Quantity     int    `form:"quantity" binding:"required"`
// 	Product_id   string `form:"product_id" binding:"required"`
// }

type VariantRequest struct {
	Variant_name string `form:"variant_name" binding:"required"`
	Quantity     int    `form:"quantity" binding:"required"`
}
