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
	fmt.Println("111111", cmp)
	//	js, _ := cmp.MapToJson(mp)
	//	fmt.Println(string(js))
}
