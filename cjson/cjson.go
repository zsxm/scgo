package cjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zsxm/scgo/log"
)

type JSON struct {
	data interface{}
}

func (this *JSON) Set(key string, value interface{}) {
	if this.data == nil {
		this.data = make(map[string]interface{})
	}
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
	} else if m, ok := this.data.([]interface{}); ok {
		var b bytes.Buffer
		b.WriteString("[")
		for i, v := range m {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(v.(string))
		}
		b.WriteString("]")
		return b.String()
	}
	return ""
}

func (this *JSON) Integer() int {
	if m, ok := this.data.(string); ok {
		if m != "" {
			r, _ := strconv.Atoi(m)

			return r
		}
		return -1
	} else if m, ok := this.data.(float64); ok {
		r, _ := strconv.Atoi(strconv.FormatFloat(m, 'f', -1, 64))
		return r
	}
	return -1
}

func (this *JSON) Integer64() int64 {
	if m, ok := this.data.(string); ok {
		if m != "" {
			r, _ := strconv.ParseInt(m, -1, 64)
			return r
		}
		return -1
	}
	return -1
}
func (this *JSON) Float() float64 {
	if m, ok := this.data.(string); ok {
		if m != "" {
			v, _ := strconv.ParseFloat(m, 64)
			return v
		}
		return -1
	} else if m, ok := this.data.(float64); ok {
		return m
	}
	return -1
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

func (this *JSON) DataMaps() map[string][]string {
	mp := this.Data()
	if mp != nil {
		result := make(map[string][]string)
		for k, v := range mp {
			if m, ok := v.(string); ok {
				result[k] = []string{m}
			} else if m, ok := v.(float64); ok {
				result[k] = []string{strconv.FormatFloat(m, 'f', -1, 64)}
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
