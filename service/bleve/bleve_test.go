package main

import (
	"fmt"
	//"os"
	"encoding/json"
	"errors"
	"github.com/blevesearch/bleve"
	"github.com/yanyiwu/gojieba"
	_ "github.com/yanyiwu/gojieba/bleve"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Test struct {
	Id    string
	Body  string
	Title string
	Type  string
}

func Example() {

	messages := make([]Test, 0)

	for i := 0; i < 20; i++ {
		message := new(Test)
		message.Id = string(i + 1)
		message.Body = "你好，世界"
		message.Type = "1"
		messages = append(messages, *message)
	}

	for _, msg := range messages {
		if err := AddIndex(msg.Id, msg); err != nil {
			panic(err)
			return
		}
	}

	var m []Test

	And("你好").Search(10, 1, &m)

}

const (
	INDEX_DIR  = "bleveindex"
	FENCI_TYPE = "gojieba"
)

var index bleve.Index

//初始化
func init() {
	indexMapping := bleve.NewIndexMapping()

	err := indexMapping.AddCustomTokenizer(FENCI_TYPE,
		map[string]interface{}{
			"dictpath":     gojieba.DICT_PATH,
			"hmmpath":      gojieba.HMM_PATH,
			"userdictpath": gojieba.USER_DICT_PATH,
			"idf":          gojieba.IDF_PATH,
			"stop_words":   gojieba.STOP_WORDS_PATH,
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

func (b *Bleve) Search(size, page int, obj interface{}) {

	log.Println("SearchQuerys:", b.querys)

	page--

	from := size * page

	query := strings.Join(b.querys, " ")

	que := bleve.NewQueryStringQuery(query)

	fmt.Println(size, from)

	req := bleve.NewSearchRequestOptions(que, size, from, false)

	req.Highlight = bleve.NewHighlight()

	res, err := index.Search(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	fmt.Println(prettify(res, obj))
}

func prettify(res *bleve.SearchResult, obj interface{}) string {

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

	e := ele.Type().Elem() //得到单个结构体type

	v := reflect.New(ele.Type().Elem()).Elem() //创建一个新的结构体value

	type Result struct {
		Id    string  `json:"id"`
		Score float64 `json:"score"`
	}
	results := []Result{}

	for _, item := range res.Hits {

		for fragmentField, fragments := range item.Fragments {
			fmt.Println("---", fragmentField)
			for _, fragment := range fragments {
				fmt.Println("===", fragment)
			}
		}
		results = append(results, Result{item.ID, item.Score})
	}
	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return string(b)

	/////

	// setField(v, name, text)

	// *ele = reflect.Append(*ele, v)

	// if !ele.IsNil() {
	//  val.Elem().Set(ele)
	// }

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

// func prettify(res *bleve.SearchResult) string {
//  type Result struct {
//      Id    string  `json:"id"`
//      Score float64 `json:"score"`
//  }
//  results := []Result{}

//  for _, item := range res.Hits {

//      for fragmentField, fragments := range item.Fragments {
//          fmt.Println(fragmentField)
//          for _, fragment := range fragments {
//              fmt.Println(fragment)
//          }
//      }

//      for otherFieldName, otherFieldValue := range item.Fields {
//          if _, ok := item.Fragments[otherFieldName]; !ok {
//              fmt.Println(otherFieldName)
//              fmt.Println(otherFieldValue)
//          }
//      }

//      results = append(results, Result{item.ID, item.Score})
//  }
//  b, err := json.Marshal(results)
//  if err != nil {
//      panic(err)
//  }
//  return string(b)
// }

func main() {
	Example()
}
