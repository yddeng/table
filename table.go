package table

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sniperHW/kendynet/event"
	"github.com/yddeng/table/pgsql"
	"sync"
	"time"
)

var (
	once       = sync.Once{}
	eventQueue = event.NewEventQueue()
	tables     = map[string]*Table{}
)

type Table struct {
	tableName  string
	version    int                 // 当前版本号
	tmpFile    *excelize.File      // 中间文件，同步给用户
	xlFile     *excelize.File      // 指令计算的结果文件，最终落地到数据库
	sessionMap map[string]*Session // session, for remoteAddr

	cmds     []map[string]interface{} // 当前cmd指令的集合
	cmdUsers []string                 // 产生当前指令的用户
}

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
	xlFile := newFile(data)
	tmpFile := newFile(data)
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
		// todo 验证正确性
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
		this.cmdUsers = []string{}

		// 保存最新数据
		b := MustJsonMarshal(getAll(this.xlFile))
		err = pgsql.UpdateTableData(this.tableName, this.version, string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func newFile(data [][]string) *excelize.File {
	file := excelize.NewFile()
	idx := file.NewSheet(Sheet)
	file.SetActiveSheet(idx)

	for row, v := range data {
		for col, value := range v {
			cellName, err := excelize.CoordinatesToCellName(col+1, row+1)
			CheckErr(err)
			_ = file.SetCellValue(Sheet, cellName, value)
		}
	}
	return file
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

func (this *Table) AddSession(session *Session) {
	if _, ok := this.sessionMap[session.RemoteAddr()]; !ok {
		this.sessionMap[session.RemoteAddr()] = session
	}
}

func (this *Table) RemoveSession(session *Session) {
	if _, ok := this.sessionMap[session.RemoteAddr()]; ok {
		delete(this.sessionMap, session.RemoteAddr())
	}
}

// 转发命令
func (this *Table) pushAllSession(msg map[string]interface{}) {
	b := MustJsonMarshal(msg)
	for _, session := range this.sessionMap {
		session.Send(msg["cmd"].(string), b)
	}
}

func (this *Table) pushOtherUser(msg map[string]interface{}, userName string) {
	b := MustJsonMarshal(msg)
	for _, session := range this.sessionMap {
		if session.UserName != userName {
			session.Send(msg["cmd"].(string), b)
		}
	}
}

// pushAll 推送所有数据
func (this *Table) pushAll() {
	data := getAll(this.tmpFile)
	resp := map[string]interface{}{
		"cmd":     "pushAll",
		"version": this.version,
		"data":    data,
	}
	this.pushAllSession(resp)
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
					_ = file.SaveTable()
					fmt.Println("delete Table", file.tableName)
					delete(tables, file.tableName)
				}
			}
		})
	}
}
