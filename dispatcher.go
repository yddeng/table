package table

import (
	"encoding/json"
	"fmt"
	"github.com/sniperHW/kendynet"
	"github.com/sniperHW/kendynet/message"
	"github.com/yddeng/table/pgsql"
)

var dispatcher map[string]func(msg map[string]interface{}, session *Session)

func Dispatcher(msg map[string]interface{}, session kendynet.StreamSession) {
	fmt.Println("Dispatcher", msg)

	cmdI := msg["cmd"]
	if cmdI == nil {
		fmt.Println("no cmd")
		return
	}

	cmd := cmdI.(string)
	if cmd == "openTable" {
		onOpenTable(msg, session)
	} else {
		sess := sessionMap[session.RemoteAddr().String()]
		if sess == nil {
			pushErr(cmdI, "no session", NewSession(session, nil, ""))
			return
		}
		handler, ok := dispatcher[cmd]
		if ok {
			handler(msg, sess)
		} else {
			panic(fmt.Sprintf("cmd:%s no register", cmd))
		}
	}
}

func onOpenTable(req map[string]interface{}, session kendynet.StreamSession) {
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
	sessionMap[sess.RemoteAddr()] = sess
	table.AddSession(sess)

	b := Must(json.Marshal(map[string]interface{}{
		"cmd":       "openTable",
		"tableName": tableName,
		"userName":  userName,
		"version":   table.version,
		"data":      getAll(table.tmpFile),
	})).([]byte)
	_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
}

func handleCellSelected(req map[string]interface{}, session *Session) {
	session.Table.pushOtherUser(map[string]interface{}{
		"cmd":      req["cmd"],
		"selected": req["selected"],
		"userName": session.UserName,
	}, session.UserName)
}

func handleInsertRow(req map[string]interface{}, session *Session) {
	session.Table.AddCmd(req, session.UserName)
	doCmd(session.Table.tmpFile, req, false)
	session.Table.pushAllSession(req)
	session.Table.pushAll()
}

func handleRemoveRow(req map[string]interface{}, session *Session) {
	session.Table.AddCmd(req, session.UserName)
	doCmd(session.Table.tmpFile, req, false)
	session.Table.pushAllSession(req)
	session.Table.pushAll()
}

func handleInsertCol(req map[string]interface{}, session *Session) {
	session.Table.AddCmd(req, session.UserName)
	doCmd(session.Table.tmpFile, req, false)
	session.Table.pushAllSession(req)
	session.Table.pushAll()
}

func handleRemoveCol(req map[string]interface{}, session *Session) {
	session.Table.AddCmd(req, session.UserName)
	doCmd(session.Table.tmpFile, req, false)
	session.Table.pushAllSession(req)
	session.Table.pushAll()
}

func handleSetCellValues(req map[string]interface{}, session *Session) {
	session.Table.AddCmd(req, session.UserName)
	doCmd(session.Table.tmpFile, req, false)
	session.Table.pushAllSession(req)
	session.Table.pushAll()
}

func handleSaveTable(req map[string]interface{}, session *Session) {
	err := session.Table.SaveTable()
	if err != nil {
		pushErr(req["cmd"], err.Error(), session)
		return
	}
	b, _ := json.Marshal(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": session.Table.version,
		"data":    getAll(session.Table.xlFile),
	})
	for _, session := range this.sessionMap {
		_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
	}
}

func handleLookHistory(req map[string]interface{}, session *Session) {
	ver := int(req["version"].(float64))
	table := session.Table

	data := getAll(table.xlFile)
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

	b, _ := json.Marshal(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": ver,
		"data":    getAll(retF),
	})
	session.Send(req["cmd"].(string), b)
}

func handleBackEditor(req map[string]interface{}, session *Session) {
	data, _ := getAll(session.Table.tmpFile)
	b, _ := json.Marshal(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": session.Table.version,
		"data":    data,
	})
	session.Send(req["cmd"].(string), b)
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
	fmt.Println("pushErr", msg)
	b, err := json.Marshal(map[string]interface{}{
		"cmd":    "pushErr",
		"doCmd":  cmd,
		"errMsg": msg,
	})
	if err == nil {
		session.Send("pushErr", b)
	}
}

func init() {
	dispatcher = map[string]func(msg map[string]interface{}, session *Session){}
	// 只做转发
	dispatcher["cellSelected"] = handleCellSelected
	// cmd，操作文件命令
	dispatcher["insertRow"] = handleInsertRow
	dispatcher["removeRow"] = handleRemoveRow
	dispatcher["insertCol"] = handleInsertCol
	dispatcher["removeCol"] = handleRemoveCol
	dispatcher["setCellValues"] = handleSetCellValues
	// 保存，将所有cmd指令存库，生成一条版本记录
	dispatcher["saveTable"] = handleSaveTable
	// 查看历史版本，返回编辑状态。消息转发过滤
	dispatcher["lookHistory"] = handleLookHistory
	dispatcher["backEditor"] = handleBackEditor
	// 回滚版本，单独生成一条版本记录
	dispatcher["rollback"] = handleRollback
}
