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

func AddOrder(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var form AddOrderForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle add order bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.AddOrderRequest{
		Buyid:  uid,
		Sellid: form.SellId,
		Gid:    form.GId,
		School: form.School,
	}

	code, err := rpc.AddOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc add order err", zap.Error(err))
		utils.FailWithMsg(c, "添加失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("buyId", int64(uid)),
		attribute.Int64("sellId", int64(form.SellId)),
		attribute.Int64("goodId", int64(form.GId)),
		attribute.String("school", form.School),
		attribute.Int64("code", int64(code)),
	)

	utils.SuccessWithMsg(c, "add order success", nil)
}

func GetBuyOrder(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.GetBuyOrderRequest{
		Buyid: uid,
	}

	code, protoList, err := rpc.GetBuyOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc get buy order err", zap.Error(err))
		utils.FailWithMsg(c, "数据获取失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("buyId", int64(uid)),
		attribute.Int64("code", int64(code)),
	)

	var list []map[string]interface{}

	for _, order := range protoList {
		list = append(list, map[string]interface{}{
			"id":       order.Id,
			"buyid":    order.Buyid,
			"buyName":  order.BuyName,
			"sellid":   order.Sellid,
			"sellName": order.SellName,
			"gid":      order.GId,
			"gName":    order.Gname,
			"school":   order.School,
			"price":    order.Price,
			"type":     order.Type,
			"cover":    order.Cover,
			"status":   order.Status,
			"time":     order.Time,
		})
	}

	dataMap := map[string]interface{}{
		"data": list,
		"len":  len(list),
	}

	misc.Logger.Info("get buy order success", zap.Any("data", dataMap))

	utils.SuccessWithMsg(c, "get buy order success", dataMap)
}

func GetSellOrder(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.GetSellOrderRequest{
		Sellid: uid,
	}

	code, protoList, err := rpc.GetSellOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc get sell order err", zap.Error(err))
		utils.FailWithMsg(c, "数据获取失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("sellId", int64(uid)),
		attribute.Int64("code", int64(code)),
	)

	var list []map[string]interface{}

	for _, order := range protoList {
		list = append(list, map[string]interface{}{
			"id":       order.Id,
			"buyid":    order.Buyid,
			"buyName":  order.BuyName,
			"sellid":   order.Sellid,
			"sellName": order.SellName,
			"gid":      order.GId,
			"gName":    order.Gname,
			"school":   order.School,
			"price":    order.Price,
			"type":     order.Type,
			"cover":    order.Cover,
			"status":   order.Status,
			"time":     order.Time,
		})
	}

	dataMap := map[string]interface{}{
		"data": list,
		"len":  len(list),
	}

	misc.Logger.Info("get buy order success", zap.Any("data", dataMap))

	utils.SuccessWithMsg(c, "get buy order success", dataMap)
}

func UpdateOrder(c *gin.Context) {

	span := trace.SpanFromContext(c.Request.Context())
	defer span.End()

	var form UpdateOrderForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle update oreder bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	uid, err := utils.GetContextUserId(c)
	if err != nil {
		misc.Logger.Error("请求Token参数错误")
		utils.FailWithMsg(c, err.Error())
		return
	}

	req := &proto.UpdateOrderRequest{
		Id:     form.Id,
		Uid:    uid,
		Status: form.Status,
	}

	code, err := rpc.UpdateOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc update order err", zap.Error(err))
		utils.FailWithMsg(c, "数据获取失败")
		return
	}

	span.SetAttributes(
		attribute.Int64("orderId", int64(req.Id)),
		attribute.Int64("status", int64(req.Status)),
	)

	utils.SuccessWithMsg(c, "update order success", nil)
}
