package util

import (
	"bytes"
	"crypto/md5"
	"fmt"
)

//左补齐字符串：原字符串不足n位则在左侧用c补齐至n位，如原字符串超过n位则直接返回
/**
*  srcStr	原字符串
*  c 		补齐的字符
*  length	最后获得的长度
**/
func LeftPad(srcStr string, c rune, n int) string {
	length := len(srcStr)
	if length == 0 {
		return srcStr
	}
	if length >= n {
		return srcStr
	}
	max := n - length //获取补齐的长度

	var buffer bytes.Buffer

	for i := 0; i < max; i++ {
		buffer.WriteRune(c)
	}

	buffer.WriteString(srcStr)

	return buffer.String()

}

//截取字符串
/**
*  start 		起点下标
*  length 	需要截取的长度
**/
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//加密字符串(MD5 32位)
/**
*  str 		要加密的字符串
**/
func MD5Hex(str string) string {
	srcData := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(srcData))
}
