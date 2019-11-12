package table

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sniperHW/kendynet/event"
	"github.com/sniperHW/kendynet/message"
	"sync"
	"time"
)

type ExcelFile struct {
	fileName    string
	xlFile      *excelize.File
	activeSheet string
	sessionMap  map[string]*Session // session, for remoteAddr
	cellLocked  map[string]*Session // session, for axis
	version     int
	mu          sync.Mutex
}

var (
	once       = sync.Once{}
	eventQueue = event.NewEventQueue()
	excelFiles = map[string]*ExcelFile{}
)

func PostTask(task func()) {
	_ = eventQueue.Post(task)
}

func CreateExcel(fileName string) error {
	file := excelize.NewFile()
	index := file.NewSheet("Sheet1")
	file.SetActiveSheet(index)
	return file.SaveAs(fmt.Sprintf("%s.xlsx", fileName))
}

func OpenExcel(fileName string) (*ExcelFile, error) {
	xlFile, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	return &ExcelFile{
		fileName:    fileName,
		xlFile:      xlFile,
		activeSheet: xlFile.GetSheetName(xlFile.GetActiveSheetIndex()),
		sessionMap:  map[string]*Session{},
		cellLocked:  map[string]*Session{},
		mu:          sync.Mutex{},
	}, nil
}

func (ef *ExcelFile) Save() error {
	return ef.xlFile.Save()
}

func (ef *ExcelFile) GetAll() ([][]string, error) {
	return ef.xlFile.GetRows(ef.activeSheet)
}

func (ef *ExcelFile) GetRow(idx int) ([]string, error) {
	rows, err := ef.xlFile.GetRows(ef.activeSheet)
	if err != nil {
		return nil, err
	}
	return rows[idx], nil
}

func (ef *ExcelFile) GetCellValue(axis string) (string, error) {
	return ef.xlFile.GetCellValue(ef.activeSheet, axis)
}

func (ef *ExcelFile) SetCellValues(fields map[string]interface{}) {
	for k, v := range fields {
		fmt.Println(k, v)
		err := ef.xlFile.SetCellValue(ef.activeSheet, k, v)
		if err != nil {
			fmt.Println("SetValues", err)
		}

	}
}

func (ef *ExcelFile) AddSession(session *Session) {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	if _, ok := ef.sessionMap[session.RemoteAddr().String()]; !ok {
		ef.sessionMap[session.RemoteAddr().String()] = session
	}
}

func (ef *ExcelFile) RemoveSession(session *Session) {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	if _, ok := ef.sessionMap[session.RemoteAddr().String()]; ok {
		delete(ef.sessionMap, session.RemoteAddr().String())
	}
	// 清理锁定
	for axis, sess := range ef.cellLocked {
		if sess == session {
			delete(ef.cellLocked, axis)
		}
	}
}

func (ef *ExcelFile) PushData() {
	data, _ := ef.GetAll()
	resp := map[string]interface{}{
		"cmd":      "pushData",
		"fileName": ef.fileName,
		"data":     data,
	}
	if len(ef.cellLocked) > 0 {
		locked := []map[string]interface{}{}
		for axis, sess := range ef.cellLocked {
			col, row, err := excelize.CellNameToCoordinates(axis)
			if err == nil {
				locked = append(locked, map[string]interface{}{
					"col":      col - 1,
					"row":      row - 1,
					"userName": sess.UserName,
				})
			}
		}
		resp["cellLocked"] = locked
	}

	b, _ := json.Marshal(resp)
	for _, sess := range ef.sessionMap {
		_ = sess.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
	}
	fmt.Println("pushData")
}

func Loop() {
	once.Do(func() {
		go func() {
			eventQueue.Run()
		}()
	})

	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		PostTask(func() {
			for _, file := range excelFiles {
				if len(file.sessionMap) == 0 {
					_ = file.Save()
					delete(excelFiles, file.fileName)
				}
			}
		})
	}
}
