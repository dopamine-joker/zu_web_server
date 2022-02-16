package router

import (
	"github.com/dopamine-joker/zu_web_server/api/handle"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Register() *gin.Engine {
	r := gin.Default()
	r.NoRoute(NoRouteFunc)
	// prometheus
	misc.StartMonitor(r)
	monitor := misc.NewPrometheusMonitor(misc.NAMESPACE, misc.SERVICE)
	r.Use(otelgin.Middleware(misc.SERVICE, otelgin.WithPropagators(otel.GetTextMapPropagator()), otelgin.WithTracerProvider(otel.GetTracerProvider())),
		monitor.PromMiddleware(), CorsMiddleware(), UserAuthMiddleware(), gin.Recovery())
	initUserRouter(r)
	initGoodsRouter(r)
	initOrderRouter(r)
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/login", handle.Login)
	userGroup.POST("/register", handle.Register)
	userGroup.POST("/tokenLogin", handle.TokenLogin)
	userGroup.POST("/logout", handle.Logout)
	userGroup.POST("/getSig", handle.GetSig)
	userGroup.POST("/update", handle.UpdateUser)
	userGroup.POST("/uploadFace", handle.UpdateFace)
}

func initGoodsRouter(r *gin.Engine) {
	goodsGroup := r.Group("/goods")
	goodsGroup.POST("/upload", handle.Upload)
	goodsGroup.POST("/getGoods", handle.GetGoods)
	goodsGroup.POST("/userGoods", handle.GetUserGoodsList)
	goodsGroup.POST("/goodsDetail", handle.GetGoodsDetail)
	goodsGroup.POST("/search", handle.SearchGoods)
	goodsGroup.POST("/delete", handle.DeleteGoods)
}

func initOrderRouter(r *gin.Engine) {
	orderGroup := r.Group("/order")
	orderGroup.POST("/add", handle.AddOrder)
	orderGroup.POST("/getBuy", handle.GetBuyOrder)
	orderGroup.POST("/getSell", handle.GetSellOrder)
	orderGroup.POST("/update", handle.UpdateOrder)
}

func NoRouteFunc(r *gin.Context) {
	utils.FailWithMsg(r, "please check request url")
}
