package pgsql

import (
	"encoding/json"
	"fmt"
)

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
