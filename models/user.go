package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

/**
*	pk		主键
*	auto 		自增值（限数值）
*	column(N)	指定字段名N
*	null		可以为非空（默认非空）
*	index 		单个字段索引
* 	unique 		唯一键
* 	auto_now_add 	第一次插入数据时自动添加当前时间
* 	auto_now 	每一次保存时自动更新当前时间
* 	type(T)		对应数据库的指定类型
*	size(S)		类型长度S
*	default(D)	默认值D（需要对应类型）
**/
type User struct {
	Id       int       `orm:"pk;auto;column(id)"`
	UserName string    `orm:"column(username);index;unique;size(64)"`
	Password string    `orm:"column(password);size(64)"`
	Email    string    `orm:"column(email);index;unique;size(64)"`
	Phone    string    `orm:"null;column(phone);size(32)"`
	AddTime  time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser  string    `orm:"column(adduser);size(64)"`
}

//自定义表名
func (u *User) TableName() string {
	return "user"
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

func (user *User) Fill() error {
	o := orm.NewOrm()
	if user.Id > 0 {
		return o.Read(user, "Id")
	}
	if user.UserName != "" {
		return o.Read(user, "UserName")
	}
	if user.Email != "" {
		return o.Read(user, "Email")
	}
	return fmt.Errorf("请确认是否传递了Id或UserName或Email", "")

}

func (user *User) String() string {
	return fmt.Sprintf("{User:{Id:%d,UserName:'%s',Email:'%s',Phone:'%s',AddTime:'%s',AddUser:'%s'}}", user.Id, user.UserName, user.Email, user.Phone, user.AddTime, user.AddUser)
}

//根据用户名得到用户信息
func GetUser(username string) (*User, error) {
	user := User{UserName: username}
	err := user.Fill()
	if err != nil {
		return &user, fmt.Errorf("用户[%s]不存在", username)
	}
	return &user, nil
}

//根据邮箱得到用户信息
func GetUserByEmail(email string) (*User, error) {
	user := User{Email: email}
	err := user.Fill()
	if err != nil {
		return &user, fmt.Errorf("邮箱[%s]未注册", email)
	}
	return &user, nil
}

//插入用户
func InsertUser(username, password, email, phone string) (int64, error) {
	o := orm.NewOrm()
	user := new(User)
	user.UserName = username
	user.Password = password
	user.Email = email
	user.Phone = phone
	user.AddTime = time.Now()
	user.AddUser = username
	return o.Insert(user)
}
