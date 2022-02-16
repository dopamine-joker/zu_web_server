package handle

import (
	"github.com/dopamine-joker/zu_web_server/api/rpc"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddOrder(c *gin.Context) {
	var form AddOrderForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle add order bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	req := &proto.AddOrderRequest{
		Buyid:  form.BuyId,
		Sellid: form.SellId,
		Gid:    form.GId,
	}

	code, err := rpc.AddOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc add order err", zap.Error(err))
		utils.FailWithMsg(c, "添加失败")
		return
	}

	//data := map[string]interface{}{
	//	"oid": oid,
	//}
	//
	//misc.Logger.Info("add order success", zap.Int32("oid", oid))

	utils.SuccessWithMsg(c, "add order success", nil)
}

func GetBuyOrder(c *gin.Context) {
	var form GetBuyOrderForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle get buy order bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	req := &proto.GetBuyOrderRequest{
		Buyid: form.BuyId,
	}

	code, protoList, err := rpc.GetBuyOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc get buy order err", zap.Error(err))
		utils.FailWithMsg(c, "数据获取失败")
		return
	}

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
			"price":    order.Price,
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
	var form GetSellOrderForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle get buy order bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	req := &proto.GetSellOrderRequest{
		Sellid: form.SellId,
	}

	code, protoList, err := rpc.GetSellOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc get sell order err", zap.Error(err))
		utils.FailWithMsg(c, "数据获取失败")
		return
	}

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
			"price":    order.Price,
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
	var form UpdateOrderForm
	var err error
	if err = c.ShouldBindJSON(&form); err != nil {
		misc.Logger.Error("handle update oreder bind json err", zap.String("err", err.Error()))
		utils.FailWithMsg(c, "参数错误")
		return
	}

	req := &proto.UpdateOrderRequest{
		Id:     form.Id,
		Status: form.Status,
	}
	code, err := rpc.UpdateOrder(c.Request.Context(), req)
	if err != nil || code == misc.CodeFail {
		misc.Logger.Error("rpc update order err", zap.Error(err))
		utils.FailWithMsg(c, "数据获取失败")
		return
	}

	utils.SuccessWithMsg(c, "update order success", nil)
}
