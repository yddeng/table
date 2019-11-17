package table

import "encoding/json"

func Min(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func CheckErr(err error) {
	if err != nil {
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
	if err != nil {
		panic(err)
	}
}
