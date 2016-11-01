package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int    `orm:"pk"`
	UserName string `orm:"null;column(UserName)"`
	Password string `orm:"null;column(Password)"`
	Email    string
	Phone    string
	AddTime  time.Time `orm:"column(AddTime)"`
	AddUser  string    `orm:"column(AddUser)"`
}

func (user *User) SetID(id int) {
	user.Id = id
}

func (user *User) GetID() int {
	return user.Id
}

func (user *User) GetUserName() string {
	return user.UserName
}

func (user *User) GetPassword() string {
	return user.Password
}

func init() {
	orm.RegisterModel(new(User))
}

func (user *User) String() string {
	return fmt.Sprintf("{User:{Id:%d,UserName:'%s',Email:'%s',Phone:'%s',AddTime:'%s',AddUser:'%s'}}", user.Id, user.UserName, user.Email, user.Phone, user.AddTime, user.AddUser)
}

//得到用户信息
func GetUser(username string) (User, error) {
	user := User{UserName: username}
	o := orm.NewOrm()
	err := o.Read(&user, "UserName")
	if err != nil {
		return user, fmt.Errorf("用户[%s]不存在", username)
	}
	return user, nil
}
