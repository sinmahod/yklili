package front

import (
	"io/ioutil"
	"os"

	"path"
	"strconv"
	"strings"

	"github.com/sinmahod/yklili/conf"
	"github.com/sinmahod/yklili/util/fileutil"
	"github.com/sinmahod/yklili/util/imageutil"
)

type ImageController struct {
	FrontController
}

func (c *ImageController) Get() {

	c.Ctx.Output.Header("Content-Type", "image/jpeg")
	c.Ctx.Output.Header("Accept-Ranges", "bytes")

	splat := c.Ctx.Input.Param(":splat")

	if splat == "" {
		return
	}

	param := strings.Split(splat, "!")

	uploadpath := conf.GetValue(conf.UploadPath)

	filepath := uploadpath + param[0]

	//限定uoload文件夹，防止读取系统文件
	if !fileutil.Exist(filepath) {
		c.GetDefaultImage()
		return
	}

	if len(param) > 1 {
		// 展示缩放图
		wh := strings.Split(param[1], "x")
		if len(wh) > 1 {
			// 先判断图片是否存在
			ext := path.Ext(param[0]) // 获取扩展名

			// 如果宽高不是数字则展示默认图片
			width, err := strconv.Atoi(wh[0])
			if err != nil {
				c.GetDefaultImage()
				return
			}
			height, err := strconv.Atoi(wh[1])
			if err != nil {
				c.GetDefaultImage()
				return
			}
			idx := strings.LastIndex(filepath, ".")
			dstpath := filepath[:idx] + "_" + wh[0] + "x" + wh[1] + ext

			// 指定格式的图片不存在则生成
			if !fileutil.Exist(dstpath) {
				imageutil.ImageCut(filepath, width, height, dstpath)
			}
			filepath = dstpath
		}
	}

	fi, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	rf, _ := ioutil.ReadAll(fi)

	c.Ctx.ResponseWriter.Write(rf)

}

// 默认图片
func (c *ImageController) GetDefaultImage() {
	fi, err := os.Open(conf.GetDefaultImage())
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	rf, _ := ioutil.ReadAll(fi)
	c.Ctx.ResponseWriter.Write(rf)
}
