package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware(), gin.Recovery())
	r.NoRoute(NoRouteFunc)
	return r
}

func NoRouteFunc(r *gin.Context) {
	FailWithMsg(r, "please check request url")
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
