package table

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"sync"
)

type ExcelFile struct {
	fileName    string
	xlFile      *excelize.File
	activeSheet string
	users       []*User
	locked      map[string]*User
	mu          sync.Mutex
}

var excelFiles = map[string]*ExcelFile{}

func CreateExcel(fileName string) error {
	file := excelize.NewFile()
	index := file.NewSheet("Sheet1")
	file.SetActiveSheet(index)
	return file.SaveAs(fmt.Sprintf("%s.xlsx", fileName))
}

func OpenExcel(fileName string) (*ExcelFile, error) {
	xlFile, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	return &ExcelFile{
		fileName:    fileName,
		xlFile:      xlFile,
		activeSheet: xlFile.GetSheetName(xlFile.GetActiveSheetIndex()),
		users:       []*User{},
		locked:      map[string]*User{},
		mu:          sync.Mutex{},
	}, nil
}

func GetExcel(fileName string) *ExcelFile {
	return excelFiles[fileName]
}

func (ef *ExcelFile) Save() error {
	return ef.xlFile.Save()
}

func (ef *ExcelFile) GetAll() ([][]string, error) {
	return ef.xlFile.GetRows(ef.activeSheet)
}

func (ef *ExcelFile) GetRow(idx int) ([]string, error) {
	rows, err := ef.xlFile.GetRows(ef.activeSheet)
	if err != nil {
		return nil, err
	}
	return rows[idx], nil
}

func (ef *ExcelFile) GetCell(axis string) (string, error) {
	return ef.xlFile.GetCellValue(ef.activeSheet, axis)
}

func (ef *ExcelFile) SetValues(fields map[string]string) {
	for k, v := range fields {
		_ = ef.xlFile.SetCellValue(ef.activeSheet, k, v)
	}
}

func (ef *ExcelFile) ClickDownCell(user *User, axis string) error {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	if u, ok := ef.locked[axis]; ok {
		if u != user {
			return fmt.Errorf("他人正在编辑")
		}
	} else {
		ef.locked[axis] = user
	}
	return nil
}

func (ef *ExcelFile) ClickUpCell(user *User, axis string) {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	delete(ef.locked, axis)
}

func (ef *ExcelFile) WriteCell(user *User, axis string, value string) error {
	ef.mu.Lock()
	defer ef.mu.Unlock()
	if u, ok := ef.locked[axis]; ok {
		if u != user {
			return fmt.Errorf("他人正在编辑")
		}
	}
	return ef.xlFile.SetCellValue(ef.activeSheet, axis, value)
}
