package router

import (
	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	TokenKey = "X-TOKEN"
	UserInfo = "X-USER"
)

var (
	noVerifyRoute = []string{"/user/login", "/user/register", "/user/tokenLogin", "/user/getSig"}
)

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

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 某些路由不需要检测
		if utils.IsContain(noVerifyRoute, c.Request.RequestURI) {
			c.Next()
			return
		}
		// 获取header的token
		token := c.GetHeader(TokenKey)
		if token == "" {
			c.Abort()
			utils.ResponseWithCode(c, misc.CodeTokenError, nil, nil)
			return
		}
		// 验证token
		req := &proto.CheckAuthRequest{AuthToken: token}
		code, _, user, err := rpc.CheckAuth(c.Request.Context(), req)
		if code == misc.CodeFail || err != nil {
			c.Abort()
			utils.ResponseWithCode(c, misc.CodeTokenError, nil, nil)
			return
		}
		c.Set(UserInfo, user)
		c.Next()
	}
}
