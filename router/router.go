package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zu_web_server/handler"
	"zu_web_server/misc"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware(), gin.Recovery())
	r.NoRoute(NoRouteFunc)
	initUserRouter(r)
	return r
}

func initUserRouter(r *gin.Engine) {
	handler.UserList = make(map[string]string)
	userGroup := r.Group("/user")
	userGroup.POST("/login", handler.Login)
	userGroup.POST("/register", handler.Register)
}

func NoRouteFunc(r *gin.Context) {
	misc.FailWithMsg(r, "please check request url")
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
		c.Set("content-type", "application/json")
		method := c.Request.Method
		// options 用于获取url所支持的方法，"GET,POST..."
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
		}
		c.Next()
	}
}
