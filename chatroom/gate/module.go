package mgate

import (
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	basegate "github.com/liangdas/mqant/gate/base"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/mitchellh/mapstructure"
)

var Module = func() module.Module {
	gate := new(Gate)
	return gate
}

type Gate struct {
	basegate.Gate
}

func (g *Gate) Version() string {
	return "v0.0.1"
}

func (g *Gate) GetType() string {
	return "gate"
}

func (g *Gate) OnInit(app module.App, settings *conf.ModuleSettings) {
	g.Gate.OnInit(g, app, settings,
		gate.WsAddr(":3653"),
		gate.TCPAddr(":3563"),
		gate.SetStorageHandler(MySessionStorage()),
		gate.SetSessionLearner(MySessionLearner()),
	)
	g.GetServer().RegisterGO("HD_chat", g.FaceChat)
	log.Info("%v模块初始化完成", g.GetType())
}

func (g *Gate) FaceChat(session gate.Session, msg map[string]interface{}) (r string, err error) {
	log.Info("接收到请求了")
	faceChat := FaceChat{}
	// 将请求参数的map interface类型 绑定在结构体上
	erro := mapstructure.Decode(msg, &faceChat)
	if erro != nil {
		log.Error(erro.Error())
	}
	// 使用App 需要是 gate 已经配置过的App
	app := g.Gate.App
	bytes, _ := g.GetStorageHandler().Query(faceChat.UserID)
	session_a, _ := basegate.NewSession(app, bytes)

	session_a.SendNR("/chat/face", []byte(faceChat.Data))
	return "nil", nil
}
