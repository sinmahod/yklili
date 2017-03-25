package pwdutil

import (
	"github.com/sinmahod/yklili/util/numberutil"
	"github.com/sinmahod/yklili/util/stringutil"
	"strconv"
	"strings"
)

//密码生成与密码校验

//生成随机的盐密码：按指定规则改变md5，可局部反向
func GeneratePWD(password string) string {

	salt := stringutil.LeftPad(strconv.Itoa(numberutil.RandInt(99999999)), '0', 8) + stringutil.LeftPad(strconv.Itoa(numberutil.RandInt(99999999)), '0', 8)

	password = stringutil.MD5Hex(password + salt)
	ss := []rune(salt)
	ps := []rune(password)
	cs := make([]rune, 48)

	for i := 0; i < 48; i += 3 {
		cs[i] = ps[i/3*2]
		c := ss[i/3]
		cs[i+1] = c
		cs[i+2] = ps[i/3*2+1]
	}
	return string(cs)

}

//校验密码是否正确
/**
*  password 		密码（明文）
*  md5 			密码（加密）
**/
func VerifyPWD(password, md5 string) bool {
	md5rune := []rune(md5)
	cs1 := make([]rune, 32)
	cs2 := make([]rune, 16)
	for i := 0; i < 48; i += 3 {
		cs1[i/3*2] = md5rune[i]
		cs1[i/3*2+1] = md5rune[i+2]
		cs2[i/3] = md5rune[i+1]
	}
	salt := string(cs2)
	return strings.EqualFold(stringutil.MD5Hex(password+salt), string(cs1))
}
