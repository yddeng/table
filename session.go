package table

import (
	"encoding/json"
	"fmt"
	"github.com/sniperHW/kendynet"
	"github.com/yddeng/table/conf"
	"path"
)

var (
	sessionMap = map[string]*Session{}
)

type Session struct {
	kendynet.StreamSession
	File     *ExcelFile
	UserName string
	doEvent  []map[string]interface{}
}

func NewSession(session kendynet.StreamSession, file *ExcelFile, name string) *Session {
	return &Session{
		StreamSession: session,
		File:          file,
		UserName:      name,
		doEvent:       []map[string]interface{}{},
	}
}

func OnClose(sess kendynet.StreamSession, reason string) {
	if session, ok := sessionMap[sess.RemoteAddr().String()]; ok {
		fmt.Println("onclose", reason)
		session.File.RemoveSession(session)
		session.File = nil
		delete(sessionMap, sess.RemoteAddr().String())
	}
}

func onOpenFile(req map[string]interface{}, session kendynet.StreamSession) {
	fmt.Println("handleOpenFile", req)

	fileName := req["fileName"].(string)
	userName := req["userName"].(string)
	ef, ok := excelFiles[fileName]
	if !ok {
		_conf := conf.GetConfig()
		ef, _ = OpenExcel(path.Join(_conf.ExcelPath, fileName))
		excelFiles[fileName] = ef
	}

	sess := NewSession(session, ef, userName)
	sessionMap[sess.RemoteAddr().String()] = sess
	ef.AddSession(sess)
	ef.PushData()
}

func (this *Session) addEvent(event map[string]interface{}) {
	this.doEvent = append(this.doEvent, event)
}

func (this *Session) SaveEvent() {
	events, err := json.Marshal(this.doEvent)
	if err != nil {
		logger.Infoln(events)
	}
}
