package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

//计算当前时间和t时间在second秒之内
//t 比较的时间，second 符合多少秒内
func TimePoor(t string, second int) bool {
	t1, _ := strconv.Atoi(t)
	t2 := int(time.Now().Unix())
	t3 := t2 - t1
	fmt.Println(t2, t1, t2-t1)
	if t3 < second {
		return true
	}
	return false
}

func Lower(v string) string {
	return strings.ToLower(v)
}

func Upper(v string) string {
	return strings.ToUpper(v)
}
