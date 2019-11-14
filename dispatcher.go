package table

import (
	"encoding/json"
	"fmt"
	"github.com/sniperHW/kendynet"
	"github.com/sniperHW/kendynet/message"
	"github.com/yddeng/table/pgsql"
)

func Dispatcher(msg map[string]interface{}, session kendynet.StreamSession) {
	cmd := msg["cmd"].(string)
	fmt.Println("Dispatcher", msg)

	if cmd == "openTable" {
		onOpenTable(msg, session)
	} else {
		sess := checkSession(session)
		if sess == nil {
			pushErr(cmd, "Session is nil", NewSession(session, nil, ""))
			return
		}

		switch cmd {
		case "saveTable":
			handleSaveTable(msg, sess)
		case "rollback":
			handleRollback(msg, sess)
		case "lookHistory":
			handleLookHistory(msg, sess)
		case "cellSelected":
			handleCellSelected(msg, sess)
		default:
			sess.Table.AddCmd(msg, sess.UserName)
			doCmd(sess.Table.tmpFile, msg, false)
			sess.Table.pushAll()
		}
	}
}

func checkSession(sess kendynet.StreamSession) *Session {
	return sessionMap[sess.RemoteAddr().String()]
}

func onOpenTable(req map[string]interface{}, session kendynet.StreamSession) {
	fmt.Println("handleOpenTable", req)

	var err error
	tableName := req["tableName"].(string)
	userName := req["userName"].(string)
	table, ok := tables[tableName]
	if !ok {
		table, err = OpenTable(tableName)
		if err != nil {
			pushErr("openTable", err.Error(), NewSession(session, nil, userName))
			return
		}
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

func handleCellSelected(req map[string]interface{}, session *Session) {
	selected := req["selected"].([]interface{})
	table := session.Table

	b, _ := json.Marshal(map[string]interface{}{
		"cmd":      "pushCellSelected",
		"selected": selected,
	})
	for _, sess := range table.sessionMap {
		if session.UserName != sess.UserName {
			_ = sess.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
		}
	}
}

func handleSaveTable(req map[string]interface{}, session *Session) {
	err := session.Table.SaveTable()
	if err != nil {
		pushErr("saveTable", err.Error(), session)
	}
}

func handleLookHistory(req map[string]interface{}, session *Session) {
	ver := int(req["version"].(float64))
	table := session.Table

	data, _ := getAll(table.xlFile)
	retF := newFile()
	cloneFile(retF, data)

	if ver < table.version {
		for i := table.version; i > ver; i-- {
			ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr("lookHistory", err.Error(), session)
				return
			}
			rollbackCmds(retF, ret)
		}
	} else if ver > table.version {
		_, err := pgsql.LoadCmd(table.tableName, ver)
		if err != nil {
			pushErr("lookHistory", "版本号错误", session)
			return
		}
		for i := table.version; i <= ver; i++ {
			ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr("lookHistory", err.Error(), session)
				return
			}
			doCmds(retF, ret)
		}
	}

	retData, _ := getAll(retF)
	b, _ := json.Marshal(map[string]interface{}{
		"cmd":     "lookHistory",
		"version": ver,
		"data":    retData,
	})
	_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))

}

func handleRollback(req map[string]interface{}, session *Session) {
	now := int(req["now"].(float64))
	ver := int(req["goto"].(float64))
	table := session.Table

	if now != table.version {
		pushErr("rollback", "版本号不一致，不能回退", session)
		table.pushAll()
		return
	}

	if len(table.sessionMap) > 1 {
		pushErr("rollback", "多人操作，不能回退", session)
		return
	}

	if ver < table.version {
		for i := table.version; i > ver; i-- {
			ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr("rollback", err.Error(), session)
				return
			}
			rollbackCmds(table.tmpFile, ret)
		}
	} else if ver > table.version {
		_, err := pgsql.LoadCmd(table.tableName, ver)
		if err != nil {
			pushErr("rollback", "版本号错误", session)
			return
		}
		for i := table.version; i <= ver; i++ {
			ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr("rollback", err.Error(), session)
				return
			}
			doCmds(table.tmpFile, ret)
		}
	}

	data, _ := getAll(table.tmpFile)
	table.version = ver
	cloneFile(table.xlFile, data)
	table.SaveTable()
	b, _ := json.Marshal(map[string]interface{}{
		"cmd":     "rollback",
		"version": ver,
		"data":    data,
	})
	_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))

}

//## pushError 返回错误信息
//```
//cmd : "pushError"
//doCmd:  string
//errMsg: string
//```
func pushErr(cmd, msg string, session *Session) {
	b, err := json.Marshal(map[string]interface{}{
		"cmd":    "pushErr",
		"doCmd":  cmd,
		"errMsg": msg,
	})
	if err == nil {
		_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
	}
}
