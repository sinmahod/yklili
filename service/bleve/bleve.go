package bleve

import (
	"errors"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	_ "github.com/yanyiwu/gojieba/bleve"
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	INDEX_DIR  = "bleveindex"
	FENCI_TYPE = "gojieba"
)

var (
	DICT_DIR        string
	DICT_PATH       string
	HMM_PATH        string
	USER_DICT_PATH  string
	IDF_PATH        string
	STOP_WORDS_PATH string
)

var (
	index        bleve.Index
	indexMapping *mapping.IndexMappingImpl
)

//初始化
func init() {

	DICT_DIR = path.Join(path.Dir(GetCurrentFilePath()), "dict")
	DICT_PATH = path.Join(DICT_DIR, "jieba.dict.utf8")
	HMM_PATH = path.Join(DICT_DIR, "hmm_model.utf8")
	USER_DICT_PATH = path.Join(DICT_DIR, "user.dict.utf8")
	IDF_PATH = path.Join(DICT_DIR, "idf.utf8")
	STOP_WORDS_PATH = path.Join(DICT_DIR, "stop_words.utf8")

	indexMapping = bleve.NewIndexMapping()

	err := indexMapping.AddCustomTokenizer(FENCI_TYPE,
		map[string]interface{}{
			"dictpath":     DICT_PATH,
			"hmmpath":      HMM_PATH,
			"userdictpath": USER_DICT_PATH,
			"idf":          IDF_PATH,
			"stop_words":   STOP_WORDS_PATH,
			"type":         FENCI_TYPE,
		},
	)
	if err != nil {
		panic(err)
	}
	err = indexMapping.AddCustomAnalyzer(FENCI_TYPE,
		map[string]interface{}{
			"type":      FENCI_TYPE,
			"tokenizer": FENCI_TYPE,
		},
	)
	if err != nil {
		panic(err)
	}

	indexMapping.DefaultAnalyzer = FENCI_TYPE

	index, err = bleve.Open(INDEX_DIR)
	if err != nil {
		if err.Error() == "cannot open index, path does not exist" {
			index, err = bleve.New(INDEX_DIR, indexMapping)
			if err != nil {
				panic(err)
			}
		}
	}
}

func GetCurrentFilePath() string {
	_, filePath, _, _ := runtime.Caller(1)
	idx := strings.LastIndex(filePath, "service")
	return filePath[:idx]
}

//清除索引文件，用于索引重建（慎用）
func ClearIndex() error {
	os.RemoveAll(INDEX_DIR)

	var err error
	index, err = bleve.New(INDEX_DIR, indexMapping)
	if err != nil {
		panic(err)
	}
	return err
}

//得到索引DOC总数量
func GetDocCount() (uint64, error) {
	return index.DocCount()
}

//添加／更新索引
func AddIndex(docId string, data interface{}) error {
	return index.Index(docId, data)
}

//删除索引
func DelIndex(docId string) error {
	return index.Delete(docId)
}

//查询
type Bleve struct {
	querys []string
}

//And形式添加检索条件
func And(query string, column ...string) *Bleve {
	q := make([]string, 0)
	if len(column) > 0 {
		for _, c := range column {
			q = append(q, "+"+c+":"+query)
		}
	} else {
		q = append(q, "+"+query)
	}

	b := new(Bleve)
	b.querys = q
	return b
}

//Not形式添加检索条件
func Not(query string, column ...string) *Bleve {
	q := make([]string, 0)
	if len(column) > 0 {
		for _, c := range column {
			q = append(q, "-"+c+":"+query)
		}
	} else {
		q = append(q, "-"+query)
	}

	b := new(Bleve)
	b.querys = q
	return b
}

//Or形式添加检索条件
func Or(query string, column ...string) *Bleve {
	q := make([]string, 0)
	if len(column) > 0 {
		for _, c := range column {
			q = append(q, c+":"+query)
		}
	} else {
		q = append(q, query)
	}

	b := new(Bleve)
	b.querys = q
	return b
}

func (b *Bleve) And(query string, column ...string) *Bleve {
	if len(column) > 0 {
		for _, c := range column {
			b.querys = append(b.querys, "+"+c+":"+query)
		}
	} else {
		b.querys = append(b.querys, "+"+query)
	}
	return b
}

func (b *Bleve) Not(query string, column ...string) *Bleve {
	if len(column) > 0 {
		for _, c := range column {
			b.querys = append(b.querys, "-"+c+":"+query)
		}
	} else {
		b.querys = append(b.querys, "-"+query)
	}
	return b
}

func (b *Bleve) Or(query string, column ...string) *Bleve {
	if len(column) > 0 {
		for _, c := range column {
			b.querys = append(b.querys, c+":"+query)
		}
	} else {
		b.querys = append(b.querys, query)
	}
	return b
}

type Data struct {
	Results           []*Result
	SearchTime        float64  `json:"searchtime"`
	SearchQuerys      []string `json:"searchquerys"`
	SearchResultTotal uint64   `json:"searchresulttotal"`
}

