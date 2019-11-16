package table

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
