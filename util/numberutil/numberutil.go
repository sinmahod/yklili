package numberutil

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//获取随机整数，取值范围为0～max-1
/**
*  max 		最大取值范围(-1)
**/
func RandInt(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

//判断是否为整数
func IsNumber(i interface{}) bool {
	if i != nil {
		tmp := fmt.Sprintf("%v", i)
		if _, err := strconv.Atoi(tmp); err == nil {
			return true
		}
	}
	return false
}

//转换为整数
func Atoi(i interface{}) int {
	tmp := fmt.Sprintf("%v", i)
	if i, err := strconv.Atoi(tmp); err == nil {
		return i
	} else {
		return 0
	}
}
