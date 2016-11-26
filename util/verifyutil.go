package util

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

const (
	NOTNULL      = "required"    //不能为空
	EMAIL        = "email"       //Email格式
	URL          = "url"         //url格式
	NEMBER       = "number"      //合法数字
	DIGITS       = "digits"      //整数
	MAX_LENGTH   = "mexlength"   //最大长度
	MIN_LENGTH   = "minlength"   //最小长度
	RANGE_LENGTH = "rangelength" //长度区间
	RANGE        = "range"       //数值区间
	MAX          = "max"         //最大数值
	MIN          = "min"         //最小数值
)

const (
	MSG_NOTNULL      = "这是必填字段"                  //不能为空
	MSG_EMAIL        = "请输入有效的电子邮件地址"            //Email格式
	MSG_URL          = "请输入有效的网址"                //url格式
	MSG_NEMBER       = "请输入有效的网址"                //合法数字
	MSG_DIGITS       = "只能输入数字"                  //整数
	MSG_MAX_LENGTH   = "最多可以输入 {0} 个字符"          //最大长度
	MSG_MIN_LENGTH   = "最少要输入 {0} 个字符"           //最小长度
	MSG_RANGE_LENGTH = "请输入长度在 {0} 到 {1} 之间的字符串" //长度区间
	MSG_RANGE        = "请输入范围在 {0} 到 {1} 之间的数值"  //数值区间
	MSG_MSG_MAX      = "请输入不大于 {0} 的数值"          //最大数值
	MSG_MIN          = "请输入不小于 {0} 的数值"          //最小数值

	//0 Form的ID
	//1 rules
	//2 messages
	SCRIPT = "jQuery(function($){$('#{0}').validate({errorClass: 'help-block',focusInvalid: false," +
		"rules:{{1}},messages:{{2}}," +
		"highlight: function (e) {" +
		"$(e).closest('.form-group').removeClass('has-info').addClass('has-error');" +
		"}," +

		"success: function (e) {" +
		"$(e).closest('.form-group').removeClass('has-error');" +
		"$(e).remove();" +
		"}," +
		"errorPlacement: function (error, element) {" +
		"if(element.is('input[type=checkbox]') || element.is('input[type=radio]')) {" +
		"var controls = element.closest('div[class*=\"col-\"]');" +
		"if(controls.find(':checkbox,:radio').length > 1) controls.append(error);" +
		"else error.insertAfter(element.nextAll('.lbl:eq(0)').eq(0));" +
		"}" +
		"else if(element.is('.select2')) {" +
		"error.insertAfter(element.siblings('[class*=\"select2-container\"]:eq(0)'));" +
		"}" +
		"else if(element.is('.chosen-select')) {" +
		"error.insertAfter(element.siblings('[class*=\"chosen-container\"]:eq(0)'));" +
		"}" +
		"else error.insertAfter(element);" +
		"}," +

		"submitHandler: function (form) {" +
		"}," +
		"invalidHandler: function (form) {" +
		"}" +
		"});" +
		"});"
)

const (
	TEXT = "<div class=\"form-group\">" +
		"<label class=\"col-sm-3 control-label no-padding-right\">" +
		"${1}" +
		"</label>" +
		"<div class=\"col-sm-9\">" +
		"<input ${2} />" +
		"</div>" +
		"</div>"
)

func AddVerifyJs(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	doc.Find("input[verify]").Each(func(i int, node *goquery.Selection) {
		id, _ := node.Attr("id")

		verify, _ := node.Attr("verify")
		fmt.Println("================222", verify, id)
		fmt.Println(node.Html())
	})

	doc.Find("goinput").Each(func(i int, node *goquery.Selection) {
		id, _ := node.Attr("id")
		text, _ := node.Html()
		verify, _ := node.Attr("verify")

		re := regexp.MustCompile("(\\w+)")
		//str := re.ReplaceAllString(TEXT, "111")
		is := []int{1000}
		fmt.Printf("%q", re.ExpandString(nil, TEXT, "Golang,World!", is))

		var buffer bytes.Buffer

		// buffer.WriteString(str)
		buffer.WriteString("\n")
		// buffer.WriteString(Script)

		for _, ss := range node.Nodes {
			for k, v := range ss.Attr {
				fmt.Println(k, v)
			}
		}

		node.ReplaceWithHtml("aaa")

		fmt.Println("================111", verify, id, text)
	})
	//fmt.Println(doc.Html())
	return "", nil
}
