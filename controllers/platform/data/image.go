package data

import (
	"beegostudy/util/dateutil"
	"beegostudy/util/fileutil"
	"beegostudy/util/stringutil"
	"fmt"
	"github.com/astaxie/beego"
	"path"
)

type ImageController struct {
	DataController
}

func (c *ImageController) UploadTest() {
	c.TplName = "platform/image/uploadDialog.html"
	c.addScript()
}

type File struct {
	//文件名称
	Name string
	//文件路径
	Path string
	//文件大小
	Size int64
	//文件类型
	Type string
}

func (c *ImageController) Upload() {

	if files, ok := c.FileMap["fileupload"]; ok {
		fs := new(File)
		for _, file := range files {
			filepath := beego.AppPath + "/upload/"

			filepath += dateutil.GetYMDPathString()

			//检查目录是否存在，不存在则创建
			if !fileutil.IsDir(filepath) {
				fileutil.CreateDir(filepath)
			}

			filepath += stringutil.GetUUID() + path.Ext(file.Filename)

			if f, err := file.Open(); err == nil {
				fileutil.WriteFileByReadCloser(filepath, f)
				fs.Name = file.Filename
				fs.Path = filepath
				fs.Size, _ = c.GetInt64("size")
				fs.Type = c.GetString("type")
				c.put("File", fs)
				c.success("上传成功")
			} else {
				fmt.Println(err)
				c.fail("上传失败")
			}
		}
	}
	c.ServeJSON()
}
