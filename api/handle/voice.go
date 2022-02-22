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
)

const (
	uploadVoiceKey = "voice"
)

func VoiceProcess(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

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

	face := make(map[string][]byte)

	for key, headers := range form.File {
		if key == uploadVoiceKey {
			for _, file := range headers {
				src, err := file.Open()
				if err != nil {
					misc.Logger.Error("voiceFile open file err", zap.Error(err))
					utils.FailWithMsg(c, "音频解码出现问题")
				}
				picBytes := make([]byte, 10*1024*1024)
				n, err := src.Read(picBytes)
				if err != nil {
					misc.Logger.Error("voiceFile file read err", zap.Error(err))
					utils.FailWithMsg(c, "音频解码出现问题")
				}
				face[file.Filename] = picBytes[:n]
				_ = src.Close()
			}
		}
	}

	var voiceFile *proto.FileStream
	for name, bytes := range face {
		voiceFile = &proto.FileStream{
			Name:    name,
			Content: bytes,
		}
	}

	req := &proto.VoiceToTxtRequest{
		VoiceFile: voiceFile,
	}

	code, txt, err := rpc.VoiceToTxt(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc upload face err", zap.Error(err))
		utils.FailWithMsg(c, "上传失败")
		return
	}

	span.SetAttributes(
		attribute.String("txt", txt),
		attribute.Int64("code", int64(code)),
	)

	misc.Logger.Info("voice to txt success", zap.String("txt", txt))

	res := map[string]interface{}{
		"txt": txt,
	}

	utils.SuccessWithMsg(c, "upload voiceFile success", res)
}
