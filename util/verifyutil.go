package util

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	NOTNULL      = "required"    //不能为空
	EMAIL        = "email"       //Email格式
	URL          = "url"         //url格式
	NEMBER       = "number"      //合法数字
	DIGITS       = "digits"      //整数
	MAX_LENGTH   = "maxlength"   //最大长度
	MIN_LENGTH   = "minlength"   //最小长度
	RANGE_LENGTH = "rangelength" //长度区间
	RANGE        = "range"       //数值区间
	MAX          = "max"         //最大数值
	MIN          = "min"         //最小数值
	EQUALTO      = "equalTo"     //两次密码输入不一致 equalTo("#ID")
)

var VERIFY_MAP = map[string]string{
	NOTNULL:      "这是必填字段",
	EMAIL:        "请输入有效的电子邮件地址",
	URL:          "请输入有效的网址",
	NEMBER:       "请输入有效的数字",
	DIGITS:       "只能输入数字",
	MAX_LENGTH:   "最多可以输入 {0} 个字符",
	MIN_LENGTH:   "最少要输入 {0} 个字符",
	RANGE_LENGTH: "请输入长度在 {0} 到 {1} 之间的字符串",
	RANGE:        "请输入范围在 {0} 到 {1} 之间的数值",
	MAX:          "请输入不大于 {0} 的数值",
	MIN:          "请输入不小于 {0} 的数值",
	EQUALTO:      "两次密码输入不一致",
}

const (
	//1 Form的ID
	//2 rules
	//3 messages
	SCRIPT = "<script>jQuery(function($){$('#$1').validate({errorClass: 'help-block',focusInvalid: false," +
		"rules:{$2},messages:{$3}," +
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
		"});</script>"
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

	vjs := new(validate)
	doc.Find("form").EachWithBreak(func(i int, node *goquery.Selection) bool {
		if id, ok := node.Attr("id"); ok {
			vjs.FormId = id
			return false
		}
		return true
	})

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

		//添加validate
		if vjs.FormId != "" {
			if verify, ok := node.Attr("verify"); ok {
				verifyForm(vjs, param["name"], verify)
			}
		}
	})
	doc.Find("body").EachWithBreak(func(i int, node *goquery.Selection) bool {
		if vjs.FormId != "" {
			// verifydate := vjs.addVerifyJs()
			// javascript := replaceContent(SCRIPT, verifydate...)
			// node.AppendHtml(javascript)
		}
		return false
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

type validate struct {
	FormId string
	CV     []*columnVerify
}

//校验js的结构
type columnVerify struct {
	Name string
	Rams []ruleAndMsg
}

type ruleAndMsg struct {
	Rule string
	Msg  string
}

func (vjs *validate) addVerifyJs() []string {
	var ruleBuffer bytes.Buffer
	var msgBuffer bytes.Buffer
	var i int
	for _, rams := range vjs.CV {
		if i > 0 {
			ruleBuffer.WriteString(",")
			msgBuffer.WriteString(",")
		}
		ruleBuffer.WriteString(rams.Name)
		ruleBuffer.WriteString(":")
		ruleBuffer.WriteString("{")
		msgBuffer.WriteString(rams.Name)
		msgBuffer.WriteString(":")
		msgBuffer.WriteString("{")
		for i, rms := range rams.Rams {
			if i > 0 {
				ruleBuffer.WriteString(",")
				msgBuffer.WriteString(",")
			}
			ruleBuffer.WriteString(rms.Rule)
			msgBuffer.WriteString(rms.Msg)
		}
		ruleBuffer.WriteString("}")
		msgBuffer.WriteString("}")
		i++
	}
	return []string{vjs.FormId, ruleBuffer.String(), msgBuffer.String()}
}

var regexRep = regexp.MustCompile(`(\()(.*)(\))`)
var regexParam = regexp.MustCompile(`([\S\s]+)\(([\S\s]+)\)`)

//传递vjs，字段名，校验字段值构造结构体
func verifyForm(vjs *validate, columnName, verifyStr string) {
	if columnName == "" {
		return
	}

	cv := new(columnVerify)
	rams := make([]ruleAndMsg, 0)
	strs := strings.Split(verifyStr, ";")

	for _, s := range strs {
		ram := new(ruleAndMsg)
		st := regexParam.FindStringSubmatch(s)
		if len(st) == 3 { //有参数
			if value := VERIFY_MAP[st[1]]; value != "" {
				ram.Msg = st[1] + `:"` + value + `"`
				if strings.Index(st[2], ",") > 0 { //区间参数
					ram.Rule = regexRep.ReplaceAllString(s, ":[$2]")
				} else { //单个参数
					ram.Rule = regexRep.ReplaceAllString(s, ":$2")
				}

			} else {
				continue
			}
		} else { //无参数
			if value := VERIFY_MAP[s]; value != "" {
				ram.Rule = s + ":true"
				ram.Msg = s + `:"` + value + `"`
			} else {
				continue
			}
		}
		rams = append(rams, *ram)
	}

	cv.Name = columnName
	cv.Rams = rams
	vjs.CV = append(vjs.CV, cv)
}
