package tools

import (
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
