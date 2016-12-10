package fileutil

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 读取文件内容
func FileToString(filepath string) string {
	fi, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

// 读取XML到结构
func XMLToStruct(filepath string, result interface{}) error {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return xml.Unmarshal(content, result)
}

/**
 *   将结构写入xml文件 (必须是xml结构) 结构内部XMLName，innerxml固定，内部属性必须导出
 *------------------------------------Struct--------------------------------------
 *	type StringResources struct {
 *		XMLName        xml.Name         `xml:"resources"`
 *		ResourceString []ResourceString `xml:"string"`
 *	}
 *
 *	type ResourceString struct {
 *		XMLName    xml.Name `xml:"string"`
 *		StringName string   `xml:"name,attr"`
 *		InnerText  string   `xml:",innerxml"`
 *	}
 *
 *------------------------------------XML------------------------------------------
 *
 *   <?xml version="1.0" encoding="UTF-8"?>
 *   <resources>
 *	<string name="VideoLoading">Loading video…</string>
 *	<string name="ApplicationName">what</string>
 *   </resources>
 *------------------------------------------------------------------------------------
 */
func XMLStructToFile(filepath string, result interface{}) error {
	//保存修改后的内容
	xmlOutPut, err := xml.MarshalIndent(result, "    ", "")
	if err != nil {
		return err
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, xmlOutPut...)
	//写入文件
	return ioutil.WriteFile(filepath, xmlOutPutData, os.ModeAppend)
}
