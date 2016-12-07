package util

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	TEXT = "<div class=\"form-group\">" +
		"<label class=\"col-sm-3 control-label no-padding-right\">" +
		"$1" +
		"</label>" +
		"<div class=\"col-sm-9\">" +
		"<input $2 />" +
		"</div>" +
		"</div>"
)

func AnalysisGoTag(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	doc.Find("gotext").Each(func(i int, node *goquery.Selection) {
		text, _ := node.Html()
		param := make(params)
		for _, ss := range node.Nodes {
			for _, attr := range ss.Attr {
				param[strings.ToLower(attr.Key)] = attr.Val
			}
		}

		newhtml := replaceContent(TEXT, text, param.join())

		node.ReplaceWithHtml(newhtml) //替换新的html
	})

	return doc.Html()
}

type params map[string]string

func (p params) join() string {
	if p["name"] == "" && p["id"] != "" {
		p["name"] = p["id"]
	}
	var buffer bytes.Buffer
	for k, v := range p {
		buffer.WriteString(" ")
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString("\"")
		buffer.WriteString(v)
		buffer.WriteString("\"")
	}
	return buffer.String()
}

//传入原内容与参数，按分组替换为新的内容
func replaceContent(content string, str ...string) string {
	var regstr bytes.Buffer
	for i, _ := range str {
		if i > 0 {
			regstr.WriteString(`\$_\$`)
		}
		regstr.WriteString(`([\S\s]+)`)
	}
	reg := regexp.MustCompile(regstr.String())
	src := strings.Join(str, "$_$")           // 源文本
	match := reg.FindStringSubmatchIndex(src) // 解析源文本
	return string(reg.ExpandString(nil, content, src, match)[:])
}
