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
	cmp.Set("name", "123.4")
	data := map[string]string{
		"name": "1111",
		"pass": "paa",
	}
	cmp.Set("cmp", data)
	fmt.Println("111111", cmp.Data())
	//	js, _ := cmp.MapToJson(mp)
	//	fmt.Println(string(js))
}
