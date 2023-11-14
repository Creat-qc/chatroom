package main

import (
	"chatroom/gate"
	"chatroom/httpgateway"
	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/registry"
	"github.com/liangdas/mqant/registry/consul"
	server2 "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"net/http"
	"time"
)

func main() {

	// 启动监听 6060端口
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	gate.JudgeGuest = func(session gate.Session) bool {
		if session.GetUserID() != "" {
			return false
		}
		return true
	}
	// 启动 consul 实现服务发现
	rs := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	// 启动 nats 准备中间件连接
	opts := &server2.Options{Host: "127.0.0.1", Port: 4222}
	ns, err := server2.NewServer(opts)

	if err != nil {
		panic(err)
	}
	go ns.Start()
	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}
	// 连接 中间件
	nc, err := nats.Connect("nats://127.0.0.1:4222", nats.MaxReconnects(10000))
	if err != nil {
		log.Error("nats error %v", err)
		return
	}
	app := mqant.CreateApp(
		module.KillWaitTTL(1*time.Minute),
		module.Debug(true),
		module.Nats(nc),
		module.Registry(rs),
		module.RegisterTTL(20*time.Second),
		module.RegisterInterval(10*time.Second),
	)
	err = app.Run(
		mgate.Module(),
		httpgateway.Module(),
	)

	if err != nil {
		log.Error(err.Error())
	}
}
