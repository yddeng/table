package table

import (
	"encoding/json"
	"fmt"
	"time"
)

func Min(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("panic", err)
		panic(err)
	}
}

func Must(msg interface{}, err error) interface{} {
	CheckErr(err)
	return msg
}

func MustJsonMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	CheckErr(err)
	return b
}

func MustJsonUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	CheckErr(err)
}

// 生成日期字符串
func GenDateTimeString(date time.Time) string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d",
		date.Year(), int(date.Month()), date.Day(), date.Hour(), date.Minute(), date.Second())
}
