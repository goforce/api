// Package conv provides data types and conversion functions for salesforce field values.
// Data types supported: nil, string, bool, *big.Rat, time.Time.
// time.Time distinguishes date, time and datetime using time.Location. For date DATE_LOCATION should be used
// and for time TIME_LOCATION should be used.
package conv

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// *time.Location to distinguish TIME data type
var LOCATION_TIME *time.Location = time.FixedZone("TIME", 0)

// *time.Location to distinguish DATE data type
var LOCATION_DATE *time.Location = time.FixedZone("DATE", 0)

const (
	FORMAT_DATE     string = "2006-01-02"
	FORMAT_TIME     string = "15:04:05.000Z07:00"
	FORMAT_DATETIME string = "2006-01-02T15:04:05.000Z07:00"
)

type Location struct {
	Latitude  *big.Rat
	Longitude *big.Rat
}

type Address struct {
	Location
	City        string
	Country     string
	CountryCode string
	PostalCode  string
	State       string
	StateCode   string
	Street      string
}

// IsBlank returns true if value is nil or emtpy string
func IsBlank(value interface{}) bool {
	s, ok := value.(string)
	return value == nil || (ok && s == "")
}

// String converts string, bool, *big.Rat, time.Time (DATETIME, DATE, TIME) to salesforce api string representation.
// nil is converted to empty string.
func String(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}
	switch value.(type) {
	case bool:
		if value.(bool) {
			return "true", nil
		} else {
			return "false", nil
		}
	case string:
		return value.(string), nil
	case *big.Rat:
		return strings.TrimSuffix(strings.TrimRight(value.(*big.Rat).FloatString(20), "0"), "."), nil
	case time.Time:
		t := value.(time.Time)
		if t.Location().String() == "DATE" {
			return t.Format(FORMAT_DATE), nil
		} else if t.Location().String() == "TIME" {
			return t.Format(FORMAT_TIME), nil
		} else {
			return t.Format(FORMAT_DATETIME), nil
		}
	}
	return "", errors.New(fmt.Sprint("cannot convert to a string: ", value))
}

// Type is set of constants of supported types with conversion functions.
type Type int

const (
	UNKNOWN Type = iota
	STRING
	BOOLEAN
	NUMBER
	DATE
	TIME
	DATETIME
	LOCATION
	ADDRESS
)

// String representation of value type (for use in error messages)
func (t Type) String() string {
	switch t {
	case UNKNOWN:
		return "unknown"
	case STRING:
		return "string"
	case BOOLEAN:
		return "boolean"
	case NUMBER:
		return "number"
	case DATE:
		return "date"
	case TIME:
		return "time"
	case DATETIME:
		return "datetime"
	case LOCATION:
		return "location"
	case ADDRESS:
		return "address"
	}
	return "unsupported"
}

// Parse string value from salesforce api representation to specified value type
func (t Type) Parse(value string) (interface{}, error) {
	if len(value) == 0 {
		return nil, nil
	}
	switch t {
	case STRING:
		return value, nil
	case BOOLEAN:
		if value == "true" {
			return true, nil
		} else if value == "false" {
			return false, nil
		} else {
			return nil, errors.New(fmt.Sprint("not a boolean: ", value))
		}
	case NUMBER:
		if r, ok := new(big.Rat).SetString(value); ok {
			return r, nil
		} else {
			return nil, errors.New(fmt.Sprint("not a number: ", value))
		}
	case DATE:
		if d, err := time.ParseInLocation(FORMAT_DATE, value, LOCATION_DATE); err == nil {
			return d, nil
		} else {
			return nil, errors.New(fmt.Sprint("not a date: ", value))
		}
	case TIME:
		if t, err := time.Parse(FORMAT_TIME, value); err == nil {
			return time.Date(1, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), LOCATION_TIME), nil
		} else {
			return nil, errors.New(fmt.Sprint("not a time: ", value))
		}
	case DATETIME:
		if dt, err := time.Parse(FORMAT_DATETIME, value); err == nil {
			return dt, nil
		} else {
			return nil, errors.New(fmt.Sprint("not a datetime: ", value))
		}
	}
	return nil, errors.New(fmt.Sprint("illegal value `", value, "` for type ", t.String()))
}

// ValueOf converts value to specified value type
func (t Type) ValueOf(value interface{}) (result interface{}, err error) {
	switch t {
	case STRING:
		if _, ok := value.(string); ok || value == nil {
			return value, nil
		} else {
			return String(value)
		}
	case BOOLEAN:
		if _, ok := value.(bool); ok || value == nil {
			return value, nil
		} else if s, ok := value.(string); ok {
			if len(s) == 0 {
				return nil, nil
			} else {
				return BOOLEAN.Parse(s)
			}
		} else {
			return nil, errors.New(fmt.Sprint("cannot convert to a boolean: ", value))
		}
	case NUMBER:
		if _, ok := value.(*big.Rat); ok || value == nil {
			return value, nil
		} else if s, ok := value.(string); ok {
			return NUMBER.Parse(s)
		} else if f, ok := value.(float64); ok {
			return new(big.Rat).SetFloat64(f), nil
		} else if i, ok := value.(int); ok {
			return big.NewRat(int64(i), 1), nil
		} else {
			return nil, errors.New(fmt.Sprint("cannot convert to a number: ", value))
		}
	case DATE:
		if value == nil {
			return nil, nil
		} else if t, ok := value.(time.Time); ok {
			return t.Truncate(time.Hour * 24), nil
		} else if s, ok := value.(string); ok {
			return DATE.Parse(s)
		} else {
			return nil, errors.New(fmt.Sprint("cannot convert to a date: ", value))
		}
	case TIME:
		if value == nil {
			return nil, nil
		} else if t, ok := value.(time.Time); ok {
			return time.Date(1, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), LOCATION_TIME), nil
		} else if s, ok := value.(string); ok {
			return TIME.Parse(s)
		} else {
			return nil, errors.New(fmt.Sprint("cannot convert to a time:", value))
		}
	case DATETIME:
		if _, ok := value.(time.Time); ok || value == nil {
			return value, nil
		} else if s, ok := value.(string); ok {
			return DATETIME.Parse(s)
		} else {
			return nil, errors.New(fmt.Sprint("cannot convert to a datetime: ", value))
		}
	case LOCATION:
		if _, ok := value.(Location); ok || value == nil {
			return value, nil
		} else {
			return nil, errors.New(fmt.Sprint("not a location: ", value))
		}
	case ADDRESS:
		if _, ok := value.(Address); ok || value == nil {
			return value, nil
		} else {
			return nil, errors.New(fmt.Sprint("not an address: ", value))
		}
	}
	return nil, errors.New(fmt.Sprint("unsupported type"))
}
