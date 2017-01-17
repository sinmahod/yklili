package models

import (
	"beegostudy/util/modelutil"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

/**
*   pk      主键
*   auto        自增值（限数值）
*   column(N)   指定字段名N
*   null        可以为非空（默认非空）
*   index       单个字段索引
*   unique      唯一键
*   auto_now_add    第一次插入数据时自动添加当前时间
*   auto_now    每一次保存时自动更新当前时间
*   type(T)     对应数据库的指定类型
*   size(S)     类型长度S
*   default(D)  默认值D（需要对应类型）
**/
type S_Site struct {
	Id         int       `orm:"pk;column(id);"`
	SiteName   string    `orm:"column(sitename);size(128)"`
	SiteHost   string    `orm:"column(sitehost);size(128)"`
	SiteDesc   string    `orm:"null;column(sitedesc);size(512)"`
	SiteLogo   string    `orm:"null;column(sitelogo);size(256)"`
	SiteBanner string    `orm:"null;column(sitebanner);size(256)"`
	Memo       string    `orm:"null;column(memo);size(512)"`
	AddTime    time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser    string    `orm:"column(adduser);size(64)"`
	ModifyTime time.Time `orm:"null;type(datetime);column(modifytime)"`
	ModifyUser string    `orm:"null;column(modifyuser);size(64)"`
}

//自定义表名
func (c *S_Site) TableName() string {
	return "s_site"
}

func (c *S_Site) GetId() int {
	return c.Id
}

func (c *S_Site) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	siteid, err := strconv.Atoi(tmpId)
	if err == nil {
		c.Id = siteid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (c *S_Site) GetName() string {
	return c.SiteName
}

func (c *S_Site) SetName(name string) {
	c.SiteName = name
}

func (c *S_Site) GetHost() string {
	return c.SiteHost
}

func (c *S_Site) SetHost(host string) {
	c.SiteHost = host
}

func (c *S_Site) GetDesc() string {
	return c.SiteDesc
}

func (c *S_Site) SetDesc(desc string) {
	c.SiteDesc = desc
}

func (c *S_Site) GetBanner() string {
	return c.SiteBanner
}

func (c *S_Site) SetBanner(banner string) {
	c.SiteBanner = banner
}

func (c *S_Site) SetMemo(m string) {
	c.Memo = m
}

func (c *S_Site) SetAddUser(uname string) {
	c.AddTime = time.Now()
	c.AddUser = uname
}

func (c *S_Site) SetModifyUser(uname string) {
	c.ModifyTime = time.Now()
	c.ModifyUser = uname
}
func init() {
	orm.RegisterModel(new(S_Site))
}

func (c *S_Site) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, c)
}

func (site *S_Site) Fill() error {
	o := orm.NewOrm()
	if site.Id > 0 {
		return o.Read(site, "Id")
	}
	return fmt.Errorf("请确认是否传递了Id")

}

func (c *S_Site) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

//得到分页
/**
*   size    每页查询长度
*   index   查询的页码
*   ordercolumn 排序字段
*   orderby     升降序:desc\asc
**/
func GetSitePage(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "configkey"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var cs []*S_Site
	o := orm.NewOrm()
	qt := o.QueryTable("s_site")
	if data["SiteName"] != nil {
		qt = qt.Filter("sitename__icontains", data["SiteName"])
	}
	_, err := qt.OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&cs)

	if err == nil {
		cnt, err := qt.Count()

		pagetotal := cnt / int64(size)

		if cnt%int64(size) > 0 {
			pagetotal++
		}

		return GetDataGrid(cs, index, int(pagetotal), cnt), err
	}

	return nil, err
}

func GetSite(id int) (*S_Site, error) {
	site := S_Site{Id: id}
	err := site.Fill()
	if err != nil {
		return &site, fmt.Errorf("站点Id[%s]不存在", id)
	}
	return &site, nil
}

//得到所有配置
func GetSites() []*S_Site {
	var cs []*S_Site
	o := orm.NewOrm()
	o.QueryTable("s_site").All(&cs)
	return cs
}
