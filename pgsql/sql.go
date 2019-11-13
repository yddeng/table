package pgsql

import (
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
        "username"  varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
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
func InsertCmd(tableName, userName, cmd string) (int, error) {
	sqlStr := `
INSERT INTO %s_cmd (username,date,cmds)  
VALUES ('%s','%s','%s')
RETURNING version;`

	date := GenDateTimeString(time.Now())
	sqlStatement := fmt.Sprintf(sqlStr, tableName, userName, date, cmd)
	row := dbConn.QueryRow(sqlStatement)
	var id int
	err := row.Scan(&id)
	return id, err
}

func LoadCmd(tableName string, v int) ([]map[string]interface{}, error) {
	sqlStr := `
SELECT cmds FROM %s_cmd 
WHERE version=%d;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName, v)
	var cmdStr string
	row := dbConn.QueryRow(sqlStatement)
	err := row.Scan(&cmdStr)
	if err != nil {
		return nil, err
	}

	var cmds []map[string]interface{}
	err = json.Unmarshal(([]byte)(cmdStr), &cmds)
	return cmds, err
}

func Insert(tableName string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO %s (%s)
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

/************************************** table_data ********************************************************/

func InsertTableData(key string, version int, data string) error {
	sqlStr := `
INSERT INTO table_data (table_name, now_version,data)
VALUES ('%s','%d','%s');`
	sqlStatement := fmt.Sprintf(sqlStr, key, version, data)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

func UpdateTableData(key string, version int, data string) error {
	sqlStr := `
UPDATE table_data 
SET now_version = '%d',data = '%s'
WHERE table_name = '%s';`
	sqlStatement := fmt.Sprintf(sqlStr, version, data, key)
	//fmt.Println(sqlStatement)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

// 加载表数据
func LoadTableData(tableName string) (int, [][]string, error) {
	sqlStatement := `
SELECT * FROM table_data 
WHERE table_name=$1;`

	var name, dataStr string
	var version int
	row := dbConn.QueryRow(sqlStatement, tableName)
	err := row.Scan(&name, &version, &dataStr)
	if err != nil {
		return 0, nil, err
	}

	var data [][]string
	err = json.Unmarshal(([]byte)(dataStr), &data)
	if err != nil {
		return 0, nil, err
	}
	return version, data, nil
}

func AllTableData() ([]map[string]interface{}, error) {
	rows, err := dbConn.Query(`SELECT table_name,now_version FROM table_data;`)
	if err != nil {
		return nil, err
	}

	ret := []map[string]interface{}{}
	for rows.Next() {
		var tableName string
		var version int
		err := rows.Scan(&tableName, &version)
		if err != nil {
			return nil, err
		}
		ret = append(ret, map[string]interface{}{
			"tableName": tableName,
			"version":   version,
		})
	}

	return ret, nil
}
