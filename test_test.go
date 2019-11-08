package table

import (
	"fmt"
	"github.com/yddeng/dutil/dhttp"
	"testing"
)

func TestCreateExcel(t *testing.T) {
	CreateExcel("test1")
	ef, _ := OpenExcel("test1.xlsx")
	ef.WriteCell(nil, "A2", "A2")
	ef.Save()
}

func TestStart(t *testing.T) {
	_, err := dhttp.Get("http://127.0.0.1:4545/setCellValue?axis=B2&value=this", 0)
	fmt.Println(err)
}
