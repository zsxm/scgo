package cjson_test

import (
	"fmt"
	"testing"

	"github.com/zsxm/scgo/cjson"
)

func TestJson(t *testing.T) {
	var jsn = `{
     "errcode" : 0,
     "errmsg" : "ok"
}`
	cmp := cjson.JsonToMap(jsn)
	cmp.Set("name", "123.4555555555555555555555555555553333333333333333333333333333333333333333333333333333335555555555")
	data := map[string]string{
		"name": "1111",
		"pass": "paa",
	}
	cmp.Set("cmp", data)
	fmt.Println("111111", cmp.Get("name").String())
	//	js, _ := cmp.MapToJson(mp)
	//	fmt.Println(string(js))
}
