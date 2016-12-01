package util

import (
	"io/ioutil"
	"os"
)

func FileToString(filepath string) string {
	fi, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}
