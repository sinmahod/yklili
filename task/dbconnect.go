package task

import (
	"yklili/service/cron"

	"github.com/astaxie/beego/orm"
)

func init() {
	cron.RegisterTask(&DBConnectTask{})
}

type DBConnectTask struct{}

func (t *DBConnectTask) GetId() string {
	return "DBConnect"
}

func (t *DBConnectTask) GetSpec() string {
	return "0  *  *  *  *  *"
}

func (t *DBConnectTask) GetDesc() string {
	return "保持数据库连接（每1分钟执行一次）"
}

func (t *DBConnectTask) Execute() {
	o := orm.NewOrm()
	o.Raw("SELECT COUNT(1) FROM s_config").QueryRow()
}
