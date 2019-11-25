package table

import (
	"fmt"
	"github.com/sniperHW/kendynet"
	"github.com/sniperHW/kendynet/message"
)

var (
	sessionMap = map[string]*Session{} // 所有的链接
	// session在status状态下，不能发送出去的消息。
	messageTocFilter map[SessStatus]map[string]struct{}
	// session在status状态下，需要处理的消息。
	messageTosFilter map[SessStatus]map[string]struct{}
)

type SessStatus int

const (
	None SessStatus = iota
	Look
	Editor
)

type Session struct {
	session  kendynet.StreamSession
	Status   SessStatus
	Table    *Table
	UserName string
}

func NewSession(session kendynet.StreamSession, file *Table, name string) *Session {
	return &Session{
		session:  session,
		Status:   Editor,
		Table:    file,
		UserName: name,
	}
}

func (this *Session) RemoteAddr() string {
	return this.session.RemoteAddr().String()
}

func (this *Session) Send(cmd string, msg []byte) error {
	v := messageTocFilter[this.Status]
	if _, ok := v[cmd]; !ok {
		return this.DirectSend(msg)
	}
	return nil
}

func (this *Session) DirectSend(msg []byte) error {
	return this.session.SendMessage(message.NewWSMessage(message.WSTextMessage, msg))
}

func (this *Session) SetStatus(status SessStatus) {
	this.Status = status
}

func OnClose(sess kendynet.StreamSession, reason string) {
	if session, ok := sessionMap[sess.RemoteAddr().String()]; ok {
		fmt.Printf("user:%s onclose %s\n", session.UserName, reason)
		session.Table.RemoveSession(session)
		session.Table = nil
		delete(sessionMap, session.RemoteAddr())
	} else {
		fmt.Printf("session:%s onclose %s\n", sess.RemoteAddr().String(), reason)
	}
}

func init() {
	// 不能发送出去的消息
	messageTocFilter = map[SessStatus]map[string]struct{}{
		Look: {
			"cellSelected": {},
			"insertRow":    {},
			"removeRow":    {},
			"insertCol":    {},
			"removeCol":    {},
			//"setCellValues": {},
			"saveTable": {},
			"pushAll":   {},
			"rollback":  {},
		},
		Editor: {
			//"rollback":    {},
		},
	}

}
