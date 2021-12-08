package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FailWithMsg(c *gin.Context, msg interface{}) {
	ResponseWithCode(c, CodeFail, msg, nil)
}

func ResponseWithCode(c *gin.Context, msgCode int, msg interface{}, data interface{}) {
	if msg == nil {
		if val, ok := MsgCodeMap[msgCode]; ok {
			msg = val
		} else {
			msg = MsgCodeMap[CodeUnknownError]
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    msgCode,
		"message": msg,
		"data":    data,
	})
}
