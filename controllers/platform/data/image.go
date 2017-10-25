package data

import (
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sinmahod/yklili/conf"
	"github.com/sinmahod/yklili/models"
	"github.com/sinmahod/yklili/models/orm"
	"github.com/sinmahod/yklili/util/fileutil"
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

func (c *ImageController) InitPage() {
	c.TplName = "platform/image/imageDialog.html"
	c.addScript()
}

func (c *ImageController) Del() {
	id, err := c.GetInt("Id")
	if err != nil {
		c.fail("操作失败，请先确定图片id是否传递正确")
		c.ServeJSON()
		return
	}
	m, err := models.GetArrachment(id)
	if err != nil {
		c.fail("操作失败，请先确定图片id是否传递正确")
		c.ServeJSON()
		return
	}

	tran := new(orm.Transaction)
	tran.Add(m, orm.DELETE)
	if err := tran.Commit(); err != nil {
		beego.Error(err)
		c.fail("操作失败，操作数据库时出现错误")
	} else {
		// 删除物理文件
		uploadpath := conf.GetValue(conf.UploadPath)
		files := fileutil.GetFilelist(uploadpath + m.FilePath)
		srcf := uploadpath + m.FilePath + m.FileNewName
		idx := strings.LastIndex(srcf, ".")
		srcf = srcf[:idx]
		for _, f := range files {
			if strings.Index(f, srcf) > -1 {
				os.Remove(f)
			}
		}
		c.success("操作成功")
	}
	c.ServeJSON()
}
