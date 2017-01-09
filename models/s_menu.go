package models

import (
	"beegostudy/util/modelutil"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
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
type S_Menu struct {
	Id         int       `orm:"pk;column(id)"`
	Pid        int       `orm:"column(pid)"`
	MenuName   string    `orm:"column(menuname);size(64)"`
	Icon       string    `orm:"null;column(icon);size(32)"`
	IsLeaf     bool      `orm:"-"`
	Link       string    `orm:"null;column(link);size(128)"`
	InnerCode  string    `orm:"column(innercode);size(128)"`
	Level      int       `orm:"column(level)"`
	OrderFlag  int       `orm:"column(orderflag)"`
	AddTime    time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser    string    `orm:"column(adduser)"`
	ModifyTime time.Time `orm:"type(datetime);column(modifytime)"`
	ModifyUser string    `orm:"column(modifyuser);size(64)"`
	Checked    bool      `orm:"-"`
	Expanded   bool      `orm:"-"`
	ChildNode  []*S_Menu `orm:"-"`
}

//自定义表名
func (m *S_Menu) TableName() string {
	return "s_menu"
}

func (menu *S_Menu) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	menuid, err := strconv.Atoi(tmpId)
	if err == nil {
		menu.Id = menuid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (menu *S_Menu) GetId() int {
	return menu.Id
}

func (menu *S_Menu) GetPid() int {
	return menu.Pid
}

func (menu *S_Menu) SetLevel(l int) {
	menu.Level = l
}

func (menu *S_Menu) SetInnerCode(code string) {
	menu.InnerCode = code
}

func (menu *S_Menu) GetInnerCode() string {
	return menu.InnerCode
}

//是否为叶子节点
func (menu *S_Menu) GetIsLeaf() bool {
	o := orm.NewOrm()
	var count int
	o.Raw("SELECT COUNT(*) FROM s_menu WHERE pid=?", menu.Id).QueryRow(&count)
	return count == 0
}

func (menu *S_Menu) SetAddTime(t time.Time) {
	menu.AddTime = t
}

func (menu *S_Menu) SetCurrentTime() {
	menu.AddTime = time.Now()
}

func (menu *S_Menu) SetAddUser(uname string) {
	menu.AddTime = time.Now()
	menu.AddUser = uname
}

func (menu *S_Menu) SetModifyUser(uname string) {
	menu.ModifyTime = time.Now()
	menu.ModifyUser = uname
}

func (menu *S_Menu) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, menu)
}

//查询数据库
func (menu *S_Menu) Fill() error {
	o := orm.NewOrm()
	if menu.Id > 0 {
		return o.Read(menu, "Id")
	}

	return fmt.Errorf("请确认是否传递了Id", "")

}

//插入
func (menu *S_Menu) Insert() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(menu)
}

//修改
func (menu *S_Menu) Update(column ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(menu, column...)
}

func init() {
	orm.RegisterModel(new(S_Menu))
}

