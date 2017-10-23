package imageutil

import (
	"github.com/disintegration/imaging"
	"log"
)

// 缩放图片  达到指定比例，超出部分截掉保留中间部分
func ImageCut(filepath string, w, h int, dstpath string) error {
	src, err := imaging.Open(filepath)
	if err == nil {
		log.Println("图片固定比例（不拉伸挤压）缩放处理中", dstpath)
		dst := imaging.Fill(src, w, h, imaging.Center, imaging.Lanczos)
		return imaging.Save(dst, dstpath)
	}
	return err
}

// 缩放图片  根据宽度及高度保证达到指定比例（会拉伸挤压）
func ImageCut2(filepath string, w, h int, dstpath string) error {
	src, err := imaging.Open(filepath)
	if err == nil {
		log.Println("图片缩放（宽高固定）处理中", dstpath)
		dst := imaging.Resize(src, w, h, imaging.Lanczos)
		return imaging.Save(dst, dstpath)
	}
	return err
}

// 缩放图片 根据宽度或高度自适配（保证图片不拉伸挤压并切高宽不会超过指定比例）
func ImageCut3(filepath string, w, h int, dstpath string) error {
	src, err := imaging.Open(filepath)
	if err == nil {
		log.Println("图片缩放（宽高自适配）处理中", dstpath)
		dst := imaging.Fit(src, w, h, imaging.Lanczos)
		return imaging.Save(dst, dstpath)
	}
	return err
}

// 直接在原图剪切指定位置指定比例的图片，不做任何缩放
func ImageCut4(filepath string, w, h int, dstpath string) error {
	src, err := imaging.Open(filepath)
	if err == nil {
		log.Println("图片剪切处理中", dstpath)
		dst := imaging.CropAnchor(src, w, h, imaging.Center)
		return imaging.Save(dst, dstpath)
	}
	return err
}

// 图片反色
func Invert(filepath, dstpath string) error {
	src, err := imaging.Open(filepath)
	if err == nil {
		log.Println("图片反色处理中", dstpath)
		dst := imaging.Invert(src)
		return imaging.Save(dst, dstpath)
	}
	return err
}

// 图片去色
func Decolourize(filepath, dstpath string) error {
	src, err := imaging.Open(filepath)
	if err == nil {
		log.Println("图片去色处理中", dstpath)
		dst := imaging.Grayscale(src)
		dst = imaging.AdjustContrast(dst, 20)
		dst = imaging.Sharpen(dst, 2)
		return imaging.Save(dst, dstpath)
	}
	return err
}

// 图片模糊
func Blur(filepath string, blur float64, dstpath string) error {
	src, err := imaging.Open(filepath)
	if err == nil {
		log.Println("图片模糊处理中", dstpath)
		dst := imaging.Blur(src, blur)
		return imaging.Save(dst, dstpath)
	}
	return err
}
