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
type S_Catalog struct {
	Id          int          `orm:"pk;column(id)"`
	Pid         int          `orm:"column(pid)"`
	CatalogName string       `orm:"column(catalogname);size(64)"`
	Logo        string       `orm:"null;column(logo);size(128)"`
	IsLeaf      bool         `orm:"-"`
	Link        string       `orm:"null;column(link);size(256)"`
	InnerCode   string       `orm:"column(innercode);size(128)"`
	Level       int          `orm:"column(level)"`
	PreviousId  int          `orm:"column(previousid)"`
	AddTime     time.Time    `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser     string       `orm:"column(adduser)"`
	ModifyTime  time.Time    `orm:"null;type(datetime);column(modifytime)"`
	ModifyUser  string       `orm:"null;column(modifyuser);size(64)"`
	Checked     bool         `orm:"-"`
	Expanded    bool         `orm:"-"`
	ChildNode   []*S_Catalog `orm:"-"`
}

//自定义表名
func (m *S_Catalog) TableName() string {
	return "s_catalog"
}

func (c *S_Catalog) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	cid, err := strconv.Atoi(tmpId)
	if err == nil {
		c.Id = cid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (c *S_Catalog) GetId() int {
	return c.Id
}

func (c *S_Catalog) GetPid() int {
	return c.Pid
}

func (c *S_Catalog) GetCatalogName() string {
	return c.CatalogName
}

func (c *S_Catalog) GetLink() string {
	return c.Link
}

func (c *S_Catalog) SetLevel(l int) {
	c.Level = l
}

func (c *S_Catalog) GetLevel() int {
	return c.Level
}

func (c *S_Catalog) SetInnerCode(code string) {
	c.InnerCode = code
}

func (c *S_Catalog) GetInnerCode() string {
	return c.InnerCode
}

func (c *S_Catalog) SetPreviousId(preid int) {
	c.PreviousId = preid
}

func (c *S_Catalog) GetPreviousId() int {
	return c.PreviousId
}

//是否为叶子节点
func (c *S_Catalog) GetIsLeaf() bool {
	o := orm.NewOrm()
	var count int
	o.Raw("SELECT COUNT(*) FROM s_catalog WHERE pid=?", c.Id).QueryRow(&count)
	return count == 0
}

func (c *S_Catalog) SetAddTime(t time.Time) {
	c.AddTime = t
}

func (c *S_Catalog) SetCurrentTime() {
	c.AddTime = time.Now()
}

func (c *S_Catalog) SetAddUser(uname string) {
	c.AddTime = time.Now()
	c.AddUser = uname
}

func (c *S_Catalog) SetModifyUser(uname string) {
	c.ModifyTime = time.Now()
	c.ModifyUser = uname
}

func (c *S_Catalog) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, c)
}

//查询数据库
func (c *S_Catalog) Fill() error {
	o := orm.NewOrm()
	if c.Id > 0 {
		return o.Read(c, "Id")
	}

	return fmt.Errorf("请确认是否传递了Id", "")

}

//插入
func (c *S_Catalog) Insert() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(c)
}

//修改
func (c *S_Catalog) Update(column ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(c, column...)
}

func init() {
	orm.RegisterModel(new(S_Catalog))
}

func (c *S_Catalog) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		beego.Warn("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

//相关函数

//根据ID得到菜单数据
func GetCatalog(id int) (*S_Catalog, error) {
	c := S_Catalog{Id: id}
	err := c.Fill()
	if err != nil {
		beego.Error("栏目Id不存在")
		return &c, err
	}
	return &c, nil
}

//得到父级InnerCode
func GetCatalogInnerCode(id int) string {
	o := orm.NewOrm()
	var c S_Catalog
	err := o.QueryTable("s_catalog").Filter("Id", id).One(&c, "innercode")
	if err == orm.ErrNoRows {
		// 没有找到记录
		return ""
	}
	return c.InnerCode
}

func GetPreviousId(level int) int {
	o := orm.NewOrm()
	var c S_Catalog
	err := o.Raw("select max(id) id from s_catalog where level = ?", level).QueryRow(&c)
	if err == orm.ErrNoRows {
		// 没有找到记录
		return 0
	}
	return c.Id
}

//得到所有的栏目
func GetCatalogs() ([]*S_Catalog, error) {
	var cs []*S_Catalog
	o := orm.NewOrm()
	_, err := o.QueryTable("s_catalog").OrderBy("-previousid").All(&cs)
	mp := make(map[int]*S_Catalog)
	for _, c := range cs {
		mp[c.PreviousId] = c
	}
	var pid = 0
	for i := 0; i < len(cs); i++ {
		if c, ok := mp[pid]; !ok {
			return cs[:i], err
		} else {
			cs[i] = c
			pid = c.Id
		}
	}
	return cs, err
}

//修改PreviousId，从id1改为id2
func GetCatalogByPrevId(id int) *S_Catalog {
	o := orm.NewOrm()
	var c S_Catalog
	err := o.Raw("select * from s_catalog where previousid = ?", id).QueryRow(&c)
	if err == orm.ErrNoRows {
		return nil
	}
	return &c
}

//得到所有栏目并按级别排序
func GetCatalogsLevel(url string) ([]*S_Catalog, error) {
	cs, err := GetCatalogs()
	if err == nil && cs != nil {
		var cslevel []*S_Catalog = make([]*S_Catalog, 0)
		//用来记录cslevel的当前位置
		idx := -1
		//线型遍历一遍
		for _, c := range cs {
			c.IsLeaf = true //默认全部是叶子节点
			if strings.EqualFold(c.Link, url) {
				c.Checked = true
			}
			if c.Pid == 0 {
				idx++
				cslevel = append(cslevel, c)
				continue
			}

			if c.Pid == cslevel[idx].Id {
				//如果当前元素是cslevel第idx个元素的子集时就放入ChildNode中
				if c.Checked {
					cslevel[idx].Checked = true
					cslevel[idx].IsLeaf = false //如果有子级则不为叶子节点
				}
				cslevel[idx].ChildNode = append(cslevel[idx].ChildNode, c)
			}
			//循环到了这里的元素都是找不到父级的元素这里直接丢弃
		}
		return cslevel, nil
	}
	return cs, err
}

//得到分页的菜单
/**
*   size    每页查询长度
*   index   查询的页码
*   ordercolumn 排序字段
*   orderby     升降序:desc\asc
**/
func GetCatalogsPage(size, index int, ordercolumn, orderby string) (*DataGrid, error) {
	if ordercolumn == "" {
		ordercolumn = "-previousid"
	} else if strings.EqualFold(orderby, "asc") {
		ordercolumn = "previousid"
	}
	var cs []*S_Catalog
	o := orm.NewOrm()
	_, err := o.QueryTable("s_catalog").OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&cs)
	if err == nil {

		mp := make(map[int]*S_Catalog)
		for _, c := range cs {
			mp[c.PreviousId] = c
		}
		var pid = 0
		for i := 0; i < len(cs); i++ {
			if c, ok := mp[pid]; !ok {
				cs = cs[:i]
				break
			} else {
				cs[i] = c
				pid = c.Id
			}
		}

		cnt, err := o.QueryTable("s_catalog").Count()
		pagetotal := cnt / int64(size)
		if cnt%int64(size) > 0 {
			pagetotal++
		}
		var tempcs []*S_Catalog = make([]*S_Catalog, len(cs))
		for i, c := range cs {
			//展开节点
			c.Expanded = true
			//设置默认都是叶子节点
			c.IsLeaf = true
			tempcs[i] = c

			if c.Level == 2 {
				for j := i - 1; j >= 0; j-- {
					if c.Pid == tempcs[j].Id {
						tempcs[j].IsLeaf = false //如果有子级则不为叶子节点
						break
					}
				}
			}
		}

		return GetDataGrid(tempcs, index, int(pagetotal), cnt), err
	}

	return nil, err
}

//返回对应级别的菜单组
func GetCatalogsByLevel(level int) ([]*S_Catalog, error) {
	var cs []*S_Catalog
	o := orm.NewOrm()
	_, err := o.QueryTable("s_catalog").Filter("Level", level).OrderBy("previousid").All(&cs)
	return cs, err
}

type CatalogSelectInit struct {
	S_Catalog
	Select string
}

//返回所有顶级菜单，并指定当前父级和去掉自身
func GetTopCatalogs(pid, self int) ([]CatalogSelectInit, error) {
	cs, err := GetCatalogsByLevel(1)
	if err != nil {
		return nil, err
	}

	cselectinit := make([]CatalogSelectInit, 0)
	for _, c := range cs {
		if c.Id == self {
			continue
		}
		selectInit := new(CatalogSelectInit)
		selectInit.S_Catalog.Id = c.Id
		selectInit.S_Catalog.CatalogName = c.CatalogName
		if c.Id == pid {
			selectInit.Select = "selected"
		}
		cselectinit = append(cselectinit, *selectInit)
	}
	return cselectinit, nil
}
