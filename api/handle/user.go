package handle

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"

	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
)

func Login(c *gin.Context) {
	var loginForm LoginForm
	var err error
	if err = c.ShouldBindJSON(&loginForm); err != nil {
		misc.Logger.Error("handle login bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	log.Println("login: email:", loginForm.Email, "pwd:", loginForm.Password)
	misc.Logger.Info("web_server login", zap.String("email", loginForm.Email))
	req := &proto.LoginRequest{
		Email:    loginForm.Email,
		Password: loginForm.Password,
	}
	code, token, user, err := rpc.Login(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc login err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}
	misc.Logger.Info("login success", zap.String("email", loginForm.Email), zap.String("token", token))

	dataMap := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	utils.SuccessWithMsg(c, "login success", dataMap)
}

func TokenLogin(c *gin.Context) {
	var tokenLoginForm TokenLoginForm
	if err := c.ShouldBindJSON(&tokenLoginForm); err != nil {
		misc.Logger.Error("handle tokenLogin bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	log.Println("token login: token:", tokenLoginForm.Token)
	misc.Logger.Info("web_server tokenLogin", zap.String("token", tokenLoginForm.Token))

	req := &proto.TokenLoginRequest{
		Token: tokenLoginForm.Token,
	}

	code, token, user, err := rpc.TokenLogin(c.Request.Context(), req)
	if code == misc.CodeFail || token == "" || err != nil {
		misc.Logger.Error("rpc tokenLogin err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}

	misc.Logger.Info("tokenLogin success", zap.String("token", token))

	dataMap := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	utils.SuccessWithMsg(c, "token login success", dataMap)
}

func Register(c *gin.Context) {
	var registerForm RegisterForm
	if err := c.ShouldBindJSON(&registerForm); err != nil {
		misc.Logger.Error("handle register bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	misc.Logger.Info("web_server register", zap.String("email", registerForm.Email),
		zap.String("pwd", registerForm.Password))

	req := &proto.RegisterRequest{
		Email:    registerForm.Email,
		Name:     registerForm.Name,
		Password: registerForm.Password,
	}

	code, err := rpc.Register(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc register err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}

	misc.Logger.Info("register success", zap.String("email", registerForm.Email),
		zap.String("password", registerForm.Password))

	utils.SuccessWithMsg(c, "register success", nil)
}

func Logout(c *gin.Context) {
	var logoutForm LogoutForm
	if err := c.ShouldBindJSON(&logoutForm); err != nil {
		misc.Logger.Error("handle logout bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	log.Println("logout, token:", logoutForm.Token)
	misc.Logger.Info("web_server logout", zap.String("token", logoutForm.Token))

	req := &proto.LogoutRequest{
		Token: logoutForm.Token,
	}

	code, err := rpc.Logout(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc logout err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}

	misc.Logger.Info("logout success", zap.String("token", logoutForm.Token))

	utils.SuccessWithMsg(c, "logout success", nil)
}

//GetSig 获取sdk初始化的签名
func GetSig(c *gin.Context) {
	var getSigForm GetSigForm
	var err error
	if err = c.ShouldBindJSON(&getSigForm); err != nil {
		misc.Logger.Error("handle getsig bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	code, sig, err := rpc.GetSig(c.Request.Context(), getSigForm.UserId, getSigForm.SdkAppId, getSigForm.Expire)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc getSIg err", zap.Error(err))
		utils.FailWithMsg(c, "获取签名失败")
		return
	}

	dataMap := map[string]interface{}{
		"sig": sig,
	}

	utils.SuccessWithMsg(c, "getSig success", dataMap)
}
