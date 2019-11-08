package pgsql

import (
	"fmt"
	"sync"
)

type Field struct {
	Name    string
	Value   interface{}
	fieldMu sync.Mutex
}

type Row struct {
	Key    interface{}
	Fields []*Field
	rowMu  sync.Mutex
}

type Design struct {
	Name   string
	Type   interface{}
	DefV   interface{}
	NotNil bool
}

type Table struct {
	TableName string
	Designs   []*Design
	Rows      []*Row
	tableMu   sync.Mutex
}

type Database struct {
	DBName string
	Tables map[string]*Table
	dbMu   sync.Mutex
}

var (
	dbInstance *Database
)

func LoadTable(tableName string) (*Table, error) {
	dbInstance.dbMu.Lock()
	table, ok := dbInstance.Tables[tableName]
	if ok {
		return table, nil
	}
	return nil, fmt.Errorf("not found table:%s", tableName)
}

func Insert(tableName string, feilds map[string]*Field) {
	table, err := LoadTable(tableName)
	if err != nil {
		return
	}

	row := &Row{
		Key:    "1",
		Fields: make([]*Field, len(table.Designs)),
		rowMu:  sync.Mutex{},
	}
	for i, de := range table.Designs {
		f, ok := feilds[de.Name]
		if !ok {
			row.Fields[i] = &Field{
				Name:    de.Name,
				Value:   de.DefV,
				fieldMu: sync.Mutex{},
			}
		} else {
			row.Fields[i] = &Field{
				Name:    de.Name,
				Value:   f.Value,
				fieldMu: sync.Mutex{},
			}
		}
	}

	table.Rows = append(table.Rows, row)
}
