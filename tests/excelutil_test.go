package test

import (
	"beegostudy/util/fileutil"
	"testing"
)

func Test_ReadExcel(t *testing.T) {
	fileutil.ReadExcel2("E:\\mywork\\go\\src\\beegostudy\\1.xlsx")
	t.Fatal("")
}
