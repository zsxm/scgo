package cjson

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zsxm/scgo/log"
)

type JSON struct {
	data interface{}
}

func (this *JSON) Set(key, value string) {
	if v, ok := this.data.(map[string]interface{}); ok {
		v[key] = value
	}
}

func (this *JSON) Get(key string) *JSON {
	m := this.Data()
	if v, ok := m[key]; ok {
		tmpJSON := &JSON{v}
		return tmpJSON
	}
	return this
}

func (this *JSON) String() string {
	if m, ok := this.data.(string); ok {
		return m
	} else if m, ok := this.data.(float64); ok {
		return strconv.FormatFloat(m, 'f', -1, 64)
	}
	return ""
}

func (this *JSON) Size() int {
	if v, ok := this.data.([]interface{}); ok {
		return len(v)
	}
	if v, ok := this.data.(map[string]interface{}); ok {
		return len(v)
	}
	return 0
}

func (this *JSON) Index(i int) *JSON {
	if v, ok := this.data.([]interface{}); ok {
		tmpIfe := v[i]
		if v1, ok1 := tmpIfe.(map[string]interface{}); ok1 {
			tmpJSON := &JSON{v1}
			return tmpJSON
		}
	}
	return this
}

func (this *JSON) Data() map[string]interface{} {
	if m, ok := (this.data).(map[string]interface{}); ok {
		return m
	}
	return nil
}

func (this *JSON) DataMap() map[string]string {
	mp := this.Data()
	if mp != nil {
		result := make(map[string]string)
		for k, v := range mp {
			if m, ok := v.(string); ok {
				result[k] = m
			} else if m, ok := v.(float64); ok {
				result[k] = strconv.FormatFloat(m, 'f', -1, 64)
			}
		}
		return result
	}
	return nil
}

func JsonToMap(data string) *JSON {
	js := &JSON{}
	d := new(interface{})
	err := json.Unmarshal([]byte(data), d)
	if err != nil {
		log.Error(err)
		fmt.Println("json str to map", err)
	}
	js.data = *d
	return js
}

func MapToJson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
