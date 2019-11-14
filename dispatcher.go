package table

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sniperHW/kendynet"
	"github.com/sniperHW/kendynet/message"
	"github.com/yddeng/table/pgsql"
	"net/http"
)

func Dispatcher(msg map[string]interface{}, session kendynet.StreamSession) {
	cmd := msg["cmd"].(string)
	fmt.Println("Dispatcher", msg)

	switch cmd {
	case "openTable":
		onOpenTable(msg, session)
	case "saveTable":
		sess := checkSession(session)
		if sess != nil {
			handleSaveTable(msg, sess)
		}
	case "rollback":
		sess := checkSession(session)
		if sess != nil {
			handleRollback(msg, sess)
		}
	case "cellSelected":
		sess := checkSession(session)
		if sess != nil {
			handleCellSelected(msg, sess)
		}
	default:
		sess := checkSession(session)
		if sess != nil {
			doCmd(sess.Table.tmpFile, msg, false)
			sess.addCmd(msg)
			sess.Table.PushData()
		} else {
			fmt.Println("no session", session.RemoteAddr().String())
		}
	}
}

func checkSession(sess kendynet.StreamSession) *Session {
	return sessionMap[sess.RemoteAddr().String()]
}

func handleCellSelected(req map[string]interface{}, session *Session) {
	selected := req["selected"].([]interface{})
	table := session.Table

	// 先清空当前session的锁定
	for axis, sess := range table.cellSelected {
		if sess == session {
			delete(table.cellSelected, axis)
		}
	}
	// 锁定当前选中
	for _, v := range selected {
		item := v.(map[string]interface{})
		cellName, err := excelize.CoordinatesToCellName(int(item["col"].(float64))+1, int(item["row"].(float64))+1)
		if err == nil {
			table.cellSelected[cellName] = session
		}
	}

	table.pushCellSelected()
}

func handleSaveTable(req map[string]interface{}, session *Session) {
	doCmds(session.Table.xlFile, session.doCmds)
	session.SaveCmd()
	session.Table.PushData()
}

func handleRollback(req map[string]interface{}, session *Session) {
	now := int(req["now"].(float64))
	ver := int(req["goto"].(float64))
	table := session.Table

	if now != table.version {
		rollbackErr("版本号不一致，不能回退", session)
		table.PushData()
		return
	}

	if len(table.sessionMap) > 1 {
		rollbackErr("多人操作，不能回退", session)
		return
	}

	if len(session.doCmds) > 0 {
		rollbackErr("当前有操作没有保存，不能回退", session)
		return
	}

	if ver < table.version {
		for i := table.version; i > ver; i-- {
			ret, err := pgsql.LoadCmd(table.fileName, i)
			if err != nil {
				rollbackErr(err.Error(), session)
				return
			}
			rollbackCmds(table.tmpFile, ret)
		}
	} else if ver > table.version {
		_, err := pgsql.LoadCmd(table.fileName, ver)
		if err != nil {
			rollbackErr("版本号错误", session)
			return
		}
		for i := table.version; i <= ver; i++ {
			ret, err := pgsql.LoadCmd(table.fileName, i)
			if err != nil {
				rollbackErr(err.Error(), session)
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
		"ok":      1,
		"version": ver,
		"data":    data,
	})
	_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))

}

func rollbackErr(err string, session *Session) {
	fmt.Println("rollbackErr", err)
	resp := map[string]interface{}{
		"cmd": "rollback",
		"ok":  0,
		"msg": err,
	}
	b, _ := json.Marshal(resp)
	_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
}

/************************************ http ***********************************************************/
// 创建文件
func HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleCreateTable request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var tableName, userName string
	if t, ok := r.Form["tableName"]; ok {
		tableName = t[0]
	}

	if p, ok := r.Form["userName"]; ok {
		userName = p[0]
	}

	if tableName == "" || userName == "" {
		httpErr("参数错误", w)
		return
	}

	// 判断名字是否存在
	_, _, err := pgsql.LoadTableData(tableName)
	if err == nil {
		httpErr("名字重复", w)
		return
	}

	// 创建指令表
	err = pgsql.CreateTableCmd(tableName)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	// todo 失败，回滚
	// 添加数据
	b, _ := json.Marshal([]string{})
	err = pgsql.InsertTableData(tableName, 0, string(b))
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok": 1,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

// 删除文件
func HandleDeleteTable(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleDeleteTable request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var tableName, userName string
	if t, ok := r.Form["tableName"]; ok {
		tableName = t[0]
	}

	if p, ok := r.Form["userName"]; ok {
		userName = p[0]
	}

	if tableName == "" || userName == "" {
		httpErr("参数错误", w)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok": 1,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

// 获取文件列表
func HandleGetAllTable(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleGetAllTable request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	ret, err := pgsql.AllTableData()
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":     1,
		"tables": ret,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

// 下载excel
func HandleDownloadTable(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleGetAllTable request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var tableName string
	if t, ok := r.Form["tableName"]; ok {
		tableName = t[0]
	}

	if tableName == "" {
		httpErr("参数错误", w)
		return
	}

	_, data, err := pgsql.LoadTableData(tableName)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	xlFile := newFile()
	cloneFile(xlFile, data)
	fileName := fmt.Sprintf("%s.xlsx", tableName)
	err = xlFile.SaveAs(fileName)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	file, err := excelize.OpenFile(fileName)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	err = file.Write(w)
	if err != nil {
		logger.Errorf("http resp err:", err)
	}
}

func httpErr(err string, w http.ResponseWriter) {
	fmt.Println("httpErr", err)
	resp := map[string]interface{}{
		"ok":  0,
		"msg": err,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorf("http resp err:", err)
	}
}
