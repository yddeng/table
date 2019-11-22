package table

import (
	"encoding/json"
	"fmt"
	"github.com/yddeng/table/pgsql"
	"net/http"
	"net/url"
	"strings"
)

// 创建文件
func HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleCreateTable request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	req := map[string]string{
		"tableName": "",
		"describe":  "",
		"token":     "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	tableName := req["tableName"]
	desc := req["describe"]
	// 判断名字是否存在
	_, _, _, _, err = pgsql.LoadTableData(tableName)
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
	err = pgsql.InsertTableData(tableName, desc)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok": 1,
	}); err != nil {
		logger.Errorln("http resp err:", err)
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

	req := map[string]string{
		"tableName": "",
		"token":     "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	// todo

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok": 1,
	}); err != nil {
		logger.Errorln("http resp err:", err)
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
		logger.Errorln("http resp err:", err)
	}
}

// 下载excel
func HandleDownloadTable(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleDownloadTable request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	req := map[string]string{
		"tableName": "",
		"token":     "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	tableName := req["tableName"]
	_, _, _, data, err := pgsql.LoadTableData(tableName)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   1,
		"data": data,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

// login登陆
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleLogin request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	req := map[string]string{
		"userName": "",
		"password": "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	userName := req["userName"]
	password := req["password"]
	pwd, _, err := pgsql.LoadUser(userName)
	if err != nil {
		httpErr("该用户不存在", w)
		return
	}

	if strings.Compare(pwd, password) != 0 {
		httpErr("密码错误", w)
		return
	}

	token := makeToken(userName)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":    1,
		"token": token,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

// 添加用户
func HandleAddUser(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleAddUser request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	req := map[string]string{
		"userName": "",
		"password": "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	userName := req["userName"]
	password := req["password"]
	_, _, err = pgsql.LoadUser(userName)
	if err == nil {
		httpErr("用户名已存在", w)
		return
	}

	err = pgsql.InsertUser(userName, password, "")
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	token := makeToken(userName)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":    1,
		"token": token,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

func checkForm(form url.Values, args map[string]string) error {
	for k := range args {
		if t, ok := form[k]; ok {
			args[k] = t[0]
		} else {
			return fmt.Errorf("key:%s not found\n", k)
		}
	}
	return nil
}

func httpErr(err string, w http.ResponseWriter) {
	logger.Errorln("httpErr", err)
	resp := map[string]interface{}{
		"ok":  0,
		"msg": err,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Errorln("http resp err:", err)
	}
}
