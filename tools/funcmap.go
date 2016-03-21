package tools

import (
	"strings"
	"text/template"
)

var (
	FuncMap = template.FuncMap{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"lowerFirst": func(s string) string {
			v := s[0:1]
			v = strings.ToLower(v) + s[1:]
			return v
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"upperFirst": func(s string) string {
			v := s[0:1]
			v = strings.ToUpper(v) + s[1:]
			return v
		},
		"isBlank": func(s string) bool {
			return IsBlank(s)
		},
		"isNotBlank": func(s string) bool {
			return IsNotBlank(s)
		},
		"fieldType": func(s string) string { //判断字段类型
			switch s {
			case "String":
				return "string"
			case "Integer":
				return "int"
			}
			return ""
		},
		"equal": func(s, e string) bool {
			return s == e
		},
		"gt": func(s, e int) bool {
			return s > e
		},
		"gteq": func(s, e int) bool {
			return s >= e
		},
		"lt": func(s, e int) bool {
			return s < e
		},
		"lteq": func(s, e int) bool {
			return s <= e
		},
	}
)
