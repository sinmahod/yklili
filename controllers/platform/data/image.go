package data

import (
	"beegostudy/conf"
	"beegostudy/models"
	"beegostudy/util/stringutil"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
)

type ImageController struct {
	DataController
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

func (c *ImageController) Img() {
	sn := c.GetString("sn")

	uploadpath := conf.GetValue(conf.UploadPath)

	if uploadpath == "" {
		uploadpath = beego.AppPath
	}

	fi, err := os.Open(uploadpath + stringutil.Decrypt(stringutil.Decode(sn)))
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	rf, _ := ioutil.ReadAll(fi)

	c.Ctx.Output.Header("Content-Type", "image/jpeg")
	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.ResponseWriter.Write(rf)
}
