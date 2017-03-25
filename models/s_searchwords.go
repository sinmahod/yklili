package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	transaction "yklili/models/orm"
	"yklili/service/progress"
	"yklili/util/fileutil"
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
type S_SearchWords struct {
	Id         int       `orm:"pk;column(id)"`
	SWords     string    `orm:"unique;column(swords);size(64)"`      //根词
	Synonym    string    `orm:"null;column(synonym);index;size(64)"` //同义词
	AddTime    time.Time `orm:"auto_now_add;type(datetime);column(addtime)"`
	AddUser    string    `orm:"column(adduser);size(64)"`
	ModifyTime time.Time `orm:"null;type(datetime);column(modifytime)"`
	ModifyUser string    `orm:"null;column(modifyuser);size(64)"`
}

//自定义表名
func (u *S_SearchWords) TableName() string {
	return "s_searchwords"
}

func init() {
	orm.RegisterModel(new(S_SearchWords))
}

func (s *S_SearchWords) SetId(id interface{}) error {
	tmpId := fmt.Sprintf("%v", id)
	sid, err := strconv.Atoi(tmpId)
	if err == nil {
		s.Id = sid
	} else {
		beego.Error("Id字段必须为正整数型【%v】\n", id)
	}
	return err
}

func (s *S_SearchWords) GetId() int {
	return s.Id
}

func (s *S_SearchWords) SetSWords(words string) {
	s.SWords = words
}

func (s *S_SearchWords) GetSWords() string {
	return s.SWords
}

func (s *S_SearchWords) SetSynonym(syno string) {
	s.Synonym = syno
}

func (s *S_SearchWords) GetSynonym() string {
	return s.Synonym
}

func (s *S_SearchWords) SetAddUser(uname string) {
	s.AddTime = time.Now()
	s.AddUser = uname
}

func (s *S_SearchWords) SetModifyUser(uname string) {
	s.ModifyTime = time.Now()
	s.ModifyUser = uname
}

func (s *S_SearchWords) SetValue(data map[string]interface{}) error {
	return modelutil.FillStruct(data, s)
}

func (s *S_SearchWords) Fill() error {
	o := orm.NewOrm()
	if s.Id > 0 {
		return o.Read(s, "Id")
	}
	if s.SWords != "" {
		return o.Read(s, "SWords")
	}
	return fmt.Errorf("请确认是否传递了Id或SWords")

}

func (c *S_SearchWords) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Printf("JSON marshaling failed: %s", err)
	}
	return fmt.Sprintf("%s\n", data)
}

func GetWordsById(id int) (*S_SearchWords, error) {
	words := S_SearchWords{Id: id}
	err := words.Fill()
	if err != nil {
		return &words, fmt.Errorf("根词Id[%s]不存在", id)
	}
	return &words, nil
}

func GetWords(swords string) (*S_SearchWords, error) {
	words := S_SearchWords{SWords: swords}
	err := words.Fill()
	if err != nil {
		return &words, fmt.Errorf("根词[%s]不存在", swords)
	}
	return &words, nil
}

//得到分页
/**
*   size    每页查询长度
*   index   查询的页码
*   ordercolumn 排序字段
*   orderby     升降序:desc\asc
**/
func GetWordsPage(size, index int, ordercolumn, orderby string, data map[string]interface{}) (*DataGrid, error) {

	if ordercolumn == "" {
		ordercolumn = "-addtime"
	} else if strings.EqualFold(orderby, "desc") {
		ordercolumn = "-" + ordercolumn
	}

	var cs []*S_SearchWords
	o := orm.NewOrm()
	qt := o.QueryTable("s_searchwords")
	if data["SWords"] != nil {
		qt = qt.Filter("swords__icontains", data["SWords"])
	}
	if data["Synonym"] != nil {
		qt = qt.Filter("synonym__icontains", data["Synonym"])
	}
	_, err := qt.OrderBy(ordercolumn).Limit(size, (index-1)*size).All(&cs)

	if err == nil {
		cnt, err := qt.Count()

		pagetotal := cnt / int64(size)

		if cnt%int64(size) > 0 {
			pagetotal++
		}

		return GetDataGrid(cs, index, int(pagetotal), cnt), err
	}

	return nil, err
}

//重新导入词典(会删除原词典包含同义词关联关系，慎用)
func ImportWords(wordfile, uname string, prog *progress.ProgressTask) error {
	//清空原词典
	o := orm.NewOrm()
	_, err := o.Raw("delete from s_searchwords").Exec()
	if err != nil {
		return err
	}

	tmpMap := make(map[string]interface{})

	tran := new(transaction.Transaction)

	str := fileutil.FileToString(wordfile)

	words := strings.Split(str, "\n")

	total := len(words)

	size := 500

	page := total / size

	if total%size > 0 {
		page++
	}

	for p := 0; p < page; p++ {
		for i := p * size; i < (p+1)*size && i < total; i++ {
			ws := strings.Split(words[i], " ")[0]
			if strings.Trim(ws, " ") == "" {
				continue
			}
			if _, ok := tmpMap[ws]; !ok {
				tmpMap[ws] = nil
				ssw := new(S_SearchWords)
				ssw.SetId(GetMaxId("S_SearchWordsID"))
				ssw.SetSWords(ws)
				ssw.SetAddUser(uname)
				tran.Add(ssw, transaction.INSERT)
			}
		}
		if err = tran.Commit(); err == nil {
			tran.Clear()
		} else {
			fmt.Println("ERROR:", err)
			o.Raw("delete from s_searchwords").Exec()
			return err
		}
		f := (p + 1) * 100 / page
		prog.SetPerc(int(f))
		prog.SetMsg("任务已执行到了%d%s", int(f), "%")
	}

	return nil
}

func WordsExists(swords string) int64 {
	o := orm.NewOrm()
	qt := o.QueryTable("s_searchwords")
	qt = qt.Filter("SWords", swords)
	cnt, _ := qt.Count()
	return cnt
}
