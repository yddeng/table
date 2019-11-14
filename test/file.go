package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
)

type File struct {
	FileName string
	ColList  []int
	RowList  []int
	Data     map[string]string // key=row:col
	xlFile   *excelize.File
}

func genKey(row, col int) string {
	return fmt.Sprintf("%d@%d", row, col)
}

//  return row,col
func keyToInt(key string) (int, int) {
	s := strings.Split(key, "@")
	row, _ := strconv.Atoi(s[0])
	col, _ := strconv.Atoi(s[1])
	return row, col
}

// 检测行列是否存在
func (this *File) checkRowCol(rowId, colId int) (bool, bool) {
	var rowOk, colOk = false, false
	for _, row := range this.RowList {
		if row == rowId {
			rowOk = true
			break
		}
	}

	for _, col := range this.ColList {
		if col == colId {
			colOk = true
			break
		}
	}

	return rowOk, colOk
}

func (this *File) InsertCol(idx, col int) {
	array := []int{col}
	array = append(array, this.ColList[idx:]...)
	this.ColList = append(this.ColList[:idx], array...)
	celHeader, _ := excelize.ColumnNumberToName(idx + 1)
	_ = this.xlFile.InsertCol(Sheet, celHeader)
}

func (this *File) RemoveCol(idx int) {
	this.ColList = append(this.ColList[:idx], this.ColList[idx+1:]...)
	celHeader, _ := excelize.ColumnNumberToName(idx + 1)
	_ = this.xlFile.RemoveCol(Sheet, celHeader)
}

func (this *File) InsertRow(idx, row int) {
	array := []int{row}
	array = append(array, this.RowList[idx:]...)
	this.RowList = append(this.RowList[:idx], array...)
	_ = this.xlFile.InsertRow(Sheet, idx)
}

func (this *File) RemoveRow(idx int) {
	this.RowList = append(this.RowList[:idx], this.RowList[idx+1:]...)
	_ = this.xlFile.InsertRow(Sheet, idx)
}

// 真实的行列序号
func (this *File) SetCellValue(rowId, colId int, value string) error {
	co, ro := this.checkRowCol(rowId, colId)
	if co && ro {
		key := genKey(rowId, colId)
		if value == "" {
			delete(this.Data, key)
		} else {
			this.Data[key] = value
		}
		return nil
	} else {
		return fmt.Errorf("不存在的行或列")
	}
}

func NewFile() *File {
	return &File{
		ColList: []int{},
		RowList: []int{},
		Data:    map[string]string{},
	}
}

func OpenFile(FileName string) *File {
	file := &File{
		ColList: []int{2, 1, 3},
		RowList: []int{1, 2, 3},
		Data: map[string]string{
			"1@2": "00",
			"2@1": "11",
		},
	}

	xF := newFile()
	for i, row := range file.RowList {
		for j, col := range file.ColList {
			cellName, _ := excelize.CoordinatesToCellName(j+1, i+1)
			value, ok := file.Data[genKey(row, col)]
			if ok {
				xF.SetCellStr(Sheet, cellName, value)
			}
		}
	}
	file.xlFile = xF

	return file
}
