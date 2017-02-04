package bleve

import (
	"fmt"
	//"os"
	"encoding/json"
	"github.com/blevesearch/bleve"
	"github.com/yanyiwu/gojieba"
	_ "github.com/yanyiwu/gojieba/bleve"
)

const INDEX_DIR = "bleveindex"

var index bleve.Index

//初始化
func init() {
	indexMapping := bleve.NewIndexMapping()

	err := indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":     gojieba.DICT_PATH,
			"hmmpath":      gojieba.HMM_PATH,
			"userdictpath": gojieba.USER_DICT_PATH,
			"idf":          gojieba.IDF_PATH,
			"stop_words":   gojieba.STOP_WORDS_PATH,
			"type":         "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	err = indexMapping.AddCustomAnalyzer("gojieba",
		map[string]interface{}{
			"type":      "gojieba",
			"tokenizer": "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}

	indexMapping.DefaultAnalyzer = "gojieba"

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
func Search() {

}

func Example() {
	messages := []struct {
		Id    string
		Body  string
		Title string
		Type  int
	}{
		{
			Id:    "6",
			Body:  "你好，世界dafdsafdsfasd",
			Title: "你好",
			Type:  1,
		},
		{
			Id:    "2",
			Body:  "世界你好啊fewfwef",
			Title: "你好",
			Type:  2,
		},
		{
			Id:    "3",
			Body:  "sdfsdaf亲口",
			Title: "你好",
			Type:  2,
		},
		{
			Id:    "4",
			Body:  "交代cdscsdc",
			Title: "你好",
			Type:  1,
		},
	}

	querys := []string{
		"你好，世界",
		"亲口交代",
	}

	for _, q := range querys {
		bleve.NewPhraseQuery(match)

		que := bleve.NewQueryStringQuery(q)

		req := bleve.NewSearchRequest(que)
		req.Highlight = bleve.NewHighlight()
		res, err := index.Search(req)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		fmt.Println(prettify(res))
	}
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

func main() {
	Example()
}
