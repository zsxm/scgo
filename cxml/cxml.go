package cxml

import (
	"encoding/xml"
	"io/ioutil"
)

type CxmlInterface interface{}

func UnmarshalConfig(cxml CxmlInterface, config string) error {
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

func Unmarshal(cxml CxmlInterface, content []byte) error {
	err := xml.Unmarshal(content, &cxml)
	if err != nil {
		return err
	}
	return nil
}

func Marshal(cxml CxmlInterface) ([]byte, error) {

	return xml.MarshalIndent(cxml, "", " ")
}
