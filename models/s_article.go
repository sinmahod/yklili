package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"yklili/service/bleve"
	"yklili/service/progress"
	"yklili/util/modelutil"
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
type S_Article struct {
	Id         int       `orm:"pk;column(id)"`
	PackageId  int       `orm:"column(packageid)"`
	Title      string    `orm:"column(title);size(128)"`
	Content    string    `orm:"null;column(content);type(text)"`
	Logo       int       `orm:"column(logo)"`
	Status     int       `orm:"column(status)"`
	AddTime    time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser    string    `orm:"column(adduser);size(64)"`
	ModifyTime time.Time `orm:"null;type(datetime);column(modifytime)"`
	ModifyUser string    `orm:"null;column(modifyuser);size(64)"`
}

const (
	DRAFT = iota
	PUBLISH
	DELETE
)

//自定义表名
func (a *S_Article) TableName() string {
	return "s_article"
}

func (a *S_Article) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	aid, err := strconv.Atoi(tmpId)
	if err == nil {
		a.Id = aid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (a *S_Article) SetPackageId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	aid, err := strconv.Atoi(tmpId)
	if err == nil {
		a.PackageId = aid
	} else {
		beego.Error("PackageId字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (a *S_Article) SetStatus(state int) {
	a.Status = state
}

func (a *S_Article) GetId() int {
	return a.Id
}

func (a *S_Article) GetPackageId() int {
	return a.PackageId
}

func (a *S_Article) SetAddTime(t time.Time) {
	a.AddTime = t
}

func (a *S_Article) SetAddUser(uname string) {
	a.AddTime = time.Now()
	a.AddUser = uname
}

func (a *S_Article) SetCurrentTime() {
	a.AddTime = time.Now()
}

func (a *S_Article) SetModifyUser(uname string) {
	a.ModifyTime = time.Now()
	a.ModifyUser = uname
}

func (a *S_Article) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, a)
}

func init() {
	orm.RegisterModel(new(S_Article))
}

func (a *S_Article) Fill() error {
	o := orm.NewOrm()
	if a.Id > 0 {
		return o.Read(a, "Id")
	}
	return fmt.Errorf("请确认是否传递了Id", "")

}

func (a *S_Article) String() string {
	data, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

func GetArticle(id int) (*S_Article, error) {
	a := S_Article{Id: id}
	err := a.Fill()
	if err != nil {
		return &a, fmt.Errorf("文章Id[%s]不存在", id)
	}
	return &a, nil
}

func GetArticleByStatus(id, status int) (*S_Article, error) {
	o := orm.NewOrm()
	a := S_Article{Id: id, Status: status}
	err := o.Read(&a, "Id", "Status")
	if err != nil {
		return &a, fmt.Errorf("文章Id[%s]不存在", id)
	}
	return &a, nil
}

//得到分页的菜单
/**
*   size    每页查询长度
*   index   查询的页码
*   ordercolumn 排序字段
*   orderby     升降序:desc\asc
**/
func GetArticlesPage(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "addtime"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var as []*S_Article
	o := orm.NewOrm()
	qt := o.QueryTable("s_article")

	if data["User"] != nil {
		user := data["User"].(*S_User)
		qt = qt.Filter("AddUser", user.GetUserName())
	}

	if data["Status"] != nil {
		qt = qt.Filter("Status", data["Status"])
	}

	if data["PackageId"] != nil {
		tmpId := fmt.Sprintf("%v", data["PackageId"])
		if pid, err := strconv.Atoi(tmpId); err == nil && pid != 0 {
			qt = qt.Filter("PackageId", pid)
		}
	}

	qt = qt.Exclude("Status", DELETE)
	_, err := qt.OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&as)

	if err == nil {
		cnt, err := qt.Count()

		pagetotal := cnt / int64(size)

		if cnt%int64(size) > 0 {
			pagetotal++
		}

		for _, a := range as {
			if a.Status == DRAFT {
				a.Title = "* " + a.Title
			}
		}

		return GetDataGrid(as, index, int(pagetotal), cnt), err
	}

	return nil, err
}

//全库重建索引
func RebuildIndex(prog *progress.ProgressTask) error {
	err := bleve.ClearIndex()
	if err != nil {
		return err
	}

	var size int64 = 100

	var as []S_Article
	o := orm.NewOrm()
	qt := o.QueryTable("s_article")
	qt = qt.Filter("Status", 1)

	cnt, err := qt.Count()

	page := cnt / size

	if cnt%size > 0 {
		page++
	}

	for i := int64(0); i < page; i++ {
		_, err = qt.Limit(size, i*size).All(&as)
		if err != nil {
			return err
		}

		for _, a := range as {
			err = bleve.AddIndex(strconv.Itoa(a.Id), a)
			if err != nil {
				return err
			}
		}

		f := (i + 1) * 100 / page

		prog.SetPerc(int(f))
		prog.SetMsg("任务已执行到了%d%s", int(f), "%")

	}

	return nil
}
