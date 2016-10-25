package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id    int
	Name  string
	Email string
	Phone string
}

func init() {
	orm.RegisterModel(new(User))
}

func (user *User) String() string {
	return fmt.Sprintf("{User:{Id:%d,Name:'%s',Email:'%s',Phone:'%s'}}", user.Id, user.Name, user.Email, user.Phone)
}
