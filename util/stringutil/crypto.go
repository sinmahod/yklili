/**
 * func main() {
 *   ss := encrypt("aaaaaaaaaa22a", "test")
 *   fmt.Println(ss)
 *   fmt.Println(decrypt(ss, "test"))
 * }
 */
package stringutil

import (
	"crypto/cipher"
	"crypto/des"
	"crypto/sha256"
	"encoding/base64"
)

const DEFAULT_KEY = "27jrWz3sxrVbR+pnyg6j"

/**
 * 字符串加密
 */
func Encrypt(src string, key ...string) string {
	var ks string
	if len(key) == 0 {
		ks = DEFAULT_KEY
	} else {
		ks = key[0]
	}
	arg := sha256.Sum224([]byte(ks))
	k := arg[:24]
	s := []byte(src)
	block, _ := des.NewTripleDESCipher(k)
	EncryptMode := cipher.NewCBCEncrypter(block, k[:8])
	s = pkcs5append(s)
	EncryptMode.CryptBlocks(s, s)
	src = base64.StdEncoding.EncodeToString(s)
	return src
}

/**
 * 字符串解密
 */
func Decrypt(src string, key ...string) string {
	var ks string
	if len(key) == 0 {
		ks = DEFAULT_KEY
	} else {
		ks = key[0]
	}
	arg := sha256.Sum224([]byte(ks))
	k := arg[:24]
	s := []byte(src)
	block, _ := des.NewTripleDESCipher(k)
	DecryptMode := cipher.NewCBCDecrypter(block, k[:8])
	s, _ = base64.StdEncoding.DecodeString(src)
	DecryptMode.CryptBlocks(s, s)
	s = pkcs5remove(s)
	return string(s)
}

func pkcs5append(plaintext []byte) []byte {
	num := 8 - len(plaintext)%8
	for i := 0; i < num; i++ {
		plaintext = append(plaintext, byte(num))
	}
	return plaintext
}

func pkcs5remove(plaintext []byte) []byte {
	length := len(plaintext)
	num := int(plaintext[length-1])
	return plaintext[:(length - num)]
}
