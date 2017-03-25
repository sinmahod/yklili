package models

import (
	"encoding/json"
	"fmt"
	"github.com/sinmahod/yklili/util/modelutil"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
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
type S_Config struct {
	ConfigKey   string    `orm:"pk;column(configkey);size(64)"`
	ConfigValue string    `orm:"column(configvalue);size(256)"`
	Memo        string    `orm:"column(memo);size(512)"`
	AddTime     time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser     string    `orm:"column(adduser);size(64)"`
	ModifyTime  time.Time `orm:"null;type(datetime);column(modifytime)"`
	ModifyUser  string    `orm:"null;column(modifyuser);size(64)"`
}

//自定义表名
func (c *S_Config) TableName() string {
	return "s_config"
}

func (c *S_Config) GetK() string {
	return c.ConfigKey
}

func (c *S_Config) SetK(k string) {
	c.ConfigKey = k
}

func (c *S_Config) GetV() string {
	return c.ConfigValue
}

func (c *S_Config) SetV(v string) {
	c.ConfigValue = v
}

func (c *S_Config) SetMemo(m string) {
	c.Memo = m
}

func (c *S_Config) SetAddUser(uname string) {
	c.AddTime = time.Now()
	c.AddUser = uname
}

func (c *S_Config) SetModifyUser(uname string) {
	c.ModifyTime = time.Now()
	c.ModifyUser = uname
}
func init() {
	orm.RegisterModel(new(S_Config))
}

func (c *S_Config) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, c)
}

func (c *S_Config) String() string {
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
func GetConfigsPage(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "configkey"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var cs []*S_Config
	o := orm.NewOrm()
	qt := o.QueryTable("s_config")
	if data["ConfigKey"] != nil {
		qt = qt.Filter("configkey__icontains", data["ConfigKey"])
	}
	_, err := qt.OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&cs, "ConfigKey", "ConfigValue")

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

//得到所有配置
func GetConfigs() []*S_Config {
	var cs []*S_Config
	o := orm.NewOrm()
	o.QueryTable("s_config").All(&cs, "ConfigKey", "ConfigValue")
	return cs
}
