package test

import (
	"beegostudy/util"
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	//s := util.LeftPad("aa", 'c', 8)
	//t.Fatal(s)
	//t.Fatal(util.RandInt(1))
	fmt.Println(util.GeneratePWD("qweqwe"))
	t.Fatal("qweqwe")
}
