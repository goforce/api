package soap

import (
	"fmt"
	"reflect"
)

var _ = fmt.Print

func (co *Connection) AddToCache(key string, obj interface{}) {
	if co.cache == nil {
		co.cache = make(map[string]interface{})
	}
	key = reflect.TypeOf(obj).Elem().String() + "/" + key
	co.cache[key] = obj
}

func (co *Connection) GetFromCache(key string, obj interface{}) (interface{}, bool) {
	key = reflect.TypeOf(obj).Elem().String() + "/" + key
	if co.cache != nil {
		if o, ok := co.cache[key]; ok {
			obj = o
			return o, true
		}
	}
	return nil, false
}
