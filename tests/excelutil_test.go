package test

import (
	"beegostudy/util/fileutil"
	"testing"
)

func Test_ReadExcel(t *testing.T) {
	fileutil.ReadExcel2("/Users/gl/Downloads/1.xlsx")
	t.Fatal("")
}
