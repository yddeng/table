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
	_, _, err := LoadTableData("df")
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
	err = UpdateTableData("55", 1, string(b))

	v, d, err := LoadTableData("ddd")
	fmt.Println(v, d, err)
}