type Result struct {
	Id      string                 `json:"id"`
	Score   float64                `json:"score"`
	Content map[string]interface{} `json:"content"`
}

//开始检索，返回字符串
func (b *Bleve) SearchToData(size, page int) (*Data, error) {
	page--

	from := size * page

	query := strings.Join(b.querys, " ")

	que := bleve.NewQueryStringQuery(query)

	req := bleve.NewSearchRequestOptions(que, size, from, false)

	req.Highlight = bleve.NewHighlight()

	s := time.Now()

	res, err := index.Search(req)
	if err != nil {
		panic(err)
		return nil, err
	}

	e := time.Since(s).Nanoseconds()

	data := new(Data)

	for _, item := range res.Hits {
		r := new(Result)
		r.Id = item.ID
		r.Score = item.Score

		doc, _ := index.Document(item.ID)

		mp := make(map[string]interface{})

		for _, field := range doc.Fields {
			mp[field.Name()] = string(field.Value())
		}

		for fragmentField, fragments := range item.Fragments {
			for _, fragment := range fragments {
				mp["Search"+fragmentField] = fragment
			}
		}
		r.Content = mp
		data.Results = append(data.Results, r)
	}

	data.SearchTime = float64(e) / 1000000
	data.SearchQuerys = b.querys
	data.SearchResultTotal = res.Total

	log.Printf("SearchTime: %f \n", float64(e)/1000000)
	log.Printf("SearchQuerys: %s \n", b.querys)
	log.Printf("SearchResultTotal: %d \n", res.Total)
	log.Printf("SearchSize: %d ,SearchForm: %d \n", size, from)

	return data, err
}

//开始检索，填充结构数组，返回总条数
func (b *Bleve) Search(size, page int, obj interface{}) (uint64, error) {
	page--

	from := size * page

	query := strings.Join(b.querys, " ")

	que := bleve.NewQueryStringQuery(query)

	req := bleve.NewSearchRequestOptions(que, size, from, false)

	req.Highlight = bleve.NewHighlight()

	s := time.Now()

	res, err := index.Search(req)
	if err != nil {
		panic(err)
		return 0, err
	}

	e := time.Since(s).Nanoseconds()

	log.Printf("SearchTime: %f \n", float64(e)/1000000)
	log.Printf("SearchQuerys: %s \n", b.querys)
	log.Printf("SearchResultTotal: %d \n", res.Total)
	log.Printf("SearchSize: %d ,SearchForm: %d \n", size, from)

	return res.Total, prettify(res, obj)
}

func prettify(res *bleve.SearchResult, obj interface{}) error {

	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("传递的对象必须为指针")
	}

	ele := val.Elem()

	if ele.Kind() != reflect.Slice {
		return fmt.Errorf("传递的对象必须为切片指针")
	}

	if ele.Type().Elem().Kind() != reflect.Struct {
		return fmt.Errorf("传递的对象必须为结构体切片指针")
	}

	//e := ele.Type().Elem() //得到单个结构体type

	v := reflect.New(ele.Type().Elem()).Elem() //创建一个新的结构体value

	type Result struct {
		Id    string  `json:"id"`
		Score float64 `json:"score"`
	}
	//results := []Result{}

	for _, item := range res.Hits {

		doc, _ := index.Document(item.ID)

		for _, field := range doc.Fields {
			setField(v, field.Name(), string(field.Value()))
		}

		for fragmentField, fragments := range item.Fragments {
			for _, fragment := range fragments {
				setField(v, fragmentField, fragment)
			}
		}
		ele = reflect.Append(ele, v)
	}

	if !ele.IsNil() {
		val.Elem().Set(ele)
	}

	return nil
}

// 设置结构体的值
func setField(structValue reflect.Value, name string, value interface{}) error {
	structFieldValue := structValue.FieldByName(name) //结构体单个属性值

	if !structFieldValue.IsValid() {
		log.Printf("[WARN] 没有找到字段: %s \n", name)
		return nil
	}

	if !structFieldValue.CanSet() {
		log.Printf("[WARN] 字段类型不可被修改：%s \n", name)
		return nil
	}

	structFieldType := structFieldValue.Type() //结构体的类型
	val := reflect.ValueOf(value)              //map值的反射值

	var err error
	if structFieldType != val.Type() {
		valStr := fmt.Sprintf("%v", value)
		ntype := structFieldValue.Type().Name()
		//如果是空字符串则只能Set到string类型，其他类型跳过
		if valStr == "" && ntype != "string" {
			return nil
		}
		val, err = typeConversion(valStr, ntype) //类型转换
		if err != nil {
			return err
		}
	}

	structFieldValue.Set(val)
	return nil
}

//类型转换
func typeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "bool" {
		i, err := strconv.ParseBool(value)
		return reflect.ValueOf(bool(i)), err
	}

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype + "\t 值：" + value)
}
