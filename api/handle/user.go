package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
	"strconv"

	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
)

func Login(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

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
	span.SetAttributes(attribute.String("email", req.GetEmail()))
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

	span := trace.SpanFromContext(c.Request.Context())

	defer span.End()

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

	span.SetAttributes(attribute.String("token", req.GetToken()))

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

func UpdateUser(c *gin.Context) {
	var form UpdateUserForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("update user bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.UpdateUserRequest{
		Uid:      form.Id,
		Email:    form.Email,
		Name:     form.Name,
		Password: form.Password,
	}

	code, err := rpc.UpdateUser(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc update user err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}

	utils.SuccessWithMsg(c, "update user success", nil)
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

const (
	uploadFaceKey = "face"
)

func UpdateFace(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		misc.Logger.Error("Upload face err", zap.Error(err))
		utils.FailWithMsg(c, "请求出错")
		return
	}

	dataMap := make(map[string]interface{})
	for key, vals := range form.Value {
		dataMap[key] = vals[0]
	}
	var uploadForm UploadFaceForm
	if err = mapstructure.Decode(dataMap, &uploadForm); err != nil {
		misc.Logger.Error("upload decode struct err", zap.Error(err))
		utils.FailWithMsg(c, "请求参数错误")
		return
	}

	face := make(map[string][]byte)

	for key, headers := range form.File {
		if key == uploadFaceKey {
			for _, file := range headers {
				src, err := file.Open()
				if err != nil {
					misc.Logger.Error("pic open file err", zap.Error(err))
					utils.FailWithMsg(c, "图片解码出现问题")
				}
				picBytes := make([]byte, 4*1024)
				n, err := src.Read(picBytes)
				if err != nil {
					misc.Logger.Error("pic file read err", zap.Error(err))
					utils.FailWithMsg(c, "图片解码出现问题")
				}
				face[file.Filename] = picBytes[:n]
				_ = src.Close()
			}
		}
	}

	var pic *proto.PicStream
	for name, bytes := range face {
		pic = &proto.PicStream{
			Name:    name,
			Content: bytes,
		}
	}

	uidInt32, err := strconv.ParseInt(uploadForm.Id, 10, 32)
	if err != nil {
		misc.Logger.Error("parse userid err", zap.String("uid", uploadForm.Id))
		utils.FailWithMsg(c, "用户id解析出错")
	}

	req := &proto.UploadFaceRequest{
		Uid: int32(uidInt32),
		Pic: pic,
	}

	code, path, err := rpc.UploadFace(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc upload face err", zap.Error(err))
		utils.FailWithMsg(c, "上传失败")
		return
	}

	misc.Logger.Info("upload face success")

	res := map[string]interface{}{
		"path": path,
	}

	utils.SuccessWithMsg(c, "upload pic success", res)
}
