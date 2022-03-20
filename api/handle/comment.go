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

func AddComment(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var form AddCommentForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle add comment bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
	}

	req := &proto.AddCommentRequest{
		Uid:     uid,
		Gid:     form.GId,
		Oid:     form.OId,
		Level:   form.Level,
		Content: form.Content,
	}

	code, cid, err := rpc.AddComment(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc add comment err", zap.Error(err))
		utils.FailWithMsg(c, "添加失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("userId", int64(uid)),
		attribute.Int64("goodsId", int64(form.GId)),
		attribute.Int64("orderId", int64(form.OId)),
		attribute.Int64("level", int64(form.Level)),
		attribute.String("content", form.Content),
		attribute.Int64("code", int64(code)),
	)

	data := map[string]interface{}{
		"cid": cid,
	}

	utils.SuccessWithMsg(c, "add comment success", data)
}

func DeleteComment(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var form DeleteCommentForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle delete comment bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
	}

	req := &proto.DeleteCommentRequest{
		Uid: uid,
		Cid: form.CId,
	}

	code, err := rpc.DeleteComment(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc delete comment err", zap.Error(err))
		utils.FailWithMsg(c, "删除失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("commentId", int64(form.CId)),
		attribute.Int64("code", int64(code)),
	)

	utils.SuccessWithMsg(c, "delete comment success", nil)
}

func GetCommentByUserId(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
	}

	req := &proto.GetCommentByUserIdRequest{
		Uid: uid,
	}

	code, protoList, err := rpc.GetCommentByUserId(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc get user comment err", zap.Error(err))
		utils.FailWithMsg(c, "获取评论失败")
		return
	}

	var list []map[string]interface{}

	for _, comment := range protoList {
		list = append(list, map[string]interface{}{
			"id":        comment.Id,
			"uid":       comment.Uid,
			"gid":       comment.Gid,
			"oid":       comment.Oid,
			"content":   comment.Content,
			"level":     comment.Level,
			"time":      comment.Time,
			"goodsName": comment.Name,
			"price":     comment.Price,
			"cover":     comment.Cover,
		})
	}

	data := map[string]interface{}{
		"len":  len(list),
		"data": list,
	}

	span.SetAttributes(
		attribute.Int64("userId", int64(uid)),
		attribute.Int64("code", int64(code)),
	)

	utils.SuccessWithMsg(c, "get user comments success", data)
}

func GetCommentByGoodsId(c *gin.Context) {
	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var form GetCommentByGoodsIdForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle get goods comment bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	req := &proto.GetCommentByGoodsIdRequest{
		Gid: form.GId,
	}

	code, protoList, err := rpc.GetCommentByGoodsId(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc get goods comment err", zap.Error(err))
		utils.FailWithMsg(c, "获取评论失败")
		return
	}

	var list []map[string]interface{}

	for _, comment := range protoList {
		list = append(list, map[string]interface{}{
			"id":       comment.Id,
			"uid":      comment.Uid,
			"gid":      comment.Gid,
			"oid":      comment.Oid,
			"content":  comment.Content,
			"level":    comment.Level,
			"time":     comment.Time,
			"userName": comment.Uname,
			"userFace": comment.Uface,
		})
	}

	data := map[string]interface{}{
		"len":  len(list),
		"data": list,
	}

	span.SetAttributes(
		attribute.Int64("goodsId", int64(form.GId)),
		attribute.Int64("code", int64(code)),
	)

	utils.SuccessWithMsg(c, "get goods comments success", data)
}
