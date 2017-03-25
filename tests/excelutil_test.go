package test

import (
	"fmt"
	"github.com/sinmahod/yklili/util/fileutil"
	"testing"
)

type ExcelData struct {
	Name string `tag:"姓名"`
}

func Test_ReadXLSX(t *testing.T) {
	ss, _ := fileutil.ReadXLSXToDT("../1.xlsx")
	fmt.Println(ss)
	// ss := excels[0]
	// fmt.Println(ss)
	t.Fatal("")
}

func Test_WriteXLSX(t *testing.T) {
	//fileutil.WriteXLSX_Test()
	t.Fatal("")
}
