package table

import (
	"fmt"
	"github.com/yddeng/table/conf"
	"net/http"
)

func Start(path string) {
	//table.InitLogger()

	file, err := OpenExcel("test.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	excelFiles[file.fileName] = file
	c1, _ := file.GetCell("A2")
	fmt.Println("cellValue1", c1)
	//file.WriteCell(nil, "A2", "A2")
	file.xlFile.SetCellValue("Sheet1", "A2", "A2")
	c2, _ := file.GetCell("A2")
	fmt.Println("cellValue2", c2)

	conf.LoadConfig(path)
	_conf := conf.GetConfig()
	fmt.Printf("start on %s, LoadDir on %s\n", _conf.HttpAddr, _conf.LoadDir)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(_conf.LoadDir))))
	http.HandleFunc("/setCellValue", cellValue)

	err = http.ListenAndServe(_conf.HttpAddr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
