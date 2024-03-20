package xmap

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type CMap map[string][]string

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v CMap) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Pop return Get value and Del key
func (v CMap) Pop(key string) string {
	s := v.Get(key)
	v.Del(key)
	return s
}

// Set sets the key to value. It replaces any existing
// values.
func (v CMap) Set(key, value string) CMap {
	v[key] = []string{value}
	return v
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v CMap) Add(key, value string) CMap {
	v[key] = append(v[key], value)
	return v
}

// Del deletes the values associated with key.
func (v CMap) Del(key string) CMap {
	delete(v, key)
	return v
}

// Obtain get all values associated with the given key.
func (v CMap) Obtain(key string) []string {
	if v == nil {
		return []string{}
	}
	return v[key]
}

// Append set the key to value if it doesn't exists. append if it exists.
func (v CMap) Append(key, value string) CMap {
	vs := v.Get(key)
	if vs == "" || len(strings.Trim(vs, " ")) == 0 {
		v.Set(key, value)
		return v
	}
	return v.Set(key, vs+value)
}

// Insert Append set the key to value if it doesn't exists. insert if it exists.
func (v CMap) Insert(key, value string) CMap {
	vs := v.Get(key)
	if vs == "" || len(strings.Trim(vs, " ")) == 0 {
		v.Set(key, value)
		return v
	}
	return v.Set(key, value+vs)
}

// Struct CMap to struct data
// value like
// type Te struct {
//		Name string `json:"name"` // must string type
//		ID   string `json:"id"` // must string type
//	}
func (v CMap) Struct(value interface{}) error {
	var m = map[string]interface{}{}
	for k, v := range v {
		m[k] = v[0]
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, value)
	if err != nil {
		return err
	}
	return nil
}

// StructToMap struct to CMap, maybe use Encode
func StructToMap(v interface{}) (values CMap) {
	values = NewCMap()
	el := reflect.ValueOf(v)
	if el.Kind() == reflect.Ptr {
		el = el.Elem()
	}
	iVal := el
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		fi := typ.Field(i)
		name := fi.Tag.Get("json")
		if name == "" {
			name = fi.Name
		}
		values.Set(name, fmt.Sprint(iVal.Field(i)))
	}
	return
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func (v CMap) Encode() string {
	return url.Values(v).Encode()
}

// CMap url.Values to CMap
// Deprecated
func (v CMap) CMap(values url.Values) CMap {
	return CMap(values)
}

// NewCMap new cmap
func NewCMap() CMap {
	return CMap{}
}
