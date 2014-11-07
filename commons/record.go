package commons

import (
	"fmt"
	"math/big"
	"time"
)

var _ = fmt.Print

type Record interface {
	SObjectType() string
	Get(name string) (interface{}, bool)
	Set(name string, value interface{}) (interface{}, error)
	Fields() []string
}

func Must(value interface{}, ok bool) interface{} {
	if !ok {
		panic("no such field")
	}
	return value
}

func MustQueryLocator(value interface{}, ok bool) QueryLocator {
	if !ok {
		panic("no such field")
	}
	if v, ok := value.(QueryLocator); ok {
		return v
	}
	panic("not a query locator")
}

func MustRecord(value interface{}, ok bool) Record {
	if !ok {
		panic("no such field")
	}
	if v, ok := value.(Record); ok {
		return v
	}
	panic("not a record")
}

func MustString(value interface{}, ok bool) string {
	if !ok {
		panic("no such field")
	}
	if value == nil {
		return ""
	}
	if v, ok := value.(string); ok {
		return v
	}
	panic("not a string")
}

func MustNumber(value interface{}, ok bool) *big.Rat {
	if !ok {
		panic("no such field")
	}
	if v, ok := value.(*big.Rat); ok {
		return v
	}
	panic("not a number")
}

func MustDate(value interface{}, ok bool) time.Time {
	if !ok {
		panic("no such field")
	}
	if v, ok := value.(time.Time); ok {
		return v
	}
	panic("not a date")
}

func MustBoolean(value interface{}, ok bool) bool {
	if !ok {
		panic("no such field")
	}
	if v, ok := value.(bool); ok {
		return v
	}
	panic("not a boolean")
}
