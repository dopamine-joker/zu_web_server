package handle

import (
	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"log"
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
	code, token, user, err := rpc.Login(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc login err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}
	misc.Logger.Info("login success", zap.String("email", loginForm.Email), zap.String("token", token))
	span.SetAttributes(
		attribute.String("email", req.Email),
		attribute.String("token", token),
		attribute.Int64("code", int64(code)),
	)

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

	code, token, user, err := rpc.TokenLogin(c.Request.Context(), req)
	if code == misc.CodeFail || token == "" || err != nil {
		misc.Logger.Error("rpc tokenLogin err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}

	span.SetAttributes(
		attribute.String("token", req.Token),
		attribute.Int64("code", int64(code)),
	)

	misc.Logger.Info("tokenLogin success", zap.String("token", token))

	dataMap := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	utils.SuccessWithMsg(c, "token login success", dataMap)
}

func UpdateUser(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	defer span.End()
	var form UpdateUserForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("update user bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.UpdateUserRequest{
		Uid:      uid,
		Email:    form.Email,
		Name:     form.Name,
		Phone:    form.Phone,
		Password: form.Password,
		School:   form.School,
		Sex:      form.Sex,
	}

	code, err := rpc.UpdateUser(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc update user err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}

	span.SetAttributes(
		attribute.Int64("id", int64(req.GetUid())),
		attribute.String("email", req.GetEmail()),
		attribute.String("phone", req.GetPhone()),
		attribute.String("name", req.GetName()),
		attribute.String("password", req.GetPassword()),
		attribute.String("school", req.GetSchool()),
		attribute.Int64("sex", int64(req.GetSex())),
		attribute.Int64("code", int64(code)),
	)

	utils.SuccessWithMsg(c, "update user success", nil)
}

func Register(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())

	defer span.End()
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

	span.SetAttributes(attribute.String("email", req.Email),
		attribute.String("name", req.Name),
		attribute.String("password", req.Password),
	)

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
	span := trace.SpanFromContext(c.Request.Context())

	defer span.End()
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

	span.SetAttributes(
		attribute.String("token", req.Token),
		attribute.Int64("code", int64(code)),
	)

	misc.Logger.Info("logout success", zap.String("token", logoutForm.Token))

	utils.SuccessWithMsg(c, "logout success", nil)
}

//GetSig 获取sdk初始化的签名
func GetSig(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var getSigForm GetSigForm
	var err error
	if err = c.ShouldBindJSON(&getSigForm); err != nil {
		misc.Logger.Error("handle getsig bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, err.Error())
		return
	}

	//uid, err := utils.GetContextUserId(c)
	//if err != nil {
	//	misc.Logger.Error("请求Token参数错误")
	//	utils.FailWithMsg(c, err.Error())
	//}

	code, sig, err := rpc.GetSig(c.Request.Context(), getSigForm.UserId, getSigForm.SdkAppId, getSigForm.Expire)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc getSIg err", zap.Error(err))
		utils.FailWithMsg(c, "获取签名失败")
		return
	}

	span.SetAttributes(
		attribute.String("sig", sig),
		attribute.Int64("code", int64(code)),
	)

	dataMap := map[string]interface{}{
		"sig": sig,
	}

	utils.SuccessWithMsg(c, "getSig success", dataMap)
}

const (
	uploadFaceKey = "face"
)

func UpdateFace(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	form, err := c.MultipartForm()
	if err != nil {
		misc.Logger.Error("Upload face err", zap.Error(err))
		utils.FailWithMsg(c, "请求出错")
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	dataMap := make(map[string]interface{})
	for key, vals := range form.Value {
		dataMap[key] = vals[0]
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
				picBytes := make([]byte, 10*1024*1024)
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

	var pic *proto.FileStream
	for name, bytes := range face {
		pic = &proto.FileStream{
			Name:    name,
			Content: bytes,
		}
	}

	req := &proto.UploadFaceRequest{
		Uid: uid,
		Pic: pic,
	}

	code, path, err := rpc.UploadFace(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc upload face err", zap.Error(err))
		utils.FailWithMsg(c, "上传失败")
		return
	}

	span.SetAttributes(
		attribute.String("path", path),
		attribute.Int64("code", int64(code)),
	)

	misc.Logger.Info("upload face success")

	res := map[string]interface{}{
		"path": path,
	}

	utils.SuccessWithMsg(c, "upload pic success", res)
}
