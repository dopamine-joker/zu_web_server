package handle

import (
	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"log"
	"strconv"
)

const (
	uploadKey      = "files"
	uploadCoverKey = "cover"
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
	cover := make(map[string][]byte)
	for key, headers := range form.File {
		if key != uploadKey && key != uploadCoverKey {
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

			if key == uploadKey {
				// 存储到map中,用于grpc,根据读取的字节数裁断
				files[file.Filename] = picBytes[:n]
			} else if key == uploadCoverKey {
				// 存储到封面byte,用于grpc,根据读取的字节数裁断
				cover[file.Filename] = picBytes[:n]
			}
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

	var coverPic *proto.PicStream
	for name, bytes := range cover {
		coverPic = &proto.PicStream{
			Name:    name,
			Content: bytes,
		}
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
		Cover:   coverPic,
		PicList: picList,
	}

	code, err := rpc.UploadGoods(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc upload err", zap.Error(err))
		utils.FailWithMsg(c, "该用户已存在相同名称的物品了")
		return
	}

	misc.Logger.Info("upload success")

	utils.SuccessWithMsg(c, "upload success", nil)
}

func GetGoods(c *gin.Context) {
	var getGoodsForm GetGoodsForm
	var err error
	if err = c.ShouldBindJSON(&getGoodsForm); err != nil {
		misc.Logger.Error("handle get goods bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	req := &proto.GetGoodsRequest{
		Page:  *getGoodsForm.Page,
		Count: *getGoodsForm.Count,
	}

	code, list, err := rpc.GetGoods(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc getGoods err", zap.Error(err))
		utils.FailWithMsg(c, "获取物品列表失败")
		return
	}

	var dataMap []map[string]interface{}
	for _, goods := range list {
		m := make(map[string]interface{})
		m["id"] = goods.Id
		m["name"] = goods.Name
		m["uname"] = goods.Uname
		m["price"] = goods.Price
		m["cover"] = goods.Cover
		dataMap = append(dataMap, m)
	}

	res := map[string]interface{}{
		"len":   len(dataMap),
		"goods": dataMap,
	}

	log.Println(res)

	utils.SuccessWithMsg(c, "get goods list success", res)
}

func GetGoodsDetail(c *gin.Context) {
	var picListForm PicListForm
	var err error
	if err = c.ShouldBindJSON(&picListForm); err != nil {
		misc.Logger.Error("handle get goods picList bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	req := &proto.GetGoodsDetailRequest{
		Gid: *picListForm.Gid,
	}

	code, goodsDetail, list, err := rpc.PicList(c.Request.Context(), req)
	if code == misc.CodeFail || err != nil {
		misc.Logger.Error("rpc getGoods picList err", zap.Error(err))
		utils.FailWithMsg(c, "拉取物品图片失败")
		return
	}

	var picList []map[string]interface{}
	for _, p := range list {
		data := map[string]interface{}{
			"id":   p.GetPid(),
			"path": p.GetPath(),
		}
		picList = append(picList, data)
	}

	dataMap := map[string]interface{}{
		"data": map[string]interface{}{
			"gid":         goodsDetail.Gid,
			"uid":         goodsDetail.Uid,
			"name":        goodsDetail.Name,
			"uname":       goodsDetail.Uname,
			"price":       goodsDetail.Price,
			"detail":      goodsDetail.Detail,
			"cover":       goodsDetail.Cover,
			"create_time": goodsDetail.CreateTime,
			"picList":     picList,
		},
	}

	log.Println(dataMap)
	misc.Logger.Info("get pic list success", zap.Int32("gid", req.GetGid()))

	utils.SuccessWithMsg(c, "get pic list success", dataMap)
}
