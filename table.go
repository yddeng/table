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
	nameMap    map[string]int      // 统计当前有哪些用户,且链接数量

	cmds       []map[string]interface{} // 当前cmd指令的集合
	cmdUsers   []string                 // 产生当前指令的用户
	verHistory map[int]*VerHistory      // 历史版本，查询20条
}

type VerHistory struct {
	version int
	users   []string
	date    string
	cmds    []map[string]interface{}
}

func PostTask(task func()) {
	_ = eventQueue.Post(task)
}

// 从数据库table_data 加载表数据
// 将数据赋值给 xlFile, tmpFile
func OpenTable(tableName string) (*Table, error) {
	//fmt.Println("loadTable", tableName)
	ret, err := pgsql.Get("table_data", fmt.Sprintf("table_name = '%s'", tableName), []string{"version", "data"})
	if err != nil {
		return nil, err
	}

	var data [][]string
	MustJsonUnmarshal(([]byte)(ret["data"].(string)), &data)
	v := int(ret["version"].(int64))
	fmt.Println(v, data, err)

	xlFile := newFile(data)
	tmpFile := newFile(data)
	table := &Table{
		tableName:  tableName,
		version:    v,
		tmpFile:    tmpFile,
		xlFile:     xlFile,
		sessionMap: map[string]*Session{},
		nameMap:    map[string]int{},
		cmds:       []map[string]interface{}{},
		cmdUsers:   []string{},
	}
	table.loadHistory()

	return table, nil
}

func (this *Table) loadHistory() {
	this.verHistory = map[int]*VerHistory{}
	for i := 0; i < 20; i++ {
		id := this.version - i
		if id <= 0 {
			break
		}
		users, date, cmds, err := pgsql.LoadCmd(this.tableName, id)
		if err != nil {
			fmt.Println(err)
			break
		}
		this.verHistory[id] = &VerHistory{
			version: id,
			users:   users,
			date:    date,
			cmds:    cmds,
		}
	}
}

func (this *Table) packHistory() []map[string]interface{} {
	ret := []map[string]interface{}{}
	for i := this.version; i > 0; i-- {
		v, ok := this.verHistory[i]
		if !ok {
			break
		}

		ret = append(ret, map[string]interface{}{
			"version": v.version,
			"users":   v.users,
			"date":    v.date,
		})
	}
	return ret
}

// 保存，这里指落地到数据库
func (this *Table) SaveTable() error {
	if len(this.cmds) > 0 {
		// todo 验证正确性
		doCmds(this.xlFile, this.cmds)

		// 保存版本指令
		cmdsStr, _ := json.Marshal(this.cmds)
		users, _ := json.Marshal(this.cmdUsers)
		v, date, err := pgsql.InsertCmd(this.tableName, string(users), string(cmdsStr))
		if err != nil {
			return err
		}
		this.version = v
		this.cmds = []map[string]interface{}{}
		this.cmdUsers = []string{}
		this.loadHistory()

		// 保存最新数据
		b := MustJsonMarshal(getAll(this.xlFile))
		err = pgsql.Update("table_data", fmt.Sprintf("table_name = '%s'", this.tableName), map[string]interface{}{
			"version": v,
			"date":    date,
			"data":    string(b),
		})
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

	if _, ok := this.nameMap[session.UserName]; !ok {
		this.nameMap[session.UserName] = 1
		this.pushUserEnter(session, true)
	} else {
		this.nameMap[session.UserName] += 1
		this.pushUserEnter(session, false)
	}
}

func (this *Table) RemoveSession(session *Session) {
	if _, ok := this.sessionMap[session.RemoteAddr()]; ok {
		delete(this.sessionMap, session.RemoteAddr())
	}

	if _, ok := this.nameMap[session.UserName]; ok {
		this.nameMap[session.UserName] -= 1
		//fmt.Println(this.nameMap)

		if this.nameMap[session.UserName] == 0 {
			delete(this.nameMap, session.UserName)
			this.pushOtherUser(map[string]interface{}{
				"cmd":      "userLeave",
				"userName": session.UserName,
			}, session.UserName)
		}
	}

}

func (this *Table) pushUserEnter(session *Session, isNew bool) {
	// 第一次登入，将自己同步给其他人
	if isNew {
		this.pushOtherUser(map[string]interface{}{
			"cmd":      "userEnter",
			"userName": session.UserName,
		}, session.UserName)
	}

	// 将其他人同步给自己
	for name := range this.nameMap {
		if name != session.UserName {
			b := MustJsonMarshal(map[string]interface{}{
				"cmd":      "userEnter",
				"userName": name,
			})
			_ = session.DirectSend(b)
		}
	}
}

// 转发命令
func (this *Table) pushAllSession(msg map[string]interface{}) {
	b := MustJsonMarshal(msg)
	for _, session := range this.sessionMap {
		_ = session.Send(msg["cmd"].(string), b)
	}
}

func (this *Table) pushOtherUser(msg map[string]interface{}, userName string) {
	b := MustJsonMarshal(msg)
	for _, session := range this.sessionMap {
		if session.UserName != userName {
			_ = session.Send(msg["cmd"].(string), b)
		}
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
					_ = file.SaveTable()
					//fmt.Println("delete Table", file.tableName)
					delete(tables, file.tableName)
				}
			}
		})
	}
}
