package soap

import (
	//"encoding/xml"
	"errors"
	"fmt"
	//"github.com/goforce/api/soap/core"
	"github.com/goforce/api/commons"
	//"io"
	"math/big"
	"strings"
	"time"
)

// NewRecord returns new instance of record of provided type.
// Optional errs provided to channel errors shortened syntax: r:=NewRecord(co.DescribeSObject(name))
func (co *Connection) NewRecord(name string, values ...map[string]interface{}) (*Record, error) {
	describe, err := co.DescribeSObject(name)
	if err != nil {
		return nil, err
	}
	r := Record{describe: describe, values: make(map[string]interface{})}
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

type Record struct {
	describe *DescribeSObject
	values   map[string]interface{}
}

func (r *Record) SObjectType() string {
	return r.describe.describe.Name
}

func (r *Record) Get(name string) (interface{}, bool) {
	var this, next string
	if i := strings.IndexRune(name, '.'); i == -1 {
		this = name
		next = ""
	} else {
		this = name[:i]
		next = name[i+1:]
	}
	this = strings.ToLower(this)
	if v, ok := r.values[this]; ok {
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

func (r *Record) MustGet(name string) interface{} {
	v, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprint("No such column '", name, "' on entity '", r.SObjectType(), "'"))
	}
	return v
}

func (r *Record) Set(name string, value interface{}) (interface{}, error) {
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
			return value, errors.New(fmt.Sprint("unsupported type of value for field:", name))
		}
		if r.describe != nil {
			if d, ok := r.describe.FieldByName(this); ok {
				converted, err := d.FieldType().ValueOf(value)
				if err != nil {
					return nil, err
				}
				r.values[this] = converted
			} else if d, ok := r.describe.FieldByRelationshipName(this); ok {
				if value == nil {
					r.values[this] = nil
				} else if nr, ok := value.(commons.Record); ok {
					so := strings.ToLower(nr.SObjectType())
					for _, rn := range d.ReferenceTo() {
						if so == strings.ToLower(rn) {
							r.values[this] = value
							return value, nil
						}
					}
					return nil, errors.New(fmt.Sprint("Entity '", nr.SObjectType(), "' not referenced by '", d.Name(), "' allowed entity types: ", d.ReferenceTo()))
				} else {
					return nil, errors.New(fmt.Sprint("value of:", this, " should be record"))
				}
			} else if _, ok := r.describe.ChildRelationshipByName(this); ok {
				if value == nil {
					r.values[this] = &QueryLocator{isLast: true, records: make(commons.Recordset, 0, 0)}
				} else if _, ok := value.(commons.QueryLocator); ok {
					// TODO validate type of query locator allowed?
					r.values[this] = value
				} else {
					return nil, errors.New(fmt.Sprint("value of:", this, " should be query locator", value))
				}
			} else {
				return nil, errors.New(fmt.Sprint("No such column '", this, "' on entity '", r.SObjectType(), "'"))
			}
		} else {
			r.values[this] = value
		}
	} else {
		if nested, ok := r.values[this]; ok {
			if t, ok := nested.(commons.Record); ok {
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

// Fields return list of names set for the Record. Names are same case as in Salesforce.
func (r *Record) Fields() []string {
	names := make([]string, 0, len(r.values))
	for name, _ := range r.values {
		if r.describe != nil {
			if d, ok := r.describe.FieldByName(name); ok {
				names = append(names, d.describe.Name)
			} else if d, ok := r.describe.FieldByRelationshipName(name); ok {
				names = append(names, d.describe.RelationshipName)
			} else if d, ok := r.describe.ChildRelationshipByName(name); ok {
				names = append(names, d.describe.RelationshipName)
			}
		} else {
			names = append(names, name)
		}
	}
	return names
}
