package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int    `orm:"pk"`
	UserName string `orm:"null;column(UserName)"`
	Password string `orm:"null"`
	Email    string
	Phone    string
	AddTime  time.Time `orm:"column(AddTime)"`
	AddUser  string    `orm:"column(AddUser)"`
}

func (user *User) SetID(id int) {
	user.Id = id
}

func init() {
	orm.RegisterModel(new(User))
}

func (user *User) String() string {
	return fmt.Sprintf("{User:{Id:%d,UserName:'%s',Email:'%s',Phone:'%s',AddTime:'%s',AddUser:'%s'}}", user.Id, user.UserName, user.Email, user.Phone, user.AddTime, user.AddUser)
}
