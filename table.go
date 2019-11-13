package table

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sniperHW/kendynet/event"
	"github.com/sniperHW/kendynet/message"
	"github.com/yddeng/table/pgsql"
	"sync"
	"time"
)

type Table struct {
	fileName     string
	version      int                 // 当前版本号
	tmpFile      *excelize.File      // 中间文件，同步给用户
	xlFile       *excelize.File      // 指令计算的结果文件，最终落地到数据库
	sessionMap   map[string]*Session // session, for remoteAddr
	cellSelected map[string]*Session // session, for axis
	mu           sync.Mutex
}

var (
	once       = sync.Once{}
	eventQueue = event.NewEventQueue()
	tables     = map[string]*Table{}
)

func PostTask(task func()) {
	_ = eventQueue.Post(task)
}

// 从数据库table_data 加载表数据
// 将数据赋值给 xlFile, tmpFile
func OpenTable(fileName string) (*Table, error) {
	fmt.Println("loadTable", fileName)
	v, data, err := pgsql.LoadTableData(fileName)
	if err != nil {
		return nil, err
	}
	fmt.Println(v, data, err)
	xlFile := newFile()
	tmpFile := newFile()
	cloneFile(xlFile, data)
	cloneFile(tmpFile, data)
	return &Table{
		fileName:     fileName,
		version:      v,
		tmpFile:      tmpFile,
		xlFile:       xlFile,
		sessionMap:   map[string]*Session{},
		cellSelected: map[string]*Session{},
		mu:           sync.Mutex{},
	}, nil
}

// 保存，这里指落地到数据库
func (this *Table) SaveTable() error {
	data, err := getAll(this.xlFile)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(data)
	return pgsql.UpdateTableData(this.fileName, this.version, string(b))
}

func newFile() *excelize.File {
	file := excelize.NewFile()
	idx := file.NewSheet(Sheet)
	file.SetActiveSheet(idx)
	return file
}

func fileSave(fileName string, file *excelize.File) {
	_ = file.SaveAs(fmt.Sprintf("%s.xlsx", fileName))
}

func cloneFile(file *excelize.File, data [][]string) {
	for row, v := range data {
		for col, value := range v {
			cellName, err := excelize.CoordinatesToCellName(col+1, row+1)
			checkErr(err)
			if err == nil {
				_ = file.SetCellValue("Sheet1", cellName, value)
			}
		}
	}
}

func (ef *Table) GetAll() ([][]string, error) {
	return ef.xlFile.GetRows(Sheet)
}

func (ef *Table) GetRow(idx int) ([]string, error) {
	rows, err := ef.xlFile.GetRows(Sheet)
	if err != nil {
		return nil, err
	}
	return rows[idx], nil
}

func (ef *Table) GetCellValue(axis string) (string, error) {
	return ef.xlFile.GetCellValue(Sheet, axis)
}

func (ef *Table) SetCellValues(fields map[string]interface{}) {
	for k, v := range fields {
		fmt.Println(k, v)
		err := ef.xlFile.SetCellValue(Sheet, k, v)
		if err != nil {
			fmt.Println("SetValues", err)
		}

	}
}

func (ef *Table) AddSession(session *Session) {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	if _, ok := ef.sessionMap[session.RemoteAddr().String()]; !ok {
		ef.sessionMap[session.RemoteAddr().String()] = session
	}
}

func (ef *Table) RemoveSession(session *Session) {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	if _, ok := ef.sessionMap[session.RemoteAddr().String()]; ok {
		delete(ef.sessionMap, session.RemoteAddr().String())
	}
	// 清理锁定
	for axis, sess := range ef.cellSelected {
		if sess == session {
			delete(ef.cellSelected, axis)
		}
	}
}

func (ef *Table) getCellSelected() []map[string]interface{} {
	locked := []map[string]interface{}{}
	for axis, sess := range ef.cellSelected {
		col, row, err := excelize.CellNameToCoordinates(axis)
		if err == nil {
			locked = append(locked, map[string]interface{}{
				"col":      col - 1,
				"row":      row - 1,
				"userName": sess.UserName,
			})
		}
	}
	return locked
}

//## pushData 推送最新的表数据, toc
//```
//cmd : "pushData"
//data:  data  //[][]string
//```

func (ef *Table) PushData() {
	data, _ := getAll(ef.tmpFile)
	resp := map[string]interface{}{
		"cmd":     "pushData",
		"version": ef.version,
		"data":    data,
	}
	b, _ := json.Marshal(resp)
	for _, sess := range ef.sessionMap {
		_ = sess.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
	}
}

//## pushCellData 推送变化的单元格数据
//```
//cmd : "pushCellData"
//cellDate:  // []map[string]string
//              []{col:int,row:int,oldValue:string,newValue:string,userName:string}
//```

func (ef *Table) pushCellData() {

}

//## pushCellSelected
//```
//cmd : "pushCellSelected"
//selected:  // []map[string]string
//              []{col:int,row:int,userName:string}
//``

func (ef *Table) pushCellSelected() {
	if len(ef.cellSelected) > 0 {
		resp := map[string]interface{}{
			"cmd":      "pushCellSelected",
			"selected": ef.getCellSelected(),
		}
		b, _ := json.Marshal(resp)
		for _, sess := range ef.sessionMap {
			_ = sess.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
		}
		fmt.Println("pushCellSelected")
	}
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
			for _, file := range tables {
				if len(file.sessionMap) == 0 {
					file.xlFile.Save()
					fmt.Println("delete Table", file.fileName)
					delete(tables, file.fileName)
				}
			}
		})
	}
}
