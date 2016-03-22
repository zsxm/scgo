package cxml

import (
	"encoding/xml"
	"io/ioutil"
)

type CxmlInterface interface{}

func Unmarshal(cxml CxmlInterface, config string) error {
	content, err := ioutil.ReadFile(config)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(content, &cxml)
	if err != nil {
		return err
	}
	return nil
}
