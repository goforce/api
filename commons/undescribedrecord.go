package commons

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type undescribedRecordField struct {
	value        interface{}
	originalName string
}

type undescribedRecord struct {
	sObjectType string
	fields      map[string]undescribedRecordField
}

func NewUndescribedRecord(sObjectType string) Record {
	return &undescribedRecord{sObjectType: sObjectType, fields: make(map[string]undescribedRecordField)}
}

func (rec *undescribedRecord) SObjectType() string {
	return rec.sObjectType
}

func (rec *undescribedRecord) Get(name string) (interface{}, bool) {
	var this, next string
	if i := strings.IndexRune(name, '.'); i == -1 {
		this = name
		next = ""
	} else {
		this = name[:i]
		next = name[i+1:]
	}
	this = strings.ToLower(this)
	if f, ok := rec.fields[this]; ok {
		if next == "" {
			return f.value, true
		} else if nestedRecord, ok := f.value.(Record); ok {
			return nestedRecord.Get(next)
		} else if f.value == nil {
			return nil, true
		}
	}
	return nil, false
}

// Set value of the field name. Returns same value and error if value contains unsupported type.
// Supported value types are: string, bool, time.Time, big.Rat, Record, QueryLocator.
// Set does not support setting values to nested subselects.
func (rec *undescribedRecord) Set(name string, value interface{}) (interface{}, error) {
	// validate type of value
	switch value.(type) {
	case nil, string, bool, time.Time, big.Rat, Record, QueryLocator:
	default:
		return value, errors.New(fmt.Sprint("unsupported type of value for field: ", name, " = ", value))
	}
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
		if pv, ok := rec.fields[this]; ok {
			rec.fields[this] = undescribedRecordField{value: value, originalName: pv.originalName}
		} else {
			rec.fields[this] = undescribedRecordField{value: value, originalName: name}
		}
	} else {
		// if nested name
		if nestedField, ok := rec.fields[this]; ok {
			if nestedRecord, ok := nestedField.value.(Record); ok {
				return nestedRecord.Set(next, value)
			} else {
				return value, errors.New(fmt.Sprint("not a record found for:", this))
			}
		} else {
			return value, errors.New(fmt.Sprint("no record found for:", this))
		}
	}
	return value, nil
}

func (rec *undescribedRecord) Fields() []string {
	result := make([]string, 0, len(rec.fields))
	for _, f := range rec.fields {
		result = append(result, f.originalName)
	}
	return result
}
