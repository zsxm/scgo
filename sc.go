package main

import (
	"fmt"

	"github.com/zsxm/scgo/cjson"
)

func main() {
	var jsn = ` {
    "touser":"OPENID",
    "msgtype":"news",
    "news":{
        "articles": [
         {
             "title":"Happy Day",
             "description":"Is Really A Happy Day",
             "url":"URL",
             "picurl":"PIC_URL"
         },
         {
             "title":"Happy Day",
             "description":"Is Really A Happy Day",
             "url":"URL",
             "picurl":"PIC_URL"
         }
         ]
    }
}`
	cmp := cjson.JsonToMap(jsn)
	fmt.Println(cmp.Get("news").Get("articles").Index(0).Get("description").String())
	mp := make(map[string]string)
	mp["111"] = "asdf"
	js, _ := cjson.MapToJson(mp)
	fmt.Println(string(js))
}
