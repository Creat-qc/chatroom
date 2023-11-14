package mgate

import (
	"github.com/liangdas/mqant/gate"
	"github.com/liangdas/mqant/log"
)

var MySessionLearner = func() gate.SessionLearner {
	s := new(SL)
	return s
}

type SL struct {
}

func (S SL) Connect(session gate.Session) {
	// 使用sessionID 的值 表示 UserID
	session.Bind(session.GetSessionID())
	log.Info("新客户端连接 ID %v IP %v", session.GetSessionID(), session.GetIP())
}

func (S SL) DisConnect(session gate.Session) {
	log.Info("新客户端释放 ID %v IP %v", session.GetSessionID(), session.GetIP())
}
