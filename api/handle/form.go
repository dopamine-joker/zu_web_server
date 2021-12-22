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

type GetSigForm struct {
	UserId   string `form:"userId" json:"userId" binding:"required"`
	SdkAppId int    `form:"sdkAppId" json:"sdkAppId" binding:"required"`
	Expire   int    `form:"expire" json:"expire" binding:"required"`
}

type UploadForm struct {
	Uid    string `form:"uid" json:"uid" binding:"required"`
	Name   string `form:"name" json:"token" binding:"required"`
	Price  string `form:"price" json:"price" binding:"required"`
	Detail string `form:"detail" json:"detail" binding:"required"`
}

type GetGoodsForm struct {
	Page  *int32 `form:"page" json:"page" binding:"required"`   //这里使用指针才能参与0值
	Count *int32 `form:"count" json:"count" binding:"required"` //同理
}

type PicListForm struct {
	Gid *int32 `form:"gid" json:"gid" binding:"required"`
}
