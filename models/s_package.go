package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"yklili/util/modelutil"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/**
*   pk      主键
*   auto        自增值（限数值）
*   column(N)   指定字段名N
*   null        非空
*   index       单个字段索引
*   unique      唯一键
*   auto_now_add    第一次插入数据时自动添加当前时间
*   auto_now    每一次保存时自动更新当前时间
*   type(T)     对应数据库的指定类型
*   size(S)     类型长度S
*   default(D)  默认值D（需要对应类型）
**/
type S_Package struct {
	Id         int       `orm:"pk;column(id)"`
	Title      string    `orm:"column(title);size(64)"`
	AddTime    time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser    string    `orm:"column(adduser)"`
	ModifyTime time.Time `orm:"null;type(datetime);column(modifytime)"`
	ModifyUser string    `orm:"null;column(modifyuser);size(64)"`
}

type FrontPackage struct {
	Id           int       `orm:"column(id)"`
	Title        string    `orm:"column(title)"`
	AddTime      time.Time `orm:"column(addtime)"`
	AddUser      string    `orm:"column(adduser)"`
	ModifyTime   time.Time `orm:"column(modifytime)"`
	ModifyUser   string    `orm:"column(modifyuser)"`
	ArticleCount int64     `orm:"column(articlecount)"`
}

//自定义表名
func (m *S_Package) TableName() string {
	return "s_package"
}

func (c *S_Package) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	cid, err := strconv.Atoi(tmpId)
	if err == nil {
		c.Id = cid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (c *S_Package) GetId() int {
	return c.Id
}

func (c *S_Package) GetTitle() string {
	return c.Title
}

func (c *S_Package) SetTitle(title string) {
	c.Title = title
}

func (c *S_Package) SetAddTime(t time.Time) {
	c.AddTime = t
}

func (c *S_Package) SetCurrentTime() {
	c.AddTime = time.Now()
}

func (c *S_Package) SetAddUser(uname string) {
	c.AddTime = time.Now()
	c.AddUser = uname
}

func (c *S_Package) SetModifyUser(uname string) {
	c.ModifyTime = time.Now()
	c.ModifyUser = uname
}

func (c *S_Package) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, c)
}

//判断是非存在文章
func (c *S_Package) IsNull() bool {
	o := orm.NewOrm()
	var cnt int64
	o.Raw("SELECT COUNT(1) FROM s_article WHERE packageid = ?", c.Id).QueryRow(&cnt)
	return cnt > 0
}

//查询数据库
func (c *S_Package) Fill() error {
	o := orm.NewOrm()
	if c.Id > 0 {
		return o.Read(c, "Id")
	}

	return fmt.Errorf("请确认是否传递了Id", "")

}

//插入
func (c *S_Package) Insert() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(c)
}

//修改
func (c *S_Package) Update(column ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(c, column...)
}

func init() {
	orm.RegisterModel(new(S_Package))
}

func (c *S_Package) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		beego.Warn("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

//相关函数

//根据ID得到菜单数据
func GetPackage(id int) (*S_Package, error) {
	c := S_Package{Id: id}
	err := c.Fill()
	if err != nil {
		beego.Error("文件夹Id不存在")
		return &c, err
	}
	return &c, nil
}

//得到所有的包
func GetPackages(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "addtime"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var as []*S_Package

	o := orm.NewOrm()

	qt := o.QueryTable("s_package")

	if data["User"] != nil {
		user := data["User"].(*S_User)
		qt = qt.Filter("AddUser", user.GetUserName())
		_, err := qt.OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&as)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("User cannot be empty")
	}

	cnt, err := qt.Count()
	if err != nil {
		return nil, err
	}

	pagetotal := cnt / int64(size)

	if cnt%int64(size) > 0 {
		pagetotal++
	}
	return GetDataGrid(as, index, int(pagetotal), cnt), err
}

//得到所有的包
func GetFrontPackages(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "addtime"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var as []*FrontPackage

	o := orm.NewOrm()

	cnt, err := o.Raw("SELECT a.*,COUNT(b.packageid) articlecount FROM s_package a LEFT JOIN s_article b ON a.id = b.packageid GROUP BY a.id").QueryRows(&as)

	if err != nil {
		return nil, err
	}

	pagetotal := cnt / int64(size)

	if cnt%int64(size) > 0 {
		pagetotal++
	}
	return GetDataGrid(as, index, int(pagetotal), cnt), err
}
