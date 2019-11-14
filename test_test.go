package table

import (
	"fmt"
	"testing"
)

func TestOpenFile(t *testing.T) {
	file := OpenFile("")
	fmt.Println(file)
	fmt.Println(file.xlFile.GetRows(Sheet))

	file.InsertCol(1, 4)
	fmt.Println(file)
	fmt.Println(file.xlFile.GetRows(Sheet))
	file.RemoveCol(0)
	fmt.Println(file)
	fmt.Println(file.xlFile.GetRows(Sheet))

	err := file.SetCellValue(4, 4, "14")
	fmt.Println(file, err)
	fmt.Println(file.xlFile.GetRows(Sheet))
}
