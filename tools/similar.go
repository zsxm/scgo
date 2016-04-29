package tools

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"unicode"
)

//两个字符串相似度比较
func SimilarDegree(str1, str2 string) float64 {
	if len(str2) > len(str1) {
		tmp := str1
		str1 = str2
		str2 = tmp
	}
	nstr1 := removeSign(str1)
	nstr2 := removeSign(str2)
	leng := math.Max(float64(len(nstr1)), float64(len(nstr2)))
	result := longestCommonSubstring(nstr1, nstr2)
	var r = float64(len(result)) * 1.0 / leng
	r, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", r), 64)
	return r
}

func longestCommonSubstring(nstr1, nstr2 string) string {
	var b1 = []byte(nstr1)
	var b2 = []byte(nstr2)
	var m = len(b1)
	var n = len(b2)
	var matrix = make([][]int, m+1)
	for i := 0; i < len(matrix); i++ {
		matrix[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if b1[i-1] == b2[j-1] {
				matrix[i][j] = matrix[i-1][j-1] + 1
			} else {
				matrix[i][j] = int(math.Max(float64(matrix[i][j-1]), float64(matrix[i-1][j])))
			}
		}
	}
	var result = make([]byte, matrix[m][n])
	var currentIndex = len(result) - 1
	for matrix[m][n] != 0 {
		if arrayEq(matrix[n], matrix[n-1]) {
			n--
		} else if matrix[m][n] == matrix[m-1][n] {
			m--
		} else {
			result[currentIndex] = b1[m-1]
			currentIndex--
			m--
			n--
		}
	}
	return string(result)
}

func arrayEq(a1, a2 []int) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func removeSign(str string) string {
	b := bytes.Buffer{}
	for _, v := range str {
		if charReg(v) {
			b.WriteRune(v)
		}
	}
	return b.String()
}

func charReg(b int32) bool {
	if unicode.Is(unicode.Scripts["Han"], b) {
		return true
	} else if b >= 'a' && b <= 'z' {
		return true
	} else if b >= 'A' && b <= 'Z' {
		return true
	} else if b >= '0' && b <= '9' {
		return true
	}
	return false
}
