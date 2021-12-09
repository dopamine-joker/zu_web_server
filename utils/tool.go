package utils

import (
	"github.com/dopamine-joker/zu_web_server/misc"
	"google.golang.org/grpc/resolver"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessWithMsg(c *gin.Context, msg interface{}, data interface{}) {
	ResponseWithCode(c, misc.CodeSuccess, msg, data)
}

func FailWithMsg(c *gin.Context, msg interface{}) {
	ResponseWithCode(c, misc.CodeFail, msg, nil)
}

func ResponseWithCode(c *gin.Context, msgCode int, msg interface{}, data interface{}) {
	if msg == nil {
		if val, ok := misc.MsgCodeMap[msgCode]; ok {
			msg = val
		} else {
			msg = misc.MsgCodeMap[misc.CodeUnknownError]
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    msgCode,
		"message": msg,
		"data":    data,
	})
}

// Remove helper function
func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

func Exist(srvAddrs []resolver.Address, addr resolver.Address) bool {
	for _, srvAddr := range srvAddrs {
		if srvAddr.Addr == addr.Addr {
			return true
		}
	}
	return false
}

func IsContain(list []string, str string) bool {
	for _, e := range list {
		if e == str {
			return true
		}
	}
	return false
}
