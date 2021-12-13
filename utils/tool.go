package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/resolver"

	"github.com/dopamine-joker/zu_web_server/misc"
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

//GetRpcMsg 提取rpc调用的错误信息
func GetRpcMsg(errMsg string) string {
	sli := strings.Split(errMsg, "desc = ")
	if len(sli) < 1 {
		return ""
	}
	return sli[1]
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
