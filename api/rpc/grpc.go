package rpc

import (
	"context"
	"fmt"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/proto"
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
		return
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
		return
	}
	code = response.GetCode()
	return
}

//CheckAuth 检查token rpc调用
func CheckAuth(ctx context.Context, req *proto.CheckAuthRequest) (code int32, authToken string, user *proto.User, err error) {
	response, err := LogicRpcClient.CheckAuth(ctx, req)
	if err != nil {
		return
	}
	code = response.GetCode()
	authToken = response.GetAuthToken()
	user = response.GetUser()
	return
}

func TokenLogin(ctx context.Context, req *proto.TokenLoginRequest) (code int32, token string, user *proto.User, err error) {
	response, err := LogicRpcClient.TokenLogin(ctx, req)
	if err != nil {
		return
	}
	code = response.GetCode()
	token = response.GetAuthToken()
	user = response.GetUser()
	return
}

func Logout(ctx context.Context, req *proto.LogoutRequest) (code int32, err error) {
	response, err := LogicRpcClient.Logout(ctx, req)
	if err != nil {
		return
	}
	code = response.GetCode()
	return
}

//UploadPic 上传文件,文件以byte数组形式上传
func UploadPic(ctx context.Context, req *proto.UploadRequest) (code int32, err error) {
	response, err := LogicRpcClient.UploadPic(ctx, req)
	if err != nil {
		return
	}
	code = response.GetCode()
	return
}
