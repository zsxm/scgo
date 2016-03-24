package tools

import (
	"os"
	"strconv"
	"strings"
)

func IsBlank(s string) bool {
	if strings.Trim(s, " ") == "" || s == "nil" || s == "null" {
		return true
	}
	return false
}
func IsNotBlank(s string) bool {
	if strings.TrimSpace(s) != "" && s != "nil" && s != "null" {
		return true
	}
	return false
}

func ParseInteger(v string) int {
	s, _ := strconv.Atoi(v)
	return s
}

//判断文件是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
