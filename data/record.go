package data

import (
	"errors"
	"fmt"
	"github.com/goforce/api/soap"
	"math/big"
	"strings"
	"time"
)

type Record interface {
	soap.Record
}

type record struct {
	values soap.Row
	info   *SObjectInfo
}

// NewRecord returns new instance of record of provided type.
func (co *Connection) NewRecord(name string, values ...map[string]interface{}) (*record, error) {
	oi, err := co.DescribeSObject(name)
	if err != nil {
		return nil, err
	}
	return oi.NewRecord(values...)
}

func (oi *SObjectInfo) NewRecord(values ...map[string]interface{}) (*record, error) {
	r := record{info: oi, values: soap.NewRow(oi.dr.Name)}
	for _, m := range values {
		for k, v := range m {
			_, err := r.Set(k, v)
			if err != nil {
				return nil, err
			}
		}
	}
	return &r, nil
}

func (r *record) SObjectType() string {
	if r.info != nil {
		return r.info.dr.Name
	}
	return r.values.SObjectType()
}

func (r *record) Get(name string) (interface{}, bool) {
	return r.values.Get(name)
}

func (r *record) MustGet(name string) interface{} {
	return r.values.MustGet(name)
}

func (r *record) Set(name string, value interface{}) (interface{}, error) {
	this, next := soap.SplitOffFirstName(name)
	if next == "" {
		// single, final name - just set value
		// validate type of value
		switch value.(type) {
		case nil, string, bool, time.Time, *big.Rat, Record, QueryResult:
		case int:
			value = big.NewRat(int64(value.(int)), 1)
		case float64:
			value = new(big.Rat).SetFloat64(value.(float64))
		default:
			return value, errors.New(fmt.Sprint("unsupported type of value for field:", name))
		}
		if r.info != nil {
			if f, ok := r.info.Fields().Get(this); ok {
				converted, err := f.Type().ValueOf(value)
				if err != nil {
					return nil, err
				}
				r.values[this] = converted
			} else if f, ok := r.info.Fields().GetByRelationshipName(this); ok {
				if value == nil {
					r.values[this] = nil
				} else if nr, ok := value.(Record); ok {
					so := strings.ToLower(nr.SObjectType())
					// use internal array to avoid copying
					for _, rn := range f.dr.ReferenceTo {
						if so == strings.ToLower(rn) {
							r.values[this] = value
							return value, nil
						}
					}
					return nil, errors.New(fmt.Sprint("Entity '", nr.SObjectType(), "' not referenced by '", f.Name(), "' allowed entity types: ", f.ReferenceTo()))
				} else {
					return nil, errors.New(fmt.Sprint("value of:", this, " should be record"))
				}
			} else if _, ok := r.info.ChildRelationships().Get(this); ok {
				if value == nil {
					r.values[this] = make(Records, 0, 0)
				} else if _, ok := value.(QueryResult); ok {
					// TODO validate type of query locator allowed?
					r.values[this] = value
				} else {
					return nil, errors.New(fmt.Sprint("value of:", this, " should be query", value))
				}
			} else {
				return nil, errors.New(fmt.Sprint("No such column '", this, "' on entity '", r.SObjectType(), "'"))
			}
		} else {
			r.values[this] = value
		}
	} else {
		if nested, ok := r.values[this]; ok {
			if t, ok := nested.(Record); ok {
				return t.Set(next, value)
			} else {
				return value, errors.New(fmt.Sprint("not a record found for:", this))
			}
		} else {
			return value, errors.New(fmt.Sprint("no record found for:", this))
		}
	}
	return value, nil
}

func (r *record) MustSet(name string, value interface{}) interface{} {
	v, err := r.Set(name, value)
	if err != nil {
		panic(fmt.Sprint("Failed to set '", name, "' on entity '", r.SObjectType(), "' to '", value, "' : ", err))
	}
	return v
}

func (r *record) Remove(name string) (interface{}, bool) {
	return r.values.Remove(name)
}

// Fields return list of names set for the Record. Names are same case as in Salesforce.
func (r *record) Fields() []string {
	names := make([]string, 0, len(r.values)-1)
	for name, _ := range r.values {
		if r.info != nil {
			if d, ok := r.info.Fields().Get(name); ok {
				names = append(names, d.Name())
			} else if d, ok := r.info.Fields().GetByRelationshipName(name); ok {
				names = append(names, d.RelationshipName())
			} else if d, ok := r.info.ChildRelationships().Get(name); ok {
				names = append(names, d.RelationshipName())
			}
		} else {
			names = append(names, name)
		}
	}
	return names
}
