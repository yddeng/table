package table

import "encoding/json"

func Min(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func Must(msg interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return msg
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MustJsonMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func MustJsonUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}
