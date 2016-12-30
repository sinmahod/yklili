package stringutil

import (
	"encoding/hex"
)

/**
 * 字符串Encode
 */
func Encode(src string) string {
	return hex.EncodeToString([]byte(src))
}

/**
 * 字符串Decode
 */
func Decode(src string) string {
	s, _ := hex.DecodeString(src)
	return string(s)
}
