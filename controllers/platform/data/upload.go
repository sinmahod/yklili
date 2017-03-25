package data

import (
	"github.com/astaxie/beego"
	"path"
	"yklili/conf"
	"yklili/models"
	"yklili/util/dateutil"
	"yklili/util/fileutil"
	"yklili/util/stringutil"
)

type UploadController struct {
	DataController
}

func (c *UploadController) Upload() {
	if files, ok := c.FileMap["fileupload"]; ok {
		for _, file := range files {

			filepath := "/upload/"

			filepath += dateutil.GetYMDPathString()

			uploadpath := conf.GetValue(conf.UploadPath)

			if uploadpath == "" {
				uploadpath = beego.AppPath
			}

			//检查目录是否存在，不存在则创建
			if !fileutil.IsDir(uploadpath + filepath) {
				fileutil.CreateDir(uploadpath + filepath)
			}

			newfilename := stringutil.GetUUID() + path.Ext(file.Filename)

			if f, err := file.Open(); err == nil {
				err = fileutil.WriteFileByReadCloser(uploadpath+filepath+newfilename, f)
				if err != nil {
					beego.Error(err)
					c.fail("上传失败")
					c.ServeJSON()
					return
				}

				filesize, _ := c.GetInt64("size")

				sysuser := c.GetSession("User").(*models.S_User)

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
