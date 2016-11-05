package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
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
	InnerCode string    `orm:"column(innercode);size(128)"`
	OrderFlag int       `orm:"column(orderflag)"`
	AddTime   time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser   string    `orm:"column(adduser)"`
	Checked   bool      `orm:"-"`
	ChildNode []*Menu   `orm:"-"`
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

//相关函数

//得到所有的菜单
func GetMenus() ([]*Menu, error) {
	var menus []*Menu
	o := orm.NewOrm()
	_, err := o.QueryTable("menu").OrderBy("innercode").All(&menus)
	return menus, err
}

//得到所有菜单并按级别排序
func GetMenusLevel(url string) ([]*Menu, error) {
	menus, err := GetMenus()
	if err == nil && menus != nil {
		var menuslevel []*Menu = make([]*Menu, len(menus), cap(menus))
		//用来记录menuslevel的当前位置
		idx := 0
		//线型遍历一遍
		for _, menu := range menus {
			if strings.EqualFold(menu.Link, url) {
				menu.Checked = true
			}
			if menuslevel[idx] == nil {
				menuslevel[idx] = menu
				continue
			}
			if menu.Pid == 0 {
				idx++
				menuslevel[idx] = menu
				continue
			}
			if menu.Pid == menuslevel[idx].Id {
				//如果当前元素是menuslevel第idx个元素的子集时就放入ChildNode中
				if menu.Checked {
					menuslevel[idx].Checked = true
				}
				menuslevel[idx].ChildNode = append(menuslevel[idx].ChildNode, menu)
			}
			//循环到了这里的元素都是找不到父级的元素这里直接丢弃
		}
		return menuslevel[:idx+1], nil
	}
	return menus, err
}
