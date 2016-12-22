package data

import (
	"beegostudy/models"
	"beegostudy/util/dateutil"
	"beegostudy/util/fileutil"
	"beegostudy/util/stringutil"
	"github.com/astaxie/beego"
	"path"
)

type UploadController struct {
	DataController
}

func (c *UploadController) Upload() {
	if files, ok := c.FileMap["fileupload"]; ok {
		for _, file := range files {

			filepath := "/upload/"

			filepath += dateutil.GetYMDPathString()

			//检查目录是否存在，不存在则创建
			if !fileutil.IsDir(beego.AppPath + filepath) {
				fileutil.CreateDir(beego.AppPath + filepath)
			}

			newfilename := stringutil.GetUUID() + path.Ext(file.Filename)

			if f, err := file.Open(); err == nil {
				err = fileutil.WriteFileByReadCloser(beego.AppPath+filepath+newfilename, f)
				if err != nil {
					beego.Error(err)
					c.fail("上传失败")
					c.ServeJSON()
					return
				}

				filesize, _ := c.GetInt64("size")

				sysuser := c.GetSession("User").(*models.User)

				m := models.AddAttachment(file.Filename, newfilename, filepath, c.GetString("type"), filesize, sysuser.GetUserName())
				c.put("File", m)
				c.success("上传成功")
			} else {
				beego.Error(err)
				c.fail("上传失败")
			}
		}
	}
	c.ServeJSON()
}
