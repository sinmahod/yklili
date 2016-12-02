package models

import "github.com/astaxie/beego/orm"

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
type MaxNo struct {
	NoName   string `orm:"pk;column(noname);size(32)"`
	NoType   string `orm:"pk;column(notype);size(64)"`
	NoValue  int    `orm:"column(novalue)"`
	NoLength int    `orm:"column(nolength)"`
}

//自定义表名
func (m *MaxNo) TableName() string {
	return "maxno"
}

func init() {
	orm.RegisterModel(new(MaxNo))
}
