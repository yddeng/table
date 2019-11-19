package pgsql

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// 生成一条初始化数据
func InsertTableData(key string, desc string) error {
	sqlStr := `
INSERT INTO table_data (table_name, version,data,date,describe)
VALUES ('%s','%d','%s','%s','%s');`

	b, _ := json.Marshal([]string{})
	date := GenDateTimeString(time.Now())
	sqlStatement := fmt.Sprintf(sqlStr, key, 0, string(b), date, desc)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

func UpdateTableData(key string, fields map[string]interface{}) error {
	sqlStr := `
UPDATE table_data 
SET %s
WHERE table_name = '%s';`

	keys := []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		keys = append(keys, fmt.Sprintf(`%s = $%d`, k, i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, strings.Join(keys, ","), key)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

// 加载表数据,describe,version,date,data
func LoadTableData(tableName string) (string, int, string, [][]string, error) {
	sqlStatement := `
SELECT describe,version,date,data FROM table_data 
WHERE table_name=$1;`

	var desc, date, dataStr string
	var version int
	row := dbConn.QueryRow(sqlStatement, tableName)
	err := row.Scan(&desc, &version, &date, &dataStr)
	if err != nil {
		return "", 0, "", nil, err
	}

	var data [][]string
	err = json.Unmarshal(([]byte)(dataStr), &data)
	if err != nil {
		return "", 0, "", nil, err
	}
	return desc, version, date, data, nil
}

func AllTableData() ([]map[string]interface{}, error) {
	rows, err := dbConn.Query(`SELECT table_name,describe,version,date FROM table_data;`)
	if err != nil {
		return nil, err
	}

	ret := []map[string]interface{}{}
	for rows.Next() {
		var name, desc, date string
		var version int
		err := rows.Scan(&name, &desc, &version, &date)
		if err != nil {
			return nil, err
		}
		ret = append(ret, map[string]interface{}{
			"tableName": name,
			"version":   version,
			"describe":  desc,
			"date":      date,
		})
	}

	return ret, nil
}
