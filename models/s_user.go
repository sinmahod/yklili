package models

import (
	"beegostudy/util/modelutil"
	"beegostudy/util/pwdutil"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
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
type S_User struct {
	Id         int       `orm:"pk;column(id)"`
	UserName   string    `orm:"column(username);index;unique;size(64)"`
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
func (u *S_User) TableName() string {
	return "s_user"
}

func (user *S_User) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	userid, err := strconv.Atoi(tmpId)
	if err == nil {
		user.Id = userid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (user *S_User) GetId() int {
	return user.Id
}

func (user *S_User) GetUserName() string {
	return user.UserName
}

func (user *S_User) GetPassword() string {
	return user.Password
}

func (user *S_User) SetPassword(pwd string) {
	user.Password = pwdutil.GeneratePWD(pwd)
}

func (user *S_User) SetAddTime(t time.Time) {
	user.AddTime = t
}

func (user *S_User) SetAddUser(uname string) {
	user.AddTime = time.Now()
	user.AddUser = uname
}

func (user *S_User) SetCurrentTime() {
	user.AddTime = time.Now()
}

func (user *S_User) SetModifyUser(uname string) {
	user.ModifyTime = time.Now()
	user.ModifyUser = uname
}

func (user *S_User) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, user)
}

func init() {
	orm.RegisterModel(new(S_User))
}

func (user *S_User) Fill() error {
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
	return fmt.Errorf("请确认是否传递了Id或UserName或Email")

}

func (user *S_User) String() string {
	data, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

//根据Id得到用户信息
func GetUserById(id int) (*S_User, error) {
	user := S_User{Id: id}
	err := user.Fill()
	if err != nil {
		return &user, fmt.Errorf("用户Id[%s]不存在", id)
	}
	return &user, nil
}

//根据用户名得到用户信息
func GetUser(username string) (*S_User, error) {
	user := S_User{UserName: username}
	err := user.Fill()
	if err != nil {
		return &user, fmt.Errorf("用户[%s]不存在", username)
	}
	return &user, nil
}

//用户是否存在
func UserExists(username string) int64 {
	o := orm.NewOrm()
	qt := o.QueryTable("s_user")
	qt = qt.Filter("UserName", username)
	cnt, _ := qt.Count()
	return cnt
}

//邮箱是否存在
func EmailExists(email string) int64 {
	o := orm.NewOrm()
	qt := o.QueryTable("s_user")
	qt = qt.Filter("Email", email)
	cnt, _ := qt.Count()
	return cnt
}

//根据邮箱得到用户信息
func GetUserByEmail(email string) (*S_User, error) {
	user := S_User{Email: email}
	err := user.Fill()
	if err != nil {
		return &user, fmt.Errorf("邮箱[%s]未注册", email)
	}
	return &user, nil
}

//得到分页的菜单
/**
*	size	每页查询长度
*	index	查询的页码
*	ordercolumn	排序字段
*	orderby		升降序:desc\asc
**/
func GetUsersPage(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "addtime"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var users []*S_User
	o := orm.NewOrm()
	qt := o.QueryTable("s_user")
	if data["UserName"] != nil {
		qt = qt.Filter("UserName__icontains", data["UserName"])
	}
	_, err := qt.OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&users, "Id", "UserName", "RealName", "Email", "Phone", "AddTime", "AddUser")

	if err == nil {
		cnt, err := qt.Count()

		pagetotal := cnt / int64(size)

		if cnt%int64(size) > 0 {
			pagetotal++
		}

		return GetDataGrid(users, index, int(pagetotal), cnt), err
	}

	return nil, err
}

//插入用户
func InsertUser(username, password, email, phone string) (int64, error) {
	o := orm.NewOrm()
	user := new(S_User)
	user.Id = GetMaxId("S_UserID")
	user.UserName = username
	user.Password = password
	user.Email = email
	user.Phone = phone
	user.AddTime = time.Now()
	user.AddUser = username
	return o.Insert(user)
}
