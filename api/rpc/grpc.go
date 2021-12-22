package rpc

import (
	"context"
	"fmt"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

var LogicRpcClient proto.RpcLogicServiceClient

func InitLogicRpcClient() {
	r, err := NewResolver(misc.Conf.EtcdCfg.Host, misc.Conf.EtcdCfg.BasePath, misc.Conf.EtcdCfg.ServerPathLogic, 5, 5)
	if err != nil {
		misc.Logger.Error("NewResolver err", zap.Error(err))
		panic(err)
	}
	resolver.Register(r)
	conn, err := grpc.Dial(fmt.Sprintf("%s://author/%s/%s", r.Scheme(), misc.Conf.EtcdCfg.BasePath, misc.Conf.EtcdCfg.ServerPathLogic),
		grpc.WithBalancerName(roundrobin.Name), grpc.WithInsecure(), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	if err != nil {
		misc.Logger.Error("fail to dial", zap.Error(err))
		panic(err)
	}
	LogicRpcClient = proto.NewRpcLogicServiceClient(conn)
}

//Login api grpc调用login
func Login(ctx context.Context, req *proto.LoginRequest) (code int32, authToken string, user *proto.User, err error) {
	misc.Logger.Info("login call rpc", zap.Any("request", req))
	response, err := LogicRpcClient.Login(ctx, req)
	if err != nil {
		return misc.CodeFail, "", nil, err
	}
	code = response.GetCode()
	authToken = response.GetAuthToken()
	user = response.GetUser()
	return
}

//Register api rpc调用Register
func Register(ctx context.Context, req *proto.RegisterRequest) (code int32, err error) {
	response, err := LogicRpcClient.Register(ctx, req)
	if err != nil {
		return misc.CodeFail, err
	}
	code = response.GetCode()
	return
}

//CheckAuth 检查token rpc调用
func CheckAuth(ctx context.Context, req *proto.CheckAuthRequest) (code int32, authToken string, user *proto.User, err error) {
	response, err := LogicRpcClient.CheckAuth(ctx, req)
	if err != nil {
		return misc.CodeFail, "", nil, err
	}
	code = response.GetCode()
	authToken = response.GetAuthToken()
	user = response.GetUser()
	return
}

func TokenLogin(ctx context.Context, req *proto.TokenLoginRequest) (code int32, token string, user *proto.User, err error) {
	response, err := LogicRpcClient.TokenLogin(ctx, req)
	if err != nil {
		return misc.CodeFail, "", nil, err
	}
	code = response.GetCode()
	token = response.GetAuthToken()
	user = response.GetUser()
	return
}

func Logout(ctx context.Context, req *proto.LogoutRequest) (code int32, err error) {
	response, err := LogicRpcClient.Logout(ctx, req)
	if err != nil {
		return misc.CodeFail, err
	}
	code = response.GetCode()
	return
}

//GetSig im sdk获取sig
func GetSig(ctx context.Context, userId string, sdkAppId, expire int) (code int32, sig string, err error) {
	sig, err = tencentyun.GenUserSig(sdkAppId, misc.Key, userId, expire)
	if err != nil {
		return misc.CodeFail, "", err
	}
	code = misc.CodeSuccess
	return
}

//UploadGoods 上传文件,文件以byte数组形式上传
func UploadGoods(ctx context.Context, req *proto.UploadRequest) (code int32, err error) {
	response, err := LogicRpcClient.UploadPic(ctx, req)
	if err != nil {
		return misc.CodeFail, err
	}
	code = response.GetCode()
	return
}

//GetGoods 获取商品基本信息
func GetGoods(ctx context.Context, req *proto.GetGoodsRequest) (code int32, goodsList []*proto.Goods, err error) {
	response, err := LogicRpcClient.GetGoods(ctx, req)
	if err != nil {
		return misc.CodeFail, nil, err
	}
	code = response.Code
	goodsList = response.GoodsList
	return
}

//PicList 根据商品id获取对应图片
func PicList(ctx context.Context, req *proto.GetGoodsDetailRequest) (code int32, goodsDetail *proto.GoodsDetail, picList []*proto.Pic, err error) {
	response, err := LogicRpcClient.GetGoodsPic(ctx, req)
	if err != nil {
		return misc.CodeFail, nil, nil, err
	}
	code = response.Code
	goodsDetail = response.GetGoods()
	picList = response.PicList
	return
}
