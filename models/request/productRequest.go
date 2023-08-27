package request

import "mime/multipart"

type ProductRequest struct {
	Name  string                `form:"name" binding:"required"`
	Image *multipart.FileHeader `form:"file"`
}
