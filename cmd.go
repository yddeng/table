package table

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

var Sheet = "Sheet1"

func doCmds(file *excelize.File, req []map[string]interface{}) {
	for _, v := range req {
		doCmd(file, v, false)
	}
}

func rollbackCmds(file *excelize.File, req []map[string]interface{}) {
	for i := len(req) - 1; i >= 0; i-- {
		doCmd(file, req[i], true)
	}
}

func doCmd(file *excelize.File, req map[string]interface{}, rollbcak bool) {
	cmd := req["cmd"].(string)
	//fmt.Println("doCmd", cmd, rollbcak)
	switch cmd {
	case "insertRow":
		insertRow(file, req, rollbcak)
	case "removeRow":
		removeRow(file, req, rollbcak)
	case "insertCol":
		insertCol(file, req, rollbcak)
	case "removeCol":
		removeCol(file, req, rollbcak)
	case "setCellValues":
		setCellValues(file, req, rollbcak)
	default:
		panic("cmd defined")
	}
}

func insertRow(file *excelize.File, req map[string]interface{}, rollbcak bool) {
	index := int(req["index"].(float64)) + 1
	var err error
	if rollbcak {
		err = file.RemoveRow(Sheet, index)
	} else {
		err = file.InsertRow(Sheet, index)
	}
	CheckErr(err)
}

func removeRow(file *excelize.File, req map[string]interface{}, rollbcak bool) {
	index := int(req["index"].(float64)) + 1
	var err error
	if rollbcak {
		err = file.InsertRow(Sheet, index)
	} else {
		err = file.RemoveRow(Sheet, index)
	}
	CheckErr(err)
}

func insertCol(file *excelize.File, req map[string]interface{}, rollbcak bool) {
	index := int(req["index"].(float64)) + 1
	celHeader, err := excelize.ColumnNumberToName(index)
	CheckErr(err)
	if rollbcak {
		err = file.RemoveCol(Sheet, celHeader)
	} else {
		err = file.InsertCol(Sheet, celHeader)
	}
	CheckErr(err)
}

func removeCol(file *excelize.File, req map[string]interface{}, rollbcak bool) {
	index := int(req["index"].(float64)) + 1
	celHeader, err := excelize.ColumnNumberToName(index)
	CheckErr(err)
	if rollbcak {
		err = file.InsertCol(Sheet, celHeader)
	} else {
		err = file.RemoveCol(Sheet, celHeader)
	}
	CheckErr(err)
}

func setCellValues(file *excelize.File, req map[string]interface{}, rollbcak bool) {
	cellValues := req["cellValues"].([]interface{})
	for _, v := range cellValues {
		item := v.(map[string]interface{})
		cellName, err := excelize.CoordinatesToCellName(int(item["col"].(float64))+1, int(item["row"].(float64))+1)
		CheckErr(err)
		if rollbcak {
			err = file.SetCellValue(Sheet, cellName, item["oldValue"])
		} else {
			err = file.SetCellValue(Sheet, cellName, item["newValue"])
		}
		CheckErr(err)
	}
}

func getAll(file *excelize.File) [][]string {
	data, err := file.GetRows(Sheet)
	CheckErr(err)
	return data
}
