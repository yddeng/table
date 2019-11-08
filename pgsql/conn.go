package pgsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

var dbConn *sql.DB

func sqlOpen() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	dbConn = db

	fmt.Println("Successfully connected!")
	return nil
}

func createTable(tableName string) {
	//client.db.
}

func insert(tableName string, fields map[string]interface{}) {

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
	fmt.Println(sqlStatement)
	row, err := dbConn.Query(sqlStatement, args...)

	fmt.Println(row, err)
	var id int
	err = row.Scan(&id)
	fmt.Println(id, err)
}

func query(tableName string) {
	sqlStr := fmt.Sprintf(`SELECT * FROM %s;`, tableName)
	rows, err := dbConn.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	out := []string{}
	for rows.Next() {
		var id string
		var age int
		err := rows.Scan(&age, &id)
		if err != nil {
			panic(err)
		}
		out = append(out, fmt.Sprintf("%s:%d", id, age))
	}
	for _, v := range out {
		fmt.Println(v)
	}
}
