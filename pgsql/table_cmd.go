package pgsql

import (
	"encoding/json"
	"fmt"
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
