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
	tableName  string
	version    int                      // 当前版本号
	tmpFile    *excelize.File           // 中间文件，同步给用户
	xlFile     *excelize.File           // 指令计算的结果文件，最终落地到数据库
	sessionMap map[string]*Session      // session, for remoteAddr
	cmds       []map[string]interface{} // 指令集合
	cmdUsers   []string
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
func OpenTable(tableName string) (*Table, error) {
	fmt.Println("loadTable", tableName)
	v, data, err := pgsql.LoadTableData(tableName)
	if err != nil {
		return nil, err
	}
	//fmt.Println(v, data, err)
	xlFile := newFile()
	tmpFile := newFile()
	cloneFile(xlFile, data)
	cloneFile(tmpFile, data)
	return &Table{
		tableName:  tableName,
		version:    v,
		tmpFile:    tmpFile,
		xlFile:     xlFile,
		sessionMap: map[string]*Session{},
		cmds:       []map[string]interface{}{},
		cmdUsers:   []string{},
	}, nil
}

// 保存，这里指落地到数据库
func (this *Table) SaveTable() error {
	if len(this.cmds) > 0 {
		doCmds(this.xlFile, this.cmds)

		// 保存版本指令
		cmdsStr, _ := json.Marshal(this.cmds)
		users, _ := json.Marshal(this.cmdUsers)
		v, err := pgsql.InsertCmd(this.tableName, string(users), string(cmdsStr))
		if err != nil {
			return err
		}
		this.version = v
		this.cmds = []map[string]interface{}{}

		// 保存最新数据
		data, _ := getAll(this.xlFile)
		b, _ := json.Marshal(data)
		err = pgsql.UpdateTableData(this.tableName, this.version, string(b))
		if err != nil {
			return err
		}

		this.pushSaveTable()
	}
	return nil
}

func newFile() *excelize.File {
	file := excelize.NewFile()
	idx := file.NewSheet(Sheet)
	file.SetActiveSheet(idx)
	return file
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

func (this *Table) AddCmd(cmd map[string]interface{}, userName string) {
	this.cmds = append(this.cmds, cmd)
	for _, name := range this.cmdUsers {
		if name == userName {
			return
		}
	}
	this.cmdUsers = append(this.cmdUsers, userName)
}

func (ef *Table) AddSession(session *Session) {
	if _, ok := ef.sessionMap[session.RemoteAddr().String()]; !ok {
		ef.sessionMap[session.RemoteAddr().String()] = session
	}
}

func (ef *Table) RemoveSession(session *Session) {
	if _, ok := ef.sessionMap[session.RemoteAddr().String()]; ok {
		delete(ef.sessionMap, session.RemoteAddr().String())
	}
}

// pushSaveTable 保存表后推送
func (this *Table) pushSaveTable() {
	data, _ := getAll(this.xlFile)
	b, _ := makeSaveTable(this.version, data)
	for _, session := range this.sessionMap {
		_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
	}
}

// pushAll 推送所有数据
func (this *Table) pushAll() {
	data, _ := getAll(this.tmpFile)
	b, _ := makePushAll(this.tableName, this.version, data)
	for _, session := range this.sessionMap {
		_ = session.SendMessage(message.NewWSMessage(message.WSTextMessage, b))
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
					file.SaveTable()
					fmt.Println("delete Table", file.tableName)
					delete(tables, file.tableName)
				}
			}
		})
	}
}
