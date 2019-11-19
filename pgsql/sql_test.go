package pgsql

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"testing"
)

func TestCreateTable(t *testing.T) {
	//CreateTableCmd("deng")

	InsertCmd("deng", "yidong", "testets")
}

func TestLoadTableData(t *testing.T) {
	_, _, _, _, err := LoadTableData("df")
	fmt.Println(err)
}

func TestInsert(t *testing.T) {
	Insert("table_data", map[string]interface{}{
		"table_name":  "ddd",
		"now_version": 3,
		"data":        "dfs",
	})
}

func TestUpdateTableData(t *testing.T) {
	file := excelize.NewFile()
	file.SetCellValue("Sheet1", "A1", "")
	data, err := file.GetRows("Sheet1")
	fmt.Println(data, err)
	b, err := json.Marshal(data)
	fmt.Println(string(b), err)
	err = UpdateTableData("ydd", map[string]interface{}{
		"version": 2,
		"date":    "sdd",
		"data":    string(b),
	})
	fmt.Println(err)
	v, _, _, d, err := LoadTableData("ddd")
	fmt.Println(v, d, err)
}

func TestNewSelect(t *testing.T) {
	c, err := NewSelect("table_data", "version,date,describe", "table_name")
	fmt.Println(err)
	if err == nil {
		c.Query(func(i map[string]interface{}, e error) {
			fmt.Println(i, e)
		}, "test")
	}
}
