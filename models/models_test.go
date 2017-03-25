package models

import (
	// "fmt"
	// "github.com/astaxie/beego/orm"
	// _ "github.com/go-sql-driver/mysql"
	// "github.com/sinmahod/yklili/util/numberutil"
	// "github.com/sinmahod/yklili/util/pwdutil"
	// "github.com/sinmahod/yklili/util/stringutil"
	"testing"
)

func Test_RunSyncdb(t *testing.T) {
	/**************自动建表***********/
	// orm.RegisterDriver("mysql", orm.DRMySQL)

	// orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(localhost:3306)/beestudy?charset=utf8")
	// //数据库别名
	// name := "default"

	// // drop table 后再建表
	// force := false

	// // 打印执行过程
	// verbose := true

	// // 遇到错误立即返回
	// err := orm.RunSyncdb(name, force, verbose)
	// if err != nil {
	// 	t.Error(err)
	// }

	t.Log("Success")
}
