package fileutil

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// 文件末尾追加
// fileName:文件名字(带全路径)
// content: 写入的内容
func AppendWriteFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	defer f.Close()
	if err == nil {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}

	return err
}

// 写入文件
func WriteFileByReadCloser(fileName string, file io.ReadCloser) error {
	defer file.Close()
	dstFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, file)
	return err
}

// 写入文件
func WriteFileByReader(fileName string, file io.Reader) error {
	dstFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, file)
	return err
}

// 写入文件
func WriteFileByByte(fileName string, data []byte) error {
	os.MkdirAll(path.Dir(fileName), os.ModePerm)
	return ioutil.WriteFile(fileName, data, 0655)
}

// 得到文件大小
func FileSize(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

//递归创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// 是否是文件
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// 是否是文件夹
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// 检查文件或目录是否存在
// 如果由 fileName 指定的文件或目录存在则返回 true，否则返回 false
func Exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

// 读取文件内容
func FileToString(filePath string) (string, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd), nil
}

// 读取XML到结构
func XMLToStruct(filePath string, result interface{}) error {
	content, err := ioutil.ReadFile(filePath)
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
func XMLStructToFile(filePath string, result interface{}) error {
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
	return ioutil.WriteFile(filePath, xmlOutPutData, os.ModeAppend)
}

func GetFilelist(path string) []string {
	var files []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil
	}
	return files
}
