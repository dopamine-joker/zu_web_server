package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dopamine-joker/zu_web_server/misc"
	"github.com/dopamine-joker/zu_web_server/utils"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

const (
	schema = "etcd"
)

type Resolver struct {
	EtcdAddrs  []string
	DiaTimeout int
	GetTimeout int

	keyPrefix  string
	basePath   string
	serverPath string
	schema     string
	srvAddrs   []resolver.Address

	watchCh clientv3.WatchChan
	closeCh chan struct{}

	client *clientv3.Client
	cc     resolver.ClientConn
}

func NewResolver(etcdAddrs []string, basePath, serverPath string, diaTimeout, getTimeOut int) (*Resolver, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddrs,
		DialTimeout: time.Duration(diaTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Resolver{
		EtcdAddrs:  etcdAddrs,
		DiaTimeout: diaTimeout,
		GetTimeout: getTimeOut,
		basePath:   basePath,
		serverPath: serverPath,
		schema:     schema,
		closeCh:    make(chan struct{}),
		client:     client,
	}, nil

}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	misc.Logger.Info("ResolveNow")
}

func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	r.keyPrefix = r.buildPrefix()
	if err := r.start(); err != nil {
		return nil, err
	}
	return r, nil

}

func (r *Resolver) Scheme() string {
	return r.schema
}

func (r *Resolver) start() error {
	resolver.Register(r)

	if err := r.sync(); err != nil {
		return err
	}

	go r.watch()

	return nil
}

func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Minute)
	r.watchCh = r.client.Watch(context.Background(), r.keyPrefix, clientv3.WithPrefix())
	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				misc.Logger.Error("sync failed", zap.Error(err))
			}

		}
	}
}

func (r *Resolver) update(event []*clientv3.Event) {
	for _, ev := range event {
		switch ev.Type {
		case mvccpb.PUT:
			misc.Logger.Info("etcd key put", zap.String("key", string(ev.Kv.Key)))
			srv, err := r.parseValue(ev.Kv.Value)
			if err != nil {
				misc.Logger.Error("json unmarshal err", zap.Error(err), zap.String("value", string(ev.Kv.Value)))
				continue
			}
			addr := resolver.Address{
				Addr: srv.Addr,
			}
			if !utils.Exist(r.srvAddrs, addr) {
				r.srvAddrs = append(r.srvAddrs, addr)
				_ = r.cc.UpdateState(resolver.State{
					Addresses: r.srvAddrs,
				})
			}
		case mvccpb.DELETE:
			misc.Logger.Info("etcd key delete", zap.String("key", string(ev.Kv.Key)))
			srvAddr := strings.TrimPrefix(string(ev.Kv.Key), r.keyPrefix)
			addr := resolver.Address{
				Addr: srvAddr,
			}
			if s, ok := utils.Remove(r.srvAddrs, addr); ok {
				r.srvAddrs = s
				_ = r.cc.UpdateState(resolver.State{
					Addresses: r.srvAddrs,
				})
			}
		}
	}
}

func (r *Resolver) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.GetTimeout)*time.Second)
	defer cancel()
	res, err := r.client.Get(ctx, r.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, v := range res.Kvs {
		srv, err := r.parseValue(v.Value)
		if err != nil {
			misc.Logger.Error("json Unmarshal data err", zap.Error(err), zap.String("data", string(v.Value)))
		}
		addr := resolver.Address{
			Addr: srv.Addr,
		}
		r.srvAddrs = append(r.srvAddrs, addr)
	}
	if err = r.cc.UpdateState(resolver.State{Addresses: r.srvAddrs}); err != nil {
		return err
	}
	return nil
}

func (r *Resolver) parseValue(value []byte) (RpcLogicServer, error) {
	var err error
	srv := RpcLogicServer{}
	if err = json.Unmarshal(value, &srv); err != nil {
		return srv, err
	}
	return srv, nil
}

func (r *Resolver) buildPrefix() string {
	return fmt.Sprintf("/%s/%s/", r.basePath, r.serverPath)
}
