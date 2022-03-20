package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/db"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
)

const (
	TokenKey   = "X-TOKEN"
	UserInfo   = "X-USER"
	UserId     = "X-UID"
	CountLimit = 20
)

var (
	noVerifyRoute = []string{"/user/login", "/user/register", "/user/tokenLogin", "/user/getSig", "/goods/search", "/metrics"}
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
		uri := c.Request.RequestURI
		// 某些路由不需要检测
		if utils.IsContain(noVerifyRoute, uri) {
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
		c.Set(UserId, user.GetId())
		// redis增加计数
		redisKey := fmt.Sprintf("%s-%s", token, uri)

		exist, err := db.RedisClient.Exists(c.Request.Context(), redisKey).Result()
		if err != nil {
			c.Abort()
			utils.ResponseWithCode(c, misc.CodeFail, "内部数据库错误", nil)
			return
		}
		cnt, err := db.RedisClient.Incr(c.Request.Context(), redisKey).Result()
		if err != nil {
			c.Abort()
			utils.ResponseWithCode(c, misc.CodeFail, "内部数据库错误", nil)
			return
		}
		if cnt > CountLimit {
			c.Abort()
			utils.ResponseWithCode(c, misc.CodeAPILimit, "访问太频繁拉，请稍后再试~", nil)
			return
		}
		// 一开始不存在，设置时间
		if exist <= 0 {
			db.RedisClient.Expire(c.Request.Context(), redisKey, time.Second)
		}
		c.Next()
	}
}
