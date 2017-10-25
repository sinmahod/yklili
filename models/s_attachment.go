package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sinmahod/yklili/util/modelutil"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

/**
*   pk      主键
*   auto        自增值（限数值）
*   column(N)   指定字段名N
*   null        非空
*   index       单个字段索引
*   unique      唯一键
*   auto_now_add    第一次插入数据时自动添加当前时间
*   auto_now    每一次保存时自动更新当前时间
*   type(T)     对应数据库的指定类型
*   size(S)     类型长度S
*   default(D)  默认值D（需要对应类型）
**/
type S_Attachment struct {
	Id          int       `orm:"pk;column(id)"`
	FileName    string    `orm:"column(filaname);size(128)"`
	FileNewName string    `orm:"column(filenewname);size(128)"`
	FilePath    string    `orm:"column(filepath);size(256)"`
	FileType    string    `orm:"column(filetype);size(64)"`
	FileSize    int64     `orm:"column(filesize)"`
	AddTime     time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser     string    `orm:"column(adduser);size(64)"`
}

//自定义表名
func (m *S_Attachment) TableName() string {
	return "s_attachment"
}

func (m *S_Attachment) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	mid, err := strconv.Atoi(tmpId)
	if err == nil {
		m.Id = mid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (m *S_Attachment) SetAddTime(t time.Time) {
	m.AddTime = t
}

func (m *S_Attachment) SetCurrentTime() {
	m.AddTime = time.Now()
}

func (m *S_Attachment) SetAddUser(uname string) {
	m.AddUser = uname
}

func (m *S_Attachment) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, m)
}

//查询数据库
func (m *S_Attachment) Fill() error {
	o := orm.NewOrm()
	if m.Id > 0 {
		return o.Read(m, "Id")
	}

	return fmt.Errorf("请确认是否传递了Id", "")

}

//插入
func (m *S_Attachment) Insert() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(m)
}

//修改
func (m *S_Attachment) Update(column ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(m, column...)
}

func init() {
	orm.RegisterModel(new(S_Attachment))
}

func (m *S_Attachment) String() string {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		beego.Warn("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

//相关函数
//根据ID查找附件
func GetArrachment(id int) (*S_Attachment, error) {
	m := S_Attachment{Id: id}
	err := m.Fill()
	if err != nil {
		beego.Error("附件Id不存在")
		return &m, err
	}
	return &m, nil
}

//新创建一个附件
func AddAttachment(filename, filenewname, filepath, filetype string, filesize int64, adduser string) *S_Attachment {
	m := new(S_Attachment)
	m.Id = GetMaxId("S_AttachmentID")
	m.FileName = filename
	m.FileNewName = filenewname
	m.FilePath = filepath
	m.FileSize = filesize
	m.FileType = filetype
	m.SetCurrentTime()
	m.AddUser = adduser
	m.Insert()
	return m
}

//得到分页的菜单
/**
*   size    每页查询长度
*   index   查询的页码
*   ordercolumn 排序字段
*   orderby     升降序:desc\asc
**/
func GetAttchmentsPage(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "-addtime"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var atta []*S_Attachment
	o := orm.NewOrm()
	qt := o.QueryTable("s_attachment")
	if data["FileName"] != nil {
		qt = qt.Filter("FileName__icontains", data["FileName"])
	}
	_, err := qt.OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&atta)

	if err == nil {
		cnt, err := qt.Count()

		pagetotal := cnt / int64(size)

		if cnt%int64(size) > 0 {
			pagetotal++
		}

		return GetDataGrid(atta, index, int(pagetotal), cnt), err
	}

	return nil, err
}
