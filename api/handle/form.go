package handle

type LoginForm struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterForm struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
}

type TokenLoginForm struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type LogoutForm struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type UploadForm struct {
	Uid    string `form:"uid" json:"uid" binding:"required"`
	Name   string `form:"name" json:"token" binding:"required"`
	Price  string `form:"price" json:"price" binding:"required"`
	Detail string `form:"detail" json:"detail" binding:"required"`
}
