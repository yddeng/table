package pgsql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/yddeng/table/conf"
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
	_conf := conf.GetConfig().DBConfig
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		_conf.DbHost, _conf.DbPort, _conf.DbUser, _conf.DbPassword, _conf.DbDataBase)
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

func Init() {
	err := sqlOpen()
	if err != nil {
		panic(err)
	}
}
