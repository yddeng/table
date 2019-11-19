package pgsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// 操作文件
func CreateTableCmd(tableName string) error {
	sql := `
    CREATE TABLE "public"."` + tableName + `_cmd" (
        "version"   SERIAL primary key ,
        "users"     varchar(128) COLLATE "pg_catalog"."default" NOT NULL,
        "date"      varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
        "cmds"      varchar(65535) COLLATE "pg_catalog"."default" NOT NULL
    );`
	smt, err := dbConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

// 插入操作
func InsertCmd(tableName, userName, cmd string) (int, string, error) {
	sqlStr := `
INSERT INTO "%s_cmd" (users,date,cmds)  
VALUES ('%s','%s','%s')
RETURNING version;`

	date := GenDateTimeString(time.Now())
	sqlStatement := fmt.Sprintf(sqlStr, tableName, userName, date, cmd)
	//fmt.Println(sqlStatement)
	row := dbConn.QueryRow(sqlStatement)
	var id int
	err := row.Scan(&id)
	return id, date, err
}

func LoadCmd(tableName string, v int) ([]string, string, []map[string]interface{}, error) {
	sqlStr := `
SELECT users,date,cmds FROM "%s_cmd" 
WHERE version=%d;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName, v)
	var userStr, dateStr, cmdStr string
	row := dbConn.QueryRow(sqlStatement)
	err := row.Scan(&userStr, &dateStr, &cmdStr)
	if err != nil {
		return nil, "", nil, err
	}

	var users []string
	_ = json.Unmarshal(([]byte)(userStr), &users)
	var cmds []map[string]interface{}
	err = json.Unmarshal(([]byte)(cmdStr), &cmds)
	return users, dateStr, cmds, nil
}

func Insert(tableName string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (%s)
VALUES (%s);`

	keys, values := []string{}, []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		keys = append(keys, k)
		values = append(values, fmt.Sprintf("$%d", i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(keys, ","), strings.Join(values, ","))
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

type SelectClient struct {
	sql    string
	fields []string
	smt    *sql.Stmt
}

func (this *SelectClient) Query(rb func(map[string]interface{}, error), args ...interface{}) {
	row := this.smt.QueryRow(args...)
	mid := make([]interface{}, len(this.fields))
	l := "&mid[0],&mid[1],&mid[2]"
	err := row.Scan(l)
	if err != nil {
		rb(nil, err)
		return
	}

	ret := map[string]interface{}{}
	for i, v := range this.fields {
		ret[v] = mid[i]
	}
	rb(ret, nil)
}

func NewSelect(tableName, fields string, key ...string) (*SelectClient, error) {
	tmp := `
SELECT %s FROM %s
`
	sqlStr := ""
	if len(key) > 0 {
		tmp += `WHERE %s=$1;`
		sqlStr = fmt.Sprintf(tmp, fields, tableName, key[0])
	} else {
		tmp += `;`
		sqlStr = fmt.Sprintf(tmp, fields, tableName)
	}

	smt, err := dbConn.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	fields_ := strings.Split(fields, ",")
	return &SelectClient{
		sql:    sqlStr,
		fields: fields_,
		smt:    smt,
	}, nil
}
