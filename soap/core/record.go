package core

import (
	"errors"
	"fmt"
	"github.com/goforce/api/commons"
	"math/big"
	"strings"
	"time"
)

// sobjecttype is used to store object type in same record map. All other names/keys should be lower case.
const sobjecttype = "SObjectType"

type Record map[string]interface{}

func NewRecord(sObjectType string, values ...map[string]interface{}) commons.Record {
	re := make(Record)
	for _, mv := range values {
		for k, v := range mv {
			re[k] = v
		}
	}
	re["SObjectType"] = sObjectType
	return re
}

func (r Record) Get(name string) (interface{}, bool) {
	var this, next string
	if i := strings.IndexRune(name, '.'); i == -1 {
		this = name
		next = ""
	} else {
		this = name[:i]
		next = name[i+1:]
	}
	this = strings.ToLower(this)
	if v, ok := r[this]; ok {
		if next == "" {
			return v, true
		} else if nested, ok := v.(commons.Record); ok {
			return nested.Get(next)
		} else if v == nil {
			return nil, true
		}
	}
	return nil, false
}

func (r Record) MustGet(name string) interface{} {
	v, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprint("No such column '", name, "' on entity '", r.SObjectType(), "'"))
	}
	return v
}

// Set value of the field name. Returns same value and error if value contains unsupported type.
// Supported value types are: string, bool, *time.Time, *big.Rat, commons.Record, commons.QueryLocator.
// Set does not support setting values to nested subselects.
func (r Record) Set(name string, value interface{}) (interface{}, error) {
	// go lower and split off first name
	var this, next string
	if i := strings.IndexRune(name, '.'); i == -1 {
		this = name
		next = ""
	} else {
		this = name[:i]
		next = name[i+1:]
	}
	this = strings.ToLower(this)
	if next == "" {
		// single, final name - just set value
		// validate type of value
		switch value.(type) {
		case nil, string, bool, time.Time, *big.Rat, commons.Record, commons.QueryLocator:
		default:
			return value, errors.New(fmt.Sprint("unsupported type of value for field: ", name, "=", value))
		}
		r[this] = value
	} else {
		// if nested name
		if v, ok := r[this]; ok {
			if nested, ok := v.(commons.Record); ok {
				return nested.Set(next, value)
			} else {
				return value, errors.New(fmt.Sprint("not a row:", this))
			}
		} else {
			return value, errors.New(fmt.Sprint("no row created:", this))
		}
	}
	return value, nil
}

func (r Record) Fields() []string {
	fields := make([]string, 0, len(r))
	for k, _ := range r {
		if k != sobjecttype {
			fields = append(fields, k)
		}
	}
	return fields
}

func (r Record) SObjectType() string {
	s, _ := r[sobjecttype].(string)
	return s
}
