package data

import (
	"beegostudy/models"
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

func (c *ImageController) Upload() {

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
					fmt.Println(err)
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
				fmt.Println(err)
				c.fail("上传失败")
			}
		}
	}
	c.ServeJSON()
}

//DataGrid列表数据加载
func (c *ImageController) List() {
	c.RequestData["filetype"] = "image/*"
	if datagrid, err := models.GetAttchmentsPage(c.PageSize, c.PageIndex, c.OrderColumn, c.OrderSord, c.RequestData); err != nil {
		beego.Error(err)
	} else {
		c.Data["json"] = datagrid
		c.ServeJSON()
	}
}
