package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"zu_web_server/misc"
)

var UserList map[string]string

func Login(c *gin.Context) {
	var loginForm LoginForm
	var err error
	if err = c.ShouldBindJSON(&loginForm); err != nil {
		misc.Logger.Error("handler login bind json err", zap.String("err", err.Error()))
		misc.FailWithMsg(c, err.Error())
		return
	}

	log.Println("login: email:", loginForm.Email, "pwd:", loginForm.Password)
	if _, ok := UserList[loginForm.Email]; !ok {
		misc.FailWithMsg(c, "no this user")
		return
	}

	if UserList[loginForm.Email] != loginForm.Password {
		misc.FailWithMsg(c, "Incorrect password")
		return
	}

	misc.Logger.Info("login success", zap.String("email", loginForm.Email),
		zap.String("password", loginForm.Password))
	misc.SuccessWithMsg(c, "login success", nil)
}

func Register(c *gin.Context) {
	var registerForm RegisterForm
	if err := c.ShouldBindJSON(&registerForm); err != nil {
		misc.Logger.Error("handler register bind json err", zap.String("err", err.Error()))
		misc.FailWithMsg(c, err.Error())
		return
	}

	log.Println("register: email:", registerForm.Email, "pwd:", registerForm.Password)

	if _, ok := UserList[registerForm.Email]; ok {
		misc.FailWithMsg(c, "this user have been register")
		return
	}

	UserList[registerForm.Email] = registerForm.Password

	misc.Logger.Info("register success", zap.String("email", registerForm.Email),
		zap.String("password", registerForm.Password))

	misc.SuccessWithMsg(c, "register success", nil)
}
