package pgsql

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	err := Set("user", map[string]interface{}{
		"user_name": "ddd",
		"password":  "erer",
	})
	fmt.Println(err)
}

func TestUpdate(t *testing.T) {
	err := Update("user", "user_name = '123'", map[string]interface{}{
		"password": "erer",
	})
	fmt.Println(err)
}

func TestSetNx(t *testing.T) {
	err := SetNx("user", "user_name", map[string]interface{}{
		"user_name":  "123",
		"password":   321,
		"permission": "",
	})
	fmt.Println(err)
}

func TestGet(t *testing.T) {
	ret, err := Get("table_data", "table_name = 'dadada'", []string{"table_name", "describe", "version"})
	fmt.Println(ret["table_name"].(string), ret["describe"].(string), ret["version"].(int64), err)
}

func TestGetAll(t *testing.T) {
	ret, err := GetAll("user", []string{"password", "user_name", "permission"})
	fmt.Println(ret, err)
}
