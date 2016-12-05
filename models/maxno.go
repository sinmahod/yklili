package models

import (
	"beegostudy/util"
	"strconv"
	"strings"
	"sync"

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
type MaxNo struct {
	NoName   string `orm:"pk;column(noname);size(32)"` //此表为双主键，目前beego v1.5.1不支持
	NoType   string `orm:"column(notype);size(64)"`
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

var mutex sync.Mutex

func GetMaxNo(noname, notype string, length int) string {
	//添加锁
	mutex.Lock()
	//解锁
	defer mutex.Unlock()
	if notype == "" {
		notype = "SN"
	}
	o := orm.NewOrm()
	var maxno string
	if v, _ := getNoValue(noname, notype); v == 0 {
		maxno = util.LeftPad("1", '0', length)
		if !strings.EqualFold(notype, "SN") {
			maxno = notype + maxno
		}
		no := MaxNo{NoName: noname, NoType: notype, NoValue: 1, NoLength: length}
		o.Insert(&no)
	} else {
		v++
		maxno = util.LeftPad(strconv.Itoa(v), '0', length)
		if !strings.EqualFold(notype, "SN") {
			maxno = notype + maxno
		}
		no := MaxNo{NoName: noname, NoType: notype, NoValue: v, NoLength: length}
		o.Update(&no)
	}
	return maxno
}

func getNoValue(noname, notype string) (int, error) {
	o := orm.NewOrm()
	var maxno MaxNo
	err := o.QueryTable("maxno").Filter("noname", noname).Filter("notype", notype).One(&maxno, "novalue")
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		beego.Error("错误：查询一条数据时返回了多条")
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		return 0, nil
	}

	return maxno.NoValue, err
}
