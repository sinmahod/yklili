package bleve

import (
	"fmt"
	//"os"
	"encoding/json"
	"github.com/blevesearch/bleve"
	"github.com/yanyiwu/gojieba"
	_ "github.com/yanyiwu/gojieba/bleve"
	"strings"
)

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

	//index, err := bleve.New(INDEX_DIR, indexMapping)
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
}

func (b *Bleve) And(query string, column ...string) *Bleve {
	if len(column) > 0 {
		for _, c := range column {
			b.querys = append(b.querys, "+"+c+":"+query)
		}
	} else {
		b.querys = append(b.querys, "+"+query)
	}
}

func (b *Bleve) Not(query string, column ...string) *Bleve {
	if len(column) > 0 {
		for _, c := range column {
			b.querys = append(b.querys, "-"+c+":"+query)
		}
	} else {
		b.querys = append(b.querys, "-"+query)
	}
}

func (b *Bleve) Or(query string, column ...string) *Bleve {
	if len(column) > 0 {
		for _, c := range column {
			b.querys = append(b.querys, c+":"+query)
		}
	} else {
		b.querys = append(b.querys, query)
	}
}

func (b *Bleve) Search(size, index int, obj interface{}) []interface{} {
	index--

	from := size * index

	query := strings.Join(b.querys, " ")

	que := bleve.NewQueryStringQuery(query)

	req := bleve.NewSearchRequestOptions(query, size, from, false)

	req.Highlight = bleve.NewHighlight()

	res, err := index.Search(req)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	fmt.Println(prettify(res))
}

func prettify(res *bleve.SearchResult) string {
	type Result struct {
		Id    string  `json:"id"`
		Score float64 `json:"score"`
	}
	results := []Result{}
	for _, item := range res.Hits {
		results = append(results, Result{item.ID, item.Score})
	}
	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return string(b)
}
