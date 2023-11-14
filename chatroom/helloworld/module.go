package helloworld

import (
	"fmt"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	this := new(HellWorld)
	return this
}

type HellWorld struct {
	basemodule.BaseModule
}

func (self *HellWorld) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "helloworld"
}
func (self *HellWorld) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}
func (self *HellWorld) OnAppConfigurationLoaded(app module.App) {
	//当App初始化时调用，这个接口不管这个模块是否在这个进程运行都会调用
	self.BaseModule.OnAppConfigurationLoaded(app)
}
func (self *HellWorld) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.GetServer().Options().Metadata["state"] = "alive"
	self.GetServer().RegisterGO("/say/hi", self.say) //handler
	self.GetServer().RegisterGO("HE_test", self.gatesay)
	log.Info("%v模块初始化完成...", self.GetType())
}

func (self *HellWorld) Run(closeSig chan bool) {
	log.Info("%v模块运行中...", self.GetType())
	<-closeSig
	log.Info("%v模块已停止...", self.GetType())
}

func (self *HellWorld) OnDestroy() {
	//一定继承
	self.BaseModule.OnDestroy()
	log.Info("%v模块已回收...", self.GetType())
}
func (self *HellWorld) say(name string) (r string, err error) {
	return fmt.Sprintf("hi %v", name), nil
}

func (self *HellWorld) gatesay(session gate.Session, msg map[string]interface{}) (r string, err error) {
	for {
		fmt.Println("session")
		fmt.Println("ip", session.GetIP())
		fmt.Println("UserID", session.GetUserID())
		fmt.Println("ServerID", session.GetServerID())
		fmt.Println("SessionID", session.GetSessionID())
		fmt.Println("Network", session.GetNetwork())

		bytes, _ := session.Serializable()
		fmt.Println("serializable", string(bytes))
		//var app module.App
		//sessionn, _ := basegate.NewSession(app, bytes)
		//fmt.Println("sessionn")
		//fmt.Println("ip", sessionn.GetIP())
		//fmt.Println("UserID", sessionn.GetUserID())
		//fmt.Println("ServerID", sessionn.GetServerID())
		//fmt.Println("SessionID", sessionn.GetSessionID())
		//fmt.Println("Network", sessionn.GetNetwork())
		//bytess, _ := sessionn.Serializable()
		//sessionn.Push()
		//fmt.Println("serializable", string(bytess))

		// 这个消息用于 下发消息到 客户端
		var input string
		fmt.Scanln(&input)
		session.SendNR("/gate/send/topic/modify", []byte(fmt.Sprintf("%v", input)))
		//sessionn.SendNR("/send/topic/modify", []byte(fmt.Sprintf("%v", input)))
	}
}
