package data

import (
	"errors"
	"fmt"
	"github.com/goforce/eval"
	"io"
	"strconv"
)

type Records []Record

// Records implement Query interface so that result of the query could be used to set subselect field
func (rs Records) Records() Records    { return rs }
func (rs Records) Reader() Reader      { return &recordsReader{rs: rs} }
func (rs Records) ExpectedLength() int { return len(rs) }

type recordsReader struct {
	rs       Records
	location int
}

func (reader *recordsReader) Read() (Record, error) {
	if reader.location < len(reader.rs) {
		r := reader.rs[reader.location]
		reader.location++
		return r, nil
	} else {
		return nil, io.EOF
	}
}

// Grow will create new Recordset increasing size by size. Recordset will be increased only if len == cap.
func (rs Records) Append(rec Record) Records {
	if len(rs) == cap(rs) {
		cp := make(Records, len(rs), cap(rs)+1000)
		copy(cp, rs)
		return append(cp, rec)
	}
	return append(rs, rec)
}

func (rs Records) FindAll(expr string, params ...interface{}) (Records, error) {
	ce, err := eval.ParseString(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprint(expr, "\n", err))
	}
	result := make(Records, 0, 100)
	for _, rec := range rs {
		context := eval.NewContext()
		context.AddValues(valueFunc(rec, params))
		v, err := ce.Eval(context)
		if err != nil {
			return nil, errors.New(fmt.Sprint("error calculating: ", expr, "\n", err))
		}
		switch v.(type) {
		case bool:
			if v.(bool) {
				result = append(result, rec)
			}
		default:
			return nil, errors.New(fmt.Sprint("error, not a boolean:", v, " as result of: ", expr))
		}
	}
	return result, nil
}

func (rs Records) MustFindAll(expr string, params ...interface{}) Records {
	re, err := rs.FindAll(expr, params...)
	if err != nil {
		panic(err)
	}
	return re
}

func (rs Records) FindFirst(expr string, params ...interface{}) (Record, error) {
	ce, err := eval.ParseString(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprint(expr, "\n", err))
	}
	for _, rec := range rs {
		context := eval.NewContext()
		context.AddValues(valueFunc(rec, params))
		v, err := ce.Eval(context)
		if err != nil {
			return nil, errors.New(fmt.Sprint("error calculating: ", expr, "\n", err))
		}
		switch v.(type) {
		case bool:
			if v.(bool) {
				return rec, nil
			}
		default:
			return nil, errors.New(fmt.Sprint("error, not a boolean:", v, " as result of: ", expr))
		}
	}
	return nil, nil
}

func (rs Records) MustFindFirst(expr string, params ...interface{}) Record {
	re, err := rs.FindFirst(expr, params...)
	if err != nil {
		panic(err)
	}
	return re
}

func (rs Records) FindSingle(expr string, params ...interface{}) (Record, error) {
	ce, err := eval.ParseString(expr)
	if err != nil {
		return nil, errors.New(fmt.Sprint(expr, "\n", err))
	}
	var re Record
	for _, rec := range rs {
		context := eval.NewContext()
		context.AddValues(valueFunc(rec, params))

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
				re = rec
			}
		default:
			return nil, errors.New(fmt.Sprint("error, not a boolean:", v, " as result of: ", expr))
		}
	}
	return re, nil
}

func (rs Records) MustFindSingle(expr string, params ...interface{}) Record {
	re, err := rs.FindSingle(expr, params...)
	if err != nil {
		panic(err)
	}
	if re == nil {
		panic(fmt.Sprint("no matching record: ", expr, " with ", params))
	}
	return re
}

func valueFunc(rec Record, params []interface{}) eval.Values {
	return func(name string) (interface{}, bool) {
		if value, ok := rec.Get(name); ok {
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
