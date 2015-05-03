package soap

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Record interface {
	SObjectType() string
	Get(field string) (interface{}, bool)
	MustGet(field string) interface{}
	Set(field string, value interface{}) (interface{}, error)
	MustSet(field string, value interface{}) interface{}
	Remove(field string) (interface{}, bool)
	Fields() []string
}

const sobjecttype = "SObjectType"

type Row map[string]interface{}

func NewRow(sObjectType string) Row {
	r := make(Row)
	r[sobjecttype] = sObjectType
	return r
}

func NewRecord(sObjectType string, values ...map[string]interface{}) Record {
	r := NewRow(sObjectType)
	for _, m := range values {
		for k, v := range m {
			r[strings.ToLower(k)] = v
		}
	}
	return r
}

func (r Row) SObjectType() string {
	return r[sobjecttype].(string)
}

func (r Row) Get(name string) (interface{}, bool) {
	this, next := SplitOffFirstName(name)
	if v, ok := r[this]; ok {
		if v == nil || next == "" {
			return v, true
		} else if nested, ok := v.(Record); ok {
			return nested.Get(next)
		}
	}
	return nil, false
}

func (r Row) MustGet(name string) interface{} {
	v, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprint("No such column '", name, "' on entity '", r.SObjectType(), "'"))
	}
	return v
}

// Set value of the field name. Returns same value and error if value contains unsupported type.
// Supported value types are: string, bool, *time.Time, *big.Rat, Record, *QueryResult.
// int and float64 are converted to *big.Rat for convinience.
// Set does not support setting values to nested subselects.
func (r Row) Set(name string, value interface{}) (interface{}, error) {
	this, next := SplitOffFirstName(name)
	if next == "" {
		// single, final name - just set value
		// validate type of value
		switch value.(type) {
		case nil, string, bool, time.Time, *big.Rat, Record, *QueryResult:
		case int:
			value = big.NewRat(int64(value.(int)), 1)
		case float64:
			value = new(big.Rat).SetFloat64(value.(float64))
		default:
			return value, errors.New(fmt.Sprint("unsupported type of value for field: ", name, "=", value))
		}
		r[this] = value
	} else {
		// if nested name
		if v, ok := r[this]; ok {
			if nested, ok := v.(Record); ok {
				return nested.Set(next, value)
			}
		}
		return value, errors.New(fmt.Sprint("not a record:", this))
	}
	return value, nil
}

func (r Row) MustSet(name string, value interface{}) interface{} {
	v, err := r.Set(name, value)
	if err != nil {
		panic(fmt.Sprint("Failed to set '", name, "' on entity '", r.SObjectType(), "' to '", value, "' : ", err))
	}
	return v
}

func (r Row) Remove(name string) (interface{}, bool) {
	this, next := SplitOffFirstName(name)
	if v, ok := r[this]; ok {
		if next == "" {
			delete(r, this)
			return v, true
		} else if nested, ok := v.(Record); ok {
			return nested.Remove(next)
		} else if v == nil {
			return nil, true
		}
	}
	return nil, false
}

func (r Row) Fields() []string {
	fields := make([]string, 0, len(r)-1)
	for k, _ := range r {
		if k != sobjecttype {
			fields = append(fields, k)
		}
	}
	return fields
}

func SplitOffFirstName(name string) (this, next string) {
	if i := strings.IndexRune(name, '.'); i == -1 {
		this = name
		next = ""
	} else {
		this = name[:i]
		next = name[i+1:]
	}
	this = strings.ToLower(this)
	return
}
