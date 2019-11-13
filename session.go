package table

import (
	"encoding/json"
	"fmt"
	"github.com/sniperHW/kendynet"
	"github.com/sniperHW/kendynet/message"
	"github.com/yddeng/table/pgsql"
	"sync"
)

var (
	sessionMap = map[string]*Session{}
)

type Session struct {
	sync.Mutex
	kendynet.StreamSession
	Table    *Table
	UserName string
	doCmds   []map[string]interface{}
}

func NewSession(session kendynet.StreamSession, file *Table, name string) *Session {
	return &Session{
		StreamSession: session,
		Table:         file,
		UserName:      name,
		doCmds:        []map[string]interface{}{},
	}
}

func OnClose(sess kendynet.StreamSession, reason string) {
	if session, ok := sessionMap[sess.RemoteAddr().String()]; ok {
		fmt.Println("onclose", reason)
		session.Table.RemoveSession(session)
		if len(session.doCmds) > 0 {
			rollbackCmds(session.Table.tmpFile, session.doCmds)
			cmdsStr, err := json.Marshal(session.doCmds)
			if err == nil {
				logger.Infoln("rollback", string(cmdsStr))
			}
			session.Table.PushData()
		}
		session.Table = nil
		delete(sessionMap, sess.RemoteAddr().String())
	}
}

func onOpenTable(req map[string]interface{}, session kendynet.StreamSession) {
	fmt.Println("handleOpenTable", req)

	tableName := req["tableName"].(string)
	userName := req["userName"].(string)
	table, ok := tables[tableName]
	if !ok {
		table, _ = OpenTable(tableName)
		tables[tableName] = table
	}

	sess := NewSession(session, table, userName)
	sessionMap[sess.RemoteAddr().String()] = sess
	table.AddSession(sess)

	resp := map[string]interface{}{
		"cmd":       "openTable",
		"tableName": tableName,
		"userName":  userName,
		"version":   table.version,
	}
	data, _ := getAll(table.tmpFile)
	resp["data"] = data
	b, _ := json.Marshal(resp)
	session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
}

func (this *Session) addCmd(cmd map[string]interface{}) {
	this.doCmds = append(this.doCmds, cmd)
}

func (this *Session) SaveCmd() {
	if len(this.doCmds) > 0 {
		cmdsStr, err := json.Marshal(this.doCmds)
		if err == nil {
			logger.Infoln(string(cmdsStr))
			v, err := pgsql.InsertCmd(this.Table.fileName, this.UserName, string(cmdsStr))
			if err == nil {
				this.Table.version = v
				this.Table.SaveTable()
			}
		}
		this.doCmds = []map[string]interface{}{}
	}
}
