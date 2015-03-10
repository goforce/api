package commons

import (
	"errors"
	"fmt"
	"github.com/goforce/eval"
	"io"
	"strconv"
)

func Must(result interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return result
}

type Connection interface {
	Server() string
	Token() string
	User() string
}

type Recordset []Record

type Record interface {
	SObjectType() string
	Get(field string) (interface{}, bool)
	MustGet(field string) interface{}
	Set(field string, value interface{}) (interface{}, error)
	Fields() []string
}

// QueryLocator is convinience wrapper around Query results.
// Subselect results are always QueryLocator.
// Next and All are supported by core api implementations only if all records
// have been retrieved (using QueryMore...).
type QueryLocator interface {

	// Records returns Recorset in current batch. If there are no records then 0 length slice is returned
	Records() Recordset

	// TotalSize returns number of records expected to be returned by all QueryLocators.
	TotalSize() int

	// IsLast returns true for last QueryLocator.
	IsLast() bool

	// Next returns next QueryLocator. After last QueryLocator error io.EOF is returned.
	Next() (QueryLocator, error)

	// All returns all records returned by query after QueryLocator retrieves all subsequent QueryLocators
	// and appends Records.
	All() (Recordset, error)
}

// QueryReader provides simple iterator over query results using Read.
// To be used as for reader:=NewQueryReader(locator);!reader.EOF();reader.Next() {...}
// Reader exposes Record interface of current record. For reader.EOF() Record is nil.
// QueryReader will panic if Next on QueryLocator will fail.
type QueryReader struct {
	Record
	locator  QueryLocator
	location int
}

func (reader *QueryReader) EOF() bool {
	return reader == nil || reader.locator == nil || reader.Record == nil
}

func (reader *QueryReader) Next() {
	reader.location += 1
	records := reader.locator.Records()
	if reader.location < len(records) {
		reader.Record = records[reader.location]
	} else if reader.locator.IsLast() {
		reader.locator = nil
		reader.Record = nil
	} else {
		locator, err := reader.locator.Next()
		if err == io.EOF {
			reader.locator = nil
			reader.Record = nil
		} else if err != nil {
			panic(err)
		} else {
			reader.locator = locator
			reader.location = 0
			reader.Record = reader.locator.Records()[reader.location]
		}
	}
}

func NewQueryReader(locator QueryLocator) *QueryReader {
	reader := QueryReader{locator: locator, location: 0}
	if locator != nil && len(locator.Records()) > 0 {
		reader.Record = locator.Records()[0]
	}
	return &reader
}

// Grow will create new Recordset increasing size by size. Recordset will be increased only if len == cap.
func (rs Recordset) GrowIfFull() Recordset {
	if len(rs) == cap(rs) {
		cp := make(Recordset, len(rs), cap(rs)+1000)
		copy(cp, rs)
		return cp
	}
	return rs
}

func (rs Recordset) FindAll(expr string, params ...interface{}) (Recordset, error) {
	ce, err := eval.ParseString(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprint(expr, "\n", err))
	}
	result := make(Recordset, 0, 100)
	for _, record := range rs {
		context := eval.NewContext()
		context.AddValues(valueFunc(record, params))
		v, err := ce.Eval(context)
		if err != nil {
			return nil, errors.New(fmt.Sprint("error calculating: ", expr, "\n", err))
		}
		switch v.(type) {
		case bool:
			if v.(bool) {
				result = append(result, record)
			}
		default:
			return nil, errors.New(fmt.Sprint("error, not a boolean:", v, " as result of: ", expr))
		}
	}
	return result, nil
}

func (rs Recordset) MustFindAll(expr string, params ...interface{}) Recordset {
	re, err := rs.FindAll(expr, params...)
	if err != nil {
		panic(err)
	}
	return re
}

func (rs Recordset) FindFirst(expr string, params ...interface{}) (Record, error) {
	ce, err := eval.ParseString(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprint(expr, "\n", err))
	}
	for _, record := range rs {
		context := eval.NewContext()
		context.AddValues(valueFunc(record, params))
		v, err := ce.Eval(context)
		if err != nil {
			return nil, errors.New(fmt.Sprint("error calculating: ", expr, "\n", err))
		}
		switch v.(type) {
		case bool:
			if v.(bool) {
				return record, nil
			}
		default:
			return nil, errors.New(fmt.Sprint("error, not a boolean:", v, " as result of: ", expr))
		}
	}
	return nil, nil
}

func (rs Recordset) MustFindFirst(expr string, params ...interface{}) Record {
	re, err := rs.FindFirst(expr, params...)
	if err != nil {
		panic(err)
	}
	return re
}

func (rs Recordset) FindSingle(expr string, params ...interface{}) (Record, error) {
	ce, err := eval.ParseString(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprint(expr, "\n", err))
	}
	var re Record
	for _, record := range rs {
		context := eval.NewContext()
		context.AddValues(valueFunc(record, params))
		v, err := ce.Eval(context)
		if err != nil {
			return nil, errors.New(fmt.Sprint("error calculating: ", expr, "\n", err))
		}
		switch v.(type) {
		case bool:
			if v.(bool) {
				if re != nil {
					return nil, errors.New(fmt.Sprint("more than one matching record: ", expr, " with ", params))
				}
				re = record
			}
		default:
			return nil, errors.New(fmt.Sprint("error, not a boolean:", v, " as result of: ", expr))
		}
	}
	if re == nil {
		return nil, errors.New(fmt.Sprint("no matching record: ", expr, " with ", params))
	}
	return re, nil
}

func (rs Recordset) MustFindSingle(expr string, params ...interface{}) Record {
	re, err := rs.FindSingle(expr, params...)
	if err != nil {
		panic(err)
	}
	return re
}

func valueFunc(record Record, params []interface{}) eval.Values {
	return func(name string) (interface{}, bool) {
		if value, ok := record.Get(name); ok {
			return value, true
		}
		if n, err := strconv.Atoi(name); err == nil {
			if 0 <= n && n < len(params) {
				return params[n], true
			}
		}
		return nil, false
	}
}
