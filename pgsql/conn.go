package pgsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "dbuser"
	password = "123456"
	dbname   = "excel"
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
	return nil
}

// 生成日期字符串
func GenDateTimeString(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d",
		date.Year(), int(date.Month()), date.Day(), date.Hour(), date.Minute(), date.Second())
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

func init() {
	err := sqlOpen()
	if err != nil {
		panic(err)
	}
}
