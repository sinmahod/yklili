package bleve

import (
	"fmt"
	"strconv"
	"testing"
)

type Test struct {
	Id    string
	Body  string
	Title string
	Type  int
}

func TestAddIndex(t *testing.T) {

	messages := make([]Test, 0)

	for i := 40; i < 60; i++ {
		message := new(Test)
		message.Id = strconv.Itoa(i + 1)
		message.Body = "jdlsjflkdsjlf初始化ldsjlfjaksdjflkjdlsnlcndslanvjksdnvkjbdskabviwheibibsdklckaqwoeijoiuorejwtrewbmbmzbxvcxkhcvhohoi"
		message.Title = "jqGrid封装-DataGrid"
		message.Type = 1
		messages = append(messages, *message)
	}

	for _, msg := range messages {
		if err := AddIndex(msg.Id, msg); err != nil {
			panic(err)
			return
		}
	}

	var as []Test

	Or("jqGrid", "Title", "Content").Search(3, 1, &as)

	fmt.Println(len(as), "+=======")

	for i, _ := range as {
		// as[i].Content = strings.Replace(as[i].Content, "<mark>", "{MARK}", -1)
		// as[i].Content = strings.Replace(as[i].Content, "</mark>", "{/MARK}", -1)
		// //转义html
		// as[i].Content = template.HTMLEscapeString(as[i].Content)
		// as[i].Content = strings.Replace(as[i].Content, "{MARK}", "<mark>", -1)
		// as[i].Content = strings.Replace(as[i].Content, "{/MARK}", "</mark>", -1)
		//
		fmt.Println(as[i])
	}

}
