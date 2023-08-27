package request

type AdminRequest struct {
	Email    string `form:"email" binding:"required"`
	Name     string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
}
