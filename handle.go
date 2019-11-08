package table

import (
	"fmt"
	"net/http"
)

func cellValue(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	//跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	var axis string
	var value string
	if a, ok := r.Form["axis"]; ok {
		axis = a[0]
	}

	if v, ok := r.Form["value"]; ok {
		value = v[0]
	}

	ef := GetExcel("test.xlsx")
	err := ef.WriteCell(nil, axis, value)
	fmt.Println(ef, axis, value, err)
	c, _ := ef.GetCell(axis)
	fmt.Println("cellValue", c)
	ef.Save()
}
