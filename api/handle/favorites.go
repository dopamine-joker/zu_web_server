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

func AddFavorites(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var form AddFavoritesForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle add favorites bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.AddFavoritesRequest{
		Uid: uid,
		Gid: form.GId,
	}

	code, fid, err := rpc.AddFavorites(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc add favorites err", zap.Error(err))
		utils.FailWithMsg(c, "添加失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("goodsId", int64(form.GId)),
		attribute.Int64("code", int64(code)),
	)

	data := map[string]interface{}{
		"fid": fid,
	}

	utils.SuccessWithMsg(c, "add favorites success", data)
}

func DeleteFavorites(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var form DeleteFavoritesForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle delete favorites bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.DeleteFavoritesRequest{
		Uid: uid,
		Fid: form.FId,
	}

	code, err := rpc.DeleteFavorites(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc delete favorites err", zap.Error(err))
		utils.FailWithMsg(c, "添加失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("favoritesId", int64(form.FId)),
		attribute.Int64("code", int64(code)),
	)

	utils.SuccessWithMsg(c, "delete favorites success", nil)
}

func GetUserFavorites(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.GetUserFavoritesRequest{
		Uid: uid,
	}

	code, protoList, err := rpc.GetUserFavorites(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc get user favorites err", zap.Error(err))
		utils.FailWithMsg(c, "添加失败")
		return
	}

	var list []map[string]interface{}

	for _, favorites := range protoList {
		list = append(list, map[string]interface{}{
			"id":    favorites.Id,
			"uid":   favorites.Uid,
			"gid":   favorites.Gid,
			"gname": favorites.Name,
			"price": favorites.Price,
			"cover": favorites.Cover,
		})
	}

	data := map[string]interface{}{
		"len":  len(list),
		"data": list,
	}

	span.SetAttributes(
		attribute.Int64("code", int64(code)),
	)

	utils.SuccessWithMsg(c, "get user favorites success", data)
}
