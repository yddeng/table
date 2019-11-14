package table

import (
	"fmt"
	"github.com/sniperHW/kendynet"
)

var (
	sessionMap = map[string]*Session{}
)

type Session struct {
	kendynet.StreamSession
	Table    *Table
	UserName string
}

func NewSession(session kendynet.StreamSession, file *Table, name string) *Session {
	return &Session{
		StreamSession: session,
		Table:         file,
		UserName:      name,
	}
}

func OnClose(sess kendynet.StreamSession, reason string) {
	if session, ok := sessionMap[sess.RemoteAddr().String()]; ok {
		fmt.Println("onclose", reason)
		session.Table.RemoveSession(session)
		session.Table = nil
		delete(sessionMap, sess.RemoteAddr().String())
	}
}
