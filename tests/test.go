package test

import (
	_ "beegostudy/models"
	_ "github.com/go-sql-driver/mysql"
	//"beegostudy/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
)

func TestMain(t *testing.T) {
	//s := util.LeftPad("aa", 'c', 8)
	//t.Fatal(s)
	//t.Fatal(util.RandInt(1))
	//fmt.Println(util.GeneratePWD("qweqwe"))
	//t.Fatal("qweqwe")
}

func TestCreateTable(t *testing.T) {
	/**************自动建表***********/
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(60.205.164.3:3306)/beestudy?charset=utf8")
	//数据库别名
	name := "default"

	// drop table 后再建表
	force := false

	// 打印执行过程
	verbose := true

	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}

	t.Fatal("Success")
}
