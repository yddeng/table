package main

import (
	"fmt"
)

func InsertGenID(tableName string, colID, rowID int) error {
	sqlStr := `
INSERT INTO gen_id (table_name, col_id,row_id)
VALUES ('%s','%d','%d');`
	sqlStatement := fmt.Sprintf(sqlStr, tableName, colID, rowID)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

func UpdateGenID(tableName string, key string, value int) error {
	sqlStr := `
UPDATE gen_id 
SET %s = '%d'
WHERE table_name = '%s';`
	sqlStatement := fmt.Sprintf(sqlStr, key, value, tableName)
	smt, err := dbConn.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

// 加载colid，rowId
func LoadGenID(tableName string) (int, int, error) {
	sqlStatement := `
SELECT col_id,row_id FROM gen_id 
WHERE table_name=$1;`

	var colID, rowID int
	row := dbConn.QueryRow(sqlStatement, tableName)
	err := row.Scan(&colID, &rowID)
	if err != nil {
		return 0, 0, err
	}

	return colID, rowID, nil
}
