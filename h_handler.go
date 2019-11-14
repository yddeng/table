package table

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/yddeng/table/pgsql"
	"net/http"
)

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
