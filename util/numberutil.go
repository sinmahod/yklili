package util

import (
	"math/rand"
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
