package router

import (
	"github.com/dopamine-joker/zu_web_server/api/handle"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware(), UserAuthMiddleware(), gin.Recovery())
	r.NoRoute(NoRouteFunc)
	initUserRouter(r)
	initGoodsRouter(r)
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/login", handle.Login)
	userGroup.POST("/register", handle.Register)
	userGroup.POST("/tokenLogin", handle.TokenLogin)
	userGroup.POST("/logout", handle.Logout)
}

func initGoodsRouter(r *gin.Engine) {
	goodsGroup := r.Group("/goods")
	goodsGroup.POST("/upload", handle.Upload)
}

func NoRouteFunc(r *gin.Context) {
	utils.FailWithMsg(r, "please check request url")
}
