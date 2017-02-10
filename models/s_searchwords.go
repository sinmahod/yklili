package models

import (
	"github.com/astaxie/beego/orm"
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
type S_SearchWords struct {
	SWords  string `orm:"pk;column(swords);size(64)"`
	Synonym string `orm:"null;column(synonym);index;size(64)"`

	RealName   string    `orm:"column(realname);size(128)"`
	Password   string    `orm:"column(password);size(64)"`
	Email      string    `orm:"column(email);index;unique;size(64)"`
	Phone      string    `orm:"null;column(phone);size(32)"`
	AddTime    time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser    string    `orm:"column(adduser);size(64)"`
	ModifyTime time.Time `orm:"null;type(datetime);column(modifytime)"`
	ModifyUser string    `orm:"null;column(modifyuser);size(64)"`
}

//自定义表名
func (u *S_SearchWords) TableName() string {
	return "s_searchwords"
}

func init() {
	orm.RegisterModel(new(S_SearchWords))
}
