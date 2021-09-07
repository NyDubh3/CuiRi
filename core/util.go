package core

import (
	"os"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func RemoveSpecialCharactar(filedata string) string {
	filedata = strings.Replace(filedata, "\\x", "",-1)
	filedata = strings.Replace(filedata, "\"", "",-1)
	filedata = strings.Replace(filedata, " ", "",-1)
	filedata = strings.Replace(filedata, "\r\n", "",-1)
	filedata = strings.Replace(filedata, ";", "",-1)
	return filedata
}
