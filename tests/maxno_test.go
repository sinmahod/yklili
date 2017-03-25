package test

import (
	"testing"
	"yklili/models"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(60.205.164.3:3306)/beestudy?charset=utf8")
}

func Test_GetMaxNo(t *testing.T) {
	orm.Debug = true
	t.Fatal(models.GetMaxNo("menu", "", 4))
}
