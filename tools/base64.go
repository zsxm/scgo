package tools

import (
	"encoding/base64"
	"fmt"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	r, err := coder.DecodeString(string(src))
	return r, err
}

func main() {
	// encode
	hello := `{"openid":"1","location":"1","managerName":"1","managerId":"1","managerTel":"1","managerArea":"中国北京朝阳"}`
	debyte := Base64Encode([]byte(hello))
	debyte = Base64Encode(debyte)
	// decode
	enbyte, _ := Base64Decode(debyte)
	enbyte, _ = Base64Decode(enbyte)

	if hello != string(enbyte) {
		fmt.Println("hello is not equal to enbyte")
	}

	fmt.Println(string(enbyte))
}
