package table

import (
	"encoding/json"
	"fmt"
	"github.com/sniperHW/kendynet"
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

	b := MustJsonMarshal(map[string]interface{}{
		"cmd":       "openTable",
		"tableName": tableName,
		"userName":  userName,
		"version":   table.version,
		"data":      getAll(table.tmpFile),
	})
	sess.DirectSend(b)
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
	session.Table.pushAll()
}

func handleSaveTable(req map[string]interface{}, session *Session) {
	err := session.Table.SaveTable()
	if err != nil {
		pushErr(req["cmd"], err.Error(), session)
		return
	}

	session.Table.pushAllSession(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": session.Table.version,
		"data":    getAll(session.Table.xlFile),
	})

}

func handleVersionList(req map[string]interface{}, session *Session) {
	b := MustJsonMarshal(map[string]interface{}{
		"cmd":  req["cmd"],
		"list": session.Table.packHistory(),
	})
	session.Send(req["cmd"].(string), b)

}

func handleLookHistory(req map[string]interface{}, session *Session) {
	ver := int(req["version"].(float64))
	table := session.Table
	session.SetStatus(Look)

	data := getAll(table.xlFile)
	newF := newFile(data)

	if ver <= table.version {
		for i := table.version; i > ver; i-- {
			_, _, ret, err := pgsql.LoadCmd(table.tableName, i)
			if err != nil {
				pushErr(req["cmd"], err.Error(), session)
				return
			}
			rollbackCmds(newF, ret)
		}
	} else {
		pushErr(req["cmd"], "版本号错误", session)
		return
	}

	b := MustJsonMarshal(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": ver,
		"data":    getAll(newF),
	})
	session.Send(req["cmd"].(string), b)
}

func handleBackEditor(req map[string]interface{}, session *Session) {
	session.SetStatus(Editor)
	b := MustJsonMarshal(map[string]interface{}{
		"cmd":     req["cmd"],
		"version": session.Table.version,
		"data":    getAll(session.Table.tmpFile),
	})
	session.Send(req["cmd"].(string), b)
}

func handleRollback(req map[string]interface{}, session *Session) {
	ver := int(req["version"].(float64))
	table := session.Table

	if ver > table.version {
		pushErr(req["cmd"], "版本号错误", session)
		return
	}

	/*
		if len(table.sessionMap) > 1 {
			pushErr(req["cmd"], "多人操作，不能回退", session)
			return
		}
	*/

	newF := newFile(getAll(table.xlFile))
	if ver < table.version {
		// 执行指令
		for id := table.version; id > ver; id-- {
			//_, _, ret, err := pgsql.LoadCmd(table.tableName, i)
			//if err != nil {
			//	pushErr(req["cmd"], err.Error(), session)
			//	return
			//}
			if v, ok := table.verHistory[id]; ok {
				rollbackCmds(newF, v.cmds)
			}
		}

		// 保存版本回退指令
		cmdStr := MustJsonMarshal([]map[string]interface{}{
			{
				"cmd":       req["cmd"],
				"tableName": table.tableName,
				"now":       table.version,
				"goto":      ver,
			},
		})
		users := MustJsonMarshal([]string{session.UserName})
		v, err := pgsql.InsertCmd(table.tableName, string(users), string(cmdStr))
		if err != nil {
			pushErr(req["cmd"], err.Error(), session)
			return
		}

		// 更新数据库
		b := MustJsonMarshal(getAll(newF))
		err = pgsql.UpdateTableData(table.tableName, table.version, string(b))
		if err != nil {
			pushErr(req["cmd"], err.Error(), session)
			// todo 应该删除回退操作日志
			return
		}
		table.version = v
	}

	// 当前表状态更改
	table.loadHistory()
	table.cmds = []map[string]interface{}{}
	table.cmdUsers = []string{}
	table.xlFile = newFile(getAll(newF))
	table.tmpFile = newF

	// 同步给所有人
	table.pushAllSession(map[string]interface{}{
		"cmd":     "rollback",
		"version": table.version,
		"data":    getAll(newF),
	})
}

//## pushError 返回错误信息
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
	// 用户选择的单元格，只做转发
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
	// 获取版本列表
	dispatcher["versionList"] = handleVersionList
}
