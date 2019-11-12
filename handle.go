package table

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sniperHW/kendynet"
)

type handFunc func(req map[string]interface{}, session *Session)

var dispatcher = map[string]handFunc{}

func DispatchMessage(msg map[string]interface{}, session kendynet.StreamSession) {
	/*err := checkBase(msg)
	if err != nil {
		b, err := json.Marshal(map[string]interface{}{
			"code": 0,
			"msg":  err.Error(),
		})
		if err == nil {
			_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
		}
		fmt.Println(err)
		return
	}*/

	cmd := msg["cmd"].(string)
	if cmd == "openFile" {
		onOpenFile(msg, session)
	} else {
		sess := checkSession(session)
		if sess != nil {
			handler, ok := dispatcher[cmd]
			if ok {
				handler(msg, sess)
			} else {
				fmt.Println("no handler", cmd)
			}
		} else {
			fmt.Println("no session", session.RemoteAddr().String())
		}
	}
}

func checkBase(req map[string]interface{}) error {
	_, ok := req["cmd"]
	if !ok {
		return fmt.Errorf("缺少参数cmd")
	}

	_, ok = req["fileName"]
	if !ok {
		return fmt.Errorf("缺少参数fileName")
	}
	return nil
}

func checkSession(sess kendynet.StreamSession) *Session {
	return sessionMap[sess.RemoteAddr().String()]
}

func handleSetCellValues(req map[string]interface{}, session *Session) {
	fmt.Println("handleSetCellValues", req)

	cellValues := req["cellValues"].([]interface{})
	file := session.File

	fields := map[string]interface{}{}
	for _, v := range cellValues {
		item := v.(map[string]interface{})
		cellName, err := excelize.CoordinatesToCellName(int(item["col"].(float64))+1, int(item["row"].(float64))+1)
		if err == nil {
			fields[cellName] = item["newValue"]
		}
	}
	file.SetCellValues(fields)
	_ = file.Save()

	file.PushData()

}

func handleCellSelected(req map[string]interface{}, session *Session) {
	fmt.Println("handleCellSelected", req)

	selected := req["selected"].([]interface{})
	file := session.File
	fmt.Println("handleCellSelected", file)

	// 先清空当前session的锁定
	for axis, sess := range file.cellLocked {
		if sess == session {
			delete(file.cellLocked, axis)
		}
	}
	// 锁定当前选中
	for _, v := range selected {
		item := v.(map[string]interface{})
		cellName, err := excelize.CoordinatesToCellName(int(item["col"].(float64))+1, int(item["row"].(float64))+1)
		if err == nil {
			file.cellLocked[cellName] = session
		}
	}
	file.PushData()
}

func handleInsertRow(req map[string]interface{}, session *Session) {
	rowIndex := req["rowIndex"].(int)
	file := session.File

	file.xlFile.InsertRow(file.activeSheet, rowIndex)
	session.addEvent(req)
}

func handleRemveRow(req map[string]interface{}, session *Session) {
	rowIndex := req["rowIndex"].(int)
	file := session.File

	file.xlFile.RemoveRow(file.activeSheet, rowIndex)
	session.addEvent(req)
}

func handleInsertCol(req map[string]interface{}, session *Session) {
	colIndex := req["colIndex"].(int)
	file := session.File

	//celHeader := excelize.
	file.xlFile.InsertCol(file.activeSheet, colIndex)
	session.addEvent(req)
}

func handleRemoveCol(req map[string]interface{}, session *Session) {
	rowIndex := req["rowIndex"].(int)
	file := session.File

	file.xlFile.RemoveCol(file.activeSheet, rowIndex)
	session.addEvent(req)
}

func init() {
	dispatcher["setCellValues"] = handleSetCellValues
	dispatcher["cellSelected"] = handleCellSelected
	dispatcher["insertRow"] = handleInsertRow
	dispatcher["removeRow"] = handleRemveRow
	dispatcher["insertCol"] = handleInsertCol
	dispatcher["removeCol"] = handleRemoveCol
}
