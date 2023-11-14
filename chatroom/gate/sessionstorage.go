package mgate

import (
	"encoding/json"
	"errors"
	"github.com/liangdas/mqant/gate"
	basegate "github.com/liangdas/mqant/gate/base"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	"os"
	"time"
)

var Path = "./gate/sessionstorage.json"

var MySessionStorage = func() gate.StorageHandler {
	s := new(St)
	return s
}

type MySession struct {
	SessionID  string
	Serialize  string
	ExpireTime string
}

type St struct {
}

func (s St) Storage(session gate.Session) (err error) {
	bytes, _ := session.Serializable()
	sessio := MySession{session.GetSessionID(), string(bytes), time.Now().Add(30 * time.Minute).Format(time.RFC3339)}

	sessionsList, _ := Read()
	sessionsList = append(sessionsList, sessio)

	_ = Write(sessionsList)
	return
}

func (s St) Delete(session gate.Session) (err error) {
	sessionID := session.GetSessionID()

	sessionsList, _ := Read()
	var index int
	for i, s := range sessionsList {
		log.Info("已存储的session %v", s)
		if s.SessionID == sessionID {
			index = i
		}
	}
	newSessionsList := append(sessionsList[:index], sessionsList[index+1])

	_ = Write(newSessionsList)
	return
}

func (s St) Query(Userid string) (data []byte, err error) {
	var res []byte
	erro := errors.New("userid not found")
	sessionsList, _ := Read()
	for _, s := range sessionsList {
		var app module.App
		session, _ := basegate.NewSession(app, []byte(s.Serialize))
		if session.GetUserID() == Userid {
			return []byte(s.Serialize), nil
		}
	}
	return res, erro
}

func (s St) Heartbeat(session gate.Session) {
	// TODO 不清楚 Heartbeat是如何触发的
	sessionsList, _ := Read()
	for _, s := range sessionsList {
		if session.GetSessionID() == s.SessionID {
			s.ExpireTime = time.Now().Add(30 * time.Minute).Format(time.RFC3339)
			break
		}
	}
	_ = Write(sessionsList)
}

func Write(sessions []MySession) (err error) {
	jsonFile, erro := os.Create(Path)
	if erro != nil {
		log.Info("session文件打开异常 %v", erro)
		return erro
	}
	defer jsonFile.Close()

	encode := json.NewEncoder(jsonFile)
	err_a := encode.Encode(sessions)
	if err_a != nil {
		log.Info("session编码异常 %v", err_a)
		return err_a
	}
	return
}

func Read() (sessions []MySession, err error) {
	var sessionList []MySession

	jsonFile, erro := os.Open(Path)
	if erro != nil {
		log.Info("session文件打开异常 %v", erro)
		return sessionList, erro
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	erro_a := decoder.Decode(&sessionList)
	if erro_a != nil {
		log.Info("session解码异常 %v", erro_a)
		return sessionList, erro_a
	}
	return sessionList, nil
}
