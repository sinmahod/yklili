package platform

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
)

type UploadController struct {
	beego.Controller
}

func (c *UploadController) Prepare() {
	if err := c.Ctx.Request.ParseForm(); err != nil {
		beego.Info(err)
	}

	for k, v := range c.Ctx.Request.Form {
		if len(v) > 0 {
			beego.Info(k, "=========", v)
		}
	}

	filemap := c.Ctx.Request.MultipartForm.File
	if len(filemap) > 0 {
		for k, v := range filemap {
			beego.Info(k, v, len(v))
			for i, vv := range v {
				beego.Info(vv.Filename)
				file, err := vv.Open()
				if err != nil {
					fmt.Println(err)
					return
				}
				defer file.Close()

				dstFile, err := os.Create("c:\\" + strconv.Itoa(i) + ".png")

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				defer dstFile.Close()

				io.Copy(dstFile, file)
			}
			reg, _ := regexp.Compile(`\[[\d]+?\]`)
			key := reg.ReplaceAllString(k, "")
			beego.Info(key)
		}
	}

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// file, handler, err = c.Ctx.Request.FormFile("file")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer file.Close()
	// fmt.Println(handler.Filename)
	// f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer f.Close()
	// io.Copy(f, file)

}
