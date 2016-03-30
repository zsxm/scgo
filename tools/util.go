package tools

import (
	"os"
	"path/filepath"
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

func CurrentDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

func EachDir(dir, pix string) []string {
	htmlTemps := make([]string, 0, 10)
	filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		path = path[strings.Index(path, pix):]
		path = strings.Replace(path, "\\", "/", -1)
		htmlTemps = append(htmlTemps, path)
		return nil
	})
	return htmlTemps
}
