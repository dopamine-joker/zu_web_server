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

type UpdateUserForm struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	School   string `form:"school" json:"school" binding:"required"`
	Sex      int32  `form:"sex" json:"sex" binding:"required"`
}

type LogoutForm struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type GetSigForm struct {
	SdkAppId int `form:"sdkAppId" json:"sdkAppId" binding:"required"`
	Expire   int `form:"expire" json:"expire" binding:"required"`
}

type UploadForm struct {
	Name   string `form:"name" json:"token" binding:"required"`
	Price  string `form:"price" json:"price" binding:"required"`
	Type   int32  `form:"type" json:"type" binding:"required"`
	School string `form:"school" json:"school" binding:"required"`
	Detail string `form:"detail" json:"detail" binding:"required"`
}

type GetGoodsForm struct {
	Page  *int32 `form:"page" json:"page" binding:"required"`   //这里使用指针才能参与0值
	Count *int32 `form:"count" json:"count" binding:"required"` //同理
}

type PicListForm struct {
	Gid *int32 `form:"gid" json:"gid" binding:"required"`
}

type SearchGoodsForm struct {
	GName string `form:"gname" json:"gname" binding:"required"`
}

type DeleteGoodsForm struct {
	Gid int32 `form:"gid" json:"gid" binding:"required"`
}

type AddOrderForm struct {
	SellId int32  `form:"sellid" json:"sellid" binding:"required"`
	GId    int32  `form:"gid" json:"gid" binding:"required"`
	School string `form:"school" json:"school" binding:"required"`
}

type UpdateOrderForm struct {
	Id     int32 `form:"id" json:"id" binding:"required"`
	Status int32 `form:"status" json:"status" binding:"required"`
}

type AddFavoritesForm struct {
	GId int32 `form:"gid" json:"gid" binding:"required"`
}

type DeleteFavoritesForm struct {
	FId int32 `form:"fid" json:"fid" binding:"required"`
}

type AddCommentForm struct {
	GId     int32  `form:"gid" json:"gid" binding:"required"`
	OId     int32  `form:"oid" json:"oid" binding:"required"`
	Level   int32  `form:"level" json:"level" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

type GetCommentByGoodsIdForm struct {
	GId int32 `form:"gid" json:"gid" binding:"required"`
}

type DeleteCommentForm struct {
	CId int32 `form:"cid" json:"cid" binding:"required"`
}
