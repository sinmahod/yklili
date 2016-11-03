package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

/**
*	pk		主键
*	auto 		自增值（限数值）
*	column(N)	指定字段名N
*	null		非空
*	index 		单个字段索引
* 	unique 		唯一键
* 	auto_now_add 	第一次插入数据时自动添加当前时间
* 	auto_now 	每一次保存时自动更新当前时间
* 	type(T)		对应数据库的指定类型
*	size(S)		类型长度S
*	default(D)	默认值D（需要对应类型）
**/
type Menu struct {
	Id        int       `orm:"pk;auto;column(id)"`
	Pid       int       `orm:"column(pid)"`
	MenuName  string    `orm:"column(menuname);size(64)"`
	Icon      string    `orm:"null;column(icon);size(32)"`
	IsRoot    bool      `orm:"column(isroot);default(true)"`
	Link      string    `orm:"null;column(link);size(128)"`
	OrderFlag int       `orm:"column(orderflag)"`
	AddTime   time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser   string    `orm:"column(adduser)"`
}

//自定义表名
func (m *Menu) TableName() string {
	return "menu"
}

func (menu *Menu) SetID(id int) {
	menu.Id = id
}

func (menu *Menu) GetID() int {
	return menu.Id
}

func init() {
	orm.RegisterModel(new(Menu))
}

func (menu *Menu) String() string {
	return fmt.Sprintf("{Menu:{Id:%d,MenuName:'%s',AddTime:'%s',AddUser:'%s'}}", menu.Id, menu.MenuName, menu.AddTime, menu.AddUser)
}
