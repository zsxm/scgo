package tools

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zsxm/scgo/tools/uuid"
)

const (
	//生成sql用到的主键标记
	GEN_SQL_IDENTIFY = "_identify"
)

//判断空
func IsBlank(s string) bool {
	if strings.Trim(s, " ") == "" || s == "nil" || s == "null" {
		return true
	}
	return false
}

//判断不是空
func IsNotBlank(s string) bool {
	if strings.TrimSpace(s) != "" && s != "nil" && s != "null" {
		return true
	}
	return false
}

//转换为int
func ParseInteger(v string) int {
	s, _ := strconv.Atoi(v)
	return s
}

//判断文件是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//获得当前路径
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

//转换小写
func Lower(v string) string {
	return strings.ToLower(v)
}

//转换大写
func Upper(v string) string {
	return strings.ToUpper(v)
}

//休眠
func Sleep(sleep int) {
	time.Sleep(time.Duration(sleep) * time.Second)
}

//两个数值比较 t1>t2
func Compare(t1, t2 int) bool {
	return t1 > t2
}

//map转sql insert
//当data map某个value值为identify时，标记它为主键标识生成一个id
//data数据 table表名
func MapToSQL_Insert(data map[string]string, table string) (string, []interface{}) {
	var sql bytes.Buffer
	var cloms bytes.Buffer
	var vals bytes.Buffer
	args := make([]interface{}, 0, len(data))
	sql.WriteString("insert into ")
	sql.WriteString(table)
	sql.WriteString(" (")
	i := 0
	for k, v := range data {
		if i > 0 {
			cloms.WriteString(",")
			vals.WriteString(",")
		}
		cloms.WriteString(k)
		vals.WriteString("?")
		if v == GEN_SQL_IDENTIFY {
			v = uuid.NewV1().String()
		}
		args = append(args, v)
		i++
	}
	sql.WriteString(cloms.String())
	sql.WriteString(") ")
	sql.WriteString("values(")
	sql.WriteString(vals.String())
	sql.WriteString(")")
	return sql.String(), args
}
