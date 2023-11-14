package chat

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	basemodule "github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	chat := new(Chat)
	return chat
}

type Chat struct {
	basemodule.BaseModule
}

func (c *Chat) Version() string {
	return "v0.0.1"
}

func (c *Chat) GetType() string {
	return "chat"
}

func (c *Chat) OnInit(app module.App, settings *conf.ModuleSettings) {
	c.BaseModule.OnInit(c, app, settings)
	log.Info("%v初始化完成", c.GetType())
}

func (c *Chat) Run(closeSig chan bool) {
	log.Info("%v模块运行中", c.GetType())
	<-closeSig
	log.Info("%v模块已终止", c.GetType())
}

func (c *Chat) SendMsg() {
	scanln, err := fmt.Scanln()
	if err != nil {
		return
	}

	log.Info("%v模块发送消息%v至中间件", c.GetType(), scanln)
}

func (c *Chat) RecMsg() {
	log.Info("%v模块接收中间件消息", c.GetType())

}
