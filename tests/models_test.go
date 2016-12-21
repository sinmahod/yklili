package test

import (
	"beegostudy/models"
	"beegostudy/util/numberutil"
	"beegostudy/util/pwdutil"
	"beegostudy/util/stringutil"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:qweqwe@tcp(60.205.164.3:3306)/beestudy?charset=utf8")
}

func Test_MenuString(t *testing.T) {
	m := new(models.Menu)
	m.Id = 1
	o := orm.NewOrm()
	o.Read(m, "Id")
	t.Logf("%s", m)
}

func Test_UserString(t *testing.T) {
	u := new(models.User)
	u.Id = 1
	o := orm.NewOrm()
	o.Read(u, "Id")
	t.Logf("%s", u)
}

func Test_LeftPad(t *testing.T) {
	s := stringutil.LeftPad("aa", 'c', 8)
	t.Fatal(s)
	t.Fatal(numberutil.RandInt(1))
	fmt.Println(pwdutil.GeneratePWD("qweqwe"))
	t.Fatal("qweqwe")
}

func Test_RunSyncdb(t *testing.T) {
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
