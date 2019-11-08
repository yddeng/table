package pgsql

import (
	"fmt"
	"testing"
)

func TestLoadTable(t *testing.T) {
	err := sqlOpen()
	if err != nil {
		panic(err)
	}

	feilds := map[string]interface{}{
		"id":  16,
		"age": 22,
	}

	insert("test", feilds)

	fmt.Println("sql ok ")
}

func TestInsert(t *testing.T) {
	err := sqlOpen()
	if err != nil {
		panic(err)
	}

	query("test")

	fmt.Println("sql ok ")
}
