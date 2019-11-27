package table

import (
	"encoding/json"
	"fmt"
	"github.com/yddeng/table/pgsql"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	_, err = pgsql.Get("table_data", fmt.Sprintf("table_name = '%s'", tableName), []string{"table_name"})
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

	// 创建tag表列
	err = pgsql.AlterAddTagData(tableName)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	// todo 失败，回滚
	err = pgsql.Set("table_data", map[string]interface{}{
		"table_name": tableName,
		"describe":   desc,
		"version":    0,
		"date":       GenDateTimeString(time.Now()),
		"data":       string(MustJsonMarshal([]string{})),
	})
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

// 更新表描述
func HandleUpdateDescribe(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleUpdateDescribe request", r.Method, r.Form)

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

	// todo
	err = pgsql.Update("table_data", fmt.Sprintf("table_name = '%s'", req["tableName"]), map[string]interface{}{
		"describe": req["describe"],
	})
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

	ret, err := pgsql.GetAll("table_data", []string{"table_name", "describe", "version", "date"})
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
	ret, err := pgsql.Get("table_data", fmt.Sprintf("table_name = '%s'", tableName), []string{"data"})
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	var data [][]string
	MustJsonUnmarshal(([]byte)(ret["data"].(string)), &data)

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
	ret, err := pgsql.Get("user", fmt.Sprintf("user_name = '%s'", userName), []string{"password"})
	if err != nil {
		httpErr("该用户不存在", w)
		return
	}

	pwd := ret["password"].(string)
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
	_, err = pgsql.Get("user", fmt.Sprintf("user_name = '%s'", userName), []string{"password"})
	if err == nil {
		httpErr("用户名已存在", w)
		return
	}

	err = pgsql.Set("user", map[string]interface{}{
		"user_name":  userName,
		"password":   password,
		"permission": "",
	})
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

func HandleAddTag(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleAddUser request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	req := map[string]string{
		"tagName": "",
		"token":   "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	sqlReq := map[string]interface{}{
		"tag_name": req["tagName"],
		"date":     GenDateTimeString(time.Now()),
	}

	desc := []string{}
	tables := r.Form["tables[]"]
	for _, name := range tables {
		ret, err := pgsql.Get("table_data", fmt.Sprintf("table_name = '%s'", name), []string{"version", "data"})
		if err != nil {
			httpErr(err.Error(), w)
			return
		}
		sqlReq[name] = ret["data"]
		desc = append(desc, fmt.Sprintf("%s:%d", name, ret["version"].(int64)))
	}

	sqlReq["describe"] = strings.Join(desc, ",")
	err = pgsql.Set("tag_data", sqlReq)
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

func HandleShowTag(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleShowTag request", r.Method, r.Form)

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	req := map[string]string{
		"token": "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	ret, err := pgsql.GetAll("tag_data", []string{"tag_name", "describe", "date"})
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   1,
		"tags": ret,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

func HandleDownTag(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	logger.Infoln("HandleDownTag request", r.Method, r.Form)
	httpHeader(&w) // 跨域

	req := map[string]string{
		"tagName": "",
		"token":   "",
	}
	err := checkForm(r.Form, req)
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	ret, err := pgsql.GetTag(req["tagName"])
	if err != nil {
		httpErr(err.Error(), w)
		return
	}

	tables := map[string]interface{}{}
	for k, v := range ret {
		if k != "tag_name" && k != "date" && k != "describe" {
			if v != nil && v.(string) != "" {
				var data [][]string
				MustJsonUnmarshal(([]byte)(v.(string)), &data)
				tables[k] = data
			}
		}
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":     1,
		"tables": tables,
	}); err != nil {
		logger.Errorf("http resp err:", err)
	}
}

func checkForm(form url.Values, args map[string]string) error {
	for k := range args {
		v := form.Get(k)
		if v != "" {
			args[k] = v
		} else {
			return fmt.Errorf("key:%s not found\n", k)
		}
	}
	return nil
}

func httpHeader(w *http.ResponseWriter) {
	//跨域
	(*w).Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	(*w).Header().Set("content-type", "application/json")             //返回数据格式是json
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
