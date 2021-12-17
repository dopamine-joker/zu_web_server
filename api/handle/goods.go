package handle

import (
	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"strconv"
)

const (
	uploadKey = "files"
)

//Upload 该请求数据格式不为json,而为multipart/form-data
func Upload(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		misc.Logger.Error("Upload multipartForm err", zap.Error(err))
		utils.FailWithMsg(c, "请求出错")
		return
	}

	// 提取除图片外的其他参数
	dataMap := make(map[string]interface{})
	for key, vals := range form.Value {
		dataMap[key] = vals[0]
	}

	var uploadForm UploadForm
	if err = mapstructure.Decode(dataMap, &uploadForm); err != nil {
		misc.Logger.Error("upload decode struct err", zap.Error(err))
		utils.FailWithMsg(c, "请求参数错误")
		return
	}

	// 提取文件,转换为byte数组后保存
	files := make(map[string][]byte)
	for key, headers := range form.File {
		if key != uploadKey {
			continue
		}
		// grpc发送文件
		for _, file := range headers { //遍历每一个文件
			// 获取文件的数据流
			src, err := file.Open()
			if err != nil {
				misc.Logger.Error("pic open file err", zap.Error(err))
				utils.FailWithMsg(c, "图片解码出现问题")
			}
			// 读取其中的byte流
			picBytes := make([]byte, 4*1024)
			n, err := src.Read(picBytes)
			if err != nil {
				misc.Logger.Error("pic file read err", zap.Error(err))
				utils.FailWithMsg(c, "图片解码出现问题")
			}

			// 存储到map中,用于grpc,根据读取的字节数裁断
			files[file.Filename] = picBytes[:n]
			_ = src.Close()
		}
	}
	var picList []*proto.PicStream
	for name, bytes := range files {
		picList = append(picList, &proto.PicStream{
			Name:    name,
			Content: bytes,
		})
	}

	uidInt32, err := strconv.ParseInt(uploadForm.Uid, 10, 32)
	if err != nil {
		misc.Logger.Error("parse userid err", zap.String("uid", uploadForm.Uid))
		utils.FailWithMsg(c, "用户id解析出错")
	}

	//构造请求
	req := &proto.UploadRequest{
		Uid:     int32(uidInt32),
		Name:    uploadForm.Name,
		Price:   uploadForm.Price,
		Detail:  uploadForm.Detail,
		PicList: picList,
	}

	code, err := rpc.UploadPic(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc upload err", zap.Error(err))
		utils.FailWithMsg(c, utils.GetRpcMsg(err.Error()))
		return
	}

	misc.Logger.Info("upload success")

	utils.SuccessWithMsg(c, "upload success", nil)
}
