package commons

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
)

func IsBlank(value interface{}) bool {
	if value == nil {
		return true
	}
	if s, ok := value.(string); ok && len(s) == 0 {
		return true
	}
	return false
}

func String(value interface{}) string {
	if value == nil {
		return ""
	}
	switch value.(type) {
	case bool:
		if value.(bool) {
			return "true"
		} else {
			return "false"
		}
	case string:
		return value.(string)
	case *big.Rat:
		return strings.TrimSuffix(strings.TrimRight(value.(*big.Rat).FloatString(20), "0"), ".")
	case time.Time:
		t := value.(time.Time)
		if t.Location().String() == "DATE" {
			return t.Format("2006-01-02")
		} else if t.Location().String() == "TIME" {
			return t.Format("15:04:05.000Z07:00")
		} else {
			return t.Format("2006-01-02T15:04:05.000Z07:00")
		}
	default:
		panic(fmt.Sprint("unsupported data type: ", reflect.TypeOf(value), " / ", value))
	}

}
