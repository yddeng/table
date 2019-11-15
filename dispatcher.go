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
		case "lookHistory":
			handleLookHistory(msg, sess)
		case "backEditor":
			handleBackEditor(msg, sess)
		case "rollback":
			handleRollback(msg, sess)
		case "cellSelected":
			sess.Table.pushOtherUser(map[string]interface{}{
				"cmd":      msg["cmd"],
				"selected": msg["selected"],
				"userName": sess.UserName,
			}, sess.UserName)
		default:
			sess.Table.AddCmd(msg, sess.UserName)
			doCmd(sess.Table.tmpFile, msg, false)
			sess.Table.pushAllSession(msg)
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
			pushErr(req["cmd"], err.Error(), NewSession(session, nil, userName))
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

func handleSaveTable(req map[string]interface{}, session *Session) {
	err := session.Table.SaveTable()
	if err != nil {
		pushErr(req["cmd"], err.Error(), session)
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
				pushErr(req["cmd"], err.Error(), session)
				return
			}
			rollbackCmds(retF, ret)
		}
	} else if ver > table.version {
		_, err := pgsql.LoadCmd(table.tableName, ver)
		if err != nil {
			pushErr(req["cmd"], "版本号错误", session)
			return
		}
		for i := table.version; i <= ver; i++ {
			ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr(req["cmd"], err.Error(), session)
				return
			}
			doCmds(retF, ret)
		}
	}

	retData, _ := getAll(retF)
	b, _ := json.Marshal(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": ver,
		"data":    retData,
	})
	_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
}

func handleBackEditor(req map[string]interface{}, session *Session) {
	data, _ := getAll(session.Table.tmpFile)
	b, _ := json.Marshal(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": session.Table.version,
		"data":    data,
	})
	_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
}

func handleRollback(req map[string]interface{}, session *Session) {
	ver := int(req["version"].(float64))
	table := session.Table

	/*
		if len(table.sessionMap) > 1 {
			pushErr(req["cmd"], "多人操作，不能回退", session)
			return
		}
	*/

	if ver < table.version {
		for i := table.version; i > ver; i-- {
			ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr(req["cmd"], err.Error(), session)
				return
			}
			rollbackCmds(table.xlFile, ret)
		}
	} else if ver > table.version {
		_, err := pgsql.LoadCmd(table.tableName, ver)
		if err != nil {
			pushErr(req["cmd"], "版本号错误", session)
			return
		}
		for i := table.version; i <= ver; i++ {
			ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr(req["cmd"], err.Error(), session)
				return
			}
			doCmds(table.xlFile, ret)
		}
	}

	data, _ := getAll(table.xlFile)

	// todo 保存版本指令
	//cmdsStr, _ := json.Marshal(this.cmds)
	//users, _ := json.Marshal(this.cmdUsers)
	//v, err := pgsql.InsertCmd(this.tableName, string(users), string(cmdsStr))
	//if err != nil {
	//	return err
	//}

	// 更新数据库
	b, _ := json.Marshal(data)
	err := pgsql.UpdateTableData(table.tableName, table.version, string(b))
	if err != nil {
		pushErr(req["cmd"], err.Error(), session)

		// 错误回退
		tmpData, _ := getAll(table.tmpFile)
		cloneFile(table.xlFile, tmpData)
		return
	}

	// 当前表状态更改
	tmpFile := newFile()
	cloneFile(tmpFile, data)
	table.tmpFile = tmpFile
	table.cmds = []map[string]interface{}{}
	table.version = ver

	// 同步给所有人
	table.pushAllSession(map[string]interface{}{
		"cmd":     "rollback",
		"version": ver,
		"data":    data,
	})
}

//## pushError 返回错误信息
//```
//cmd : "pushError"
//doCmd:  string
//errMsg: string
//```
func pushErr(cmd interface{}, msg string, session *Session) {
	b, err := json.Marshal(map[string]interface{}{
		"cmd":    "pushErr",
		"doCmd":  cmd,
		"errMsg": msg,
	})
	if err == nil {
		_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
	}
}