func (menu *S_Menu) String() string {
	data, err := json.MarshalIndent(menu, "", "    ")
	if err != nil {
		beego.Warn("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

//相关函数

//返回第一个url
func GetIndexLink() string {
	o := orm.NewOrm()
	var menu S_Menu
	err := o.QueryTable("s_menu").Exclude("link", "").OrderBy("innercode").One(&menu, "link")
	if err == orm.ErrNoRows {
		// 没有找到记录
		return ""
	}
	return menu.Link
}

//根据ID得到菜单数据
func GetMenu(id int) (*S_Menu, error) {
	menu := S_Menu{Id: id}
	err := menu.Fill()
	if err != nil {
		beego.Error("菜单Id不存在")
		return &menu, err
	}
	return &menu, nil
}

//得到父级InnerCode
func GetInnerCode(id int) string {
	o := orm.NewOrm()
	var menu S_Menu
	err := o.QueryTable("s_menu").Filter("Id", id).One(&menu, "innercode")
	if err == orm.ErrNoRows {
		// 没有找到记录
		return ""
	}
	return menu.InnerCode
}

//得到所有的菜单
func GetMenus() ([]*S_Menu, error) {
	var menus []*S_Menu
	o := orm.NewOrm()
	_, err := o.QueryTable("s_menu").OrderBy("innercode").All(&menus)
	return menus, err
}

//得到所有菜单并按级别排序
func GetMenusLevel(url string) ([]*S_Menu, error) {
	menus, err := GetMenus()
	if err == nil && menus != nil {
		var menuslevel []*S_Menu = make([]*S_Menu, 0)
		//用来记录menuslevel的当前位置
		idx := -1
		//线型遍历一遍
		for _, menu := range menus {
			menu.IsLeaf = true //默认全部是叶子节点
			if strings.EqualFold(menu.Link, url) {
				menu.Checked = true
			}
			if menu.Pid == 0 {
				idx++
				menuslevel = append(menuslevel, menu)
				continue
			}

			if menu.Pid == menuslevel[idx].Id {
				//如果当前元素是menuslevel第idx个元素的子集时就放入ChildNode中
				if menu.Checked {
					menuslevel[idx].Checked = true
					menuslevel[idx].IsLeaf = false //如果有子级则不为叶子节点
				}
				menuslevel[idx].ChildNode = append(menuslevel[idx].ChildNode, menu)
			}
			//循环到了这里的元素都是找不到父级的元素这里直接丢弃
		}
		return menuslevel, nil
	}
	return menus, err
}

//得到分页的菜单
/**
*	size	每页查询长度
*	index	查询的页码
*	ordercolumn	排序字段
*	orderby		升降序:desc\asc
**/
func GetMenusPage(size, index int, ordercolumn, orderby string) (*DataGrid, error) {
	if ordercolumn == "" {
		ordercolumn = "innercode"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}
	var menus []*S_Menu
	o := orm.NewOrm()
	_, err := o.QueryTable("s_menu").OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&menus)
	if err == nil {
		cnt, err := o.QueryTable("s_menu").Count()
		pagetotal := cnt / int64(size)
		if cnt%int64(size) > 0 {
			pagetotal++
		}
		var tempmenus []*S_Menu = make([]*S_Menu, len(menus))
		for i, menu := range menus {
			//展开节点
			menu.Expanded = true
			//设置默认都是叶子节点
			menu.IsLeaf = true
			tempmenus[i] = menu

			if menu.Level == 2 {
				for j := i - 1; j >= 0; j-- {
					if menu.Pid == tempmenus[j].Id {
						tempmenus[j].IsLeaf = false //如果有子级则不为叶子节点
						break
					}
				}
			}
		}

		return GetDataGrid(tempmenus, index, int(pagetotal), cnt), err
	}

	return nil, err
}

//返回对应级别的菜单组
func GetMenusByLevel(level int) ([]*S_Menu, error) {
	var menus []*S_Menu
	o := orm.NewOrm()
	_, err := o.QueryTable("s_menu").Filter("Level", level).OrderBy("innercode").All(&menus)
	return menus, err
}

type MenuSelectInit struct {
	S_Menu
	Select string
}

//返回所有顶级菜单，并指定当前父级和去掉自身
func GetTopMenus(pid, self int) ([]MenuSelectInit, error) {
	menus, err := GetMenusByLevel(1)
	if err != nil {
		return nil, err
	}

	menuselectinit := make([]MenuSelectInit, 0)
	for _, menu := range menus {
		if menu.Id == self {
			continue
		}
		selectInit := new(MenuSelectInit)
		selectInit.S_Menu.Id = menu.Id
		selectInit.S_Menu.MenuName = menu.MenuName
		if menu.Id == pid {
			selectInit.Select = "selected"
		}
		menuselectinit = append(menuselectinit, *selectInit)
	}
	return menuselectinit, nil
}
