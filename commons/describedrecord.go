package commons

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// timezone called DATE is used to distinguish DATE fields
var dateLocation = time.FixedZone("DATE", 0)

type describedField struct {
	value         interface{}
	describe      *DescribeSObjectFieldResult
	childDescribe *ChildRelationship
}

type describedRecord struct {
	describe *DescribeSObjectResult
	fields   map[string]describedField
}

// NewDescribedRecord returns new instance of record of provided type.
// Optional errs provided to channel errors shortened syntax: r:=NewDescribedRecord(DescribeSObject(connection,name))
func NewDescribedRecord(describe *DescribeSObjectResult, errs ...error) (Record, error) {
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	if describe == nil {
		return nil, errors.New("empty describe info.")
	}
	return &describedRecord{describe: describe, fields: make(map[string]describedField)}, nil
}

func (rec *describedRecord) SObjectType() string {
	return rec.describe.Name
}

func (rec *describedRecord) Get(name string) (interface{}, bool) {
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

func (rec *describedRecord) Set(name string, value interface{}) (interface{}, error) {
	// validate type of value
	//	switch value.(type) {
	//	case string, bool, time.Time, *big.Rat, Record, QueryLocator, nil:
	//	default:
	//		return value, errors.New(fmt.Sprint("unsupported type of value for field:", name))
	//	}
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
		// single, final name - validate name, value and set
		if fd := rec.describe.Get(this); fd != nil {
			switch fd.Type {
			case "string", "picklist", "multipicklist", "combobox", "reference", "base64", "textarea", "phone", "id", "url", "email", "encryptedstring", "datacategorygroupreference":
				if _, ok := value.(string); ok || value == nil {
					rec.fields[this] = describedField{value: value, describe: fd}
					return value, nil
				} else {
					return nil, errors.New(fmt.Sprint("value passed for:", this, ": not a string"))
				}
			case "boolean":
				if _, ok := value.(bool); ok || value == nil {
					rec.fields[this] = describedField{value: value, describe: fd}
					return value, nil
				} else if s, ok := value.(string); ok {
					if s == "true" {
						rec.fields[this] = describedField{value: true, describe: fd}
						return value, nil
					} else if s == "false" {
						rec.fields[this] = describedField{value: false, describe: fd}
						return value, nil
					} else {
						return nil, errors.New(fmt.Sprint("invalid value for:", this, " not a true or false:", value))
					}
				} else {
					return nil, errors.New(fmt.Sprint("value passed for:", this, ": not a boolean"))
				}
			case "currency", "int", "double", "percent":
				if _, ok := value.(*big.Rat); ok || value == nil {
					rec.fields[this] = describedField{value: value, describe: fd}
					return value, nil
				} else if s, ok := value.(string); ok {
					if r, ok := new(big.Rat).SetString(s); ok {
						rec.fields[this] = describedField{value: r, describe: fd}
						return value, nil
					} else {
						return nil, errors.New(fmt.Sprint("failed to convert to number for:", this, " value:", s))
					}
				} else {
					return nil, errors.New(fmt.Sprint("value passed for:", this, ": not a number"))
				}
			case "date":
				if IsBlank(value) {
					rec.fields[this] = describedField{value: nil, describe: fd}
					return value, nil
				} else if t, ok := value.(time.Time); ok {
					rec.fields[this] = describedField{value: t.Truncate(time.Hour * 24), describe: fd}
					return value, nil
				} else if s, ok := value.(string); ok {
					if d, err := time.ParseInLocation("2006-01-02", s, dateLocation); err == nil {
						rec.fields[this] = describedField{value: d, describe: fd}
						return value, nil
					} else {
						return nil, errors.New(fmt.Sprint("failed to convert to date for:", this, " value:", s))
					}
				} else {
					return nil, errors.New(fmt.Sprint("value passed for:", this, ": not a date:", value))
				}
			case "datetime", "time":
				if IsBlank(value) {
					rec.fields[this] = describedField{value: nil, describe: fd}
					return value, nil
				} else if _, ok := value.(time.Time); ok {
					rec.fields[this] = describedField{value: value, describe: fd}
					return value, nil
				} else if s, ok := value.(string); ok {
					if d, err := time.Parse("2006-01-02T15:04:05.999Z0700", s); err == nil {
						rec.fields[this] = describedField{value: d, describe: fd}
						return value, nil
					} else {
						return nil, errors.New(fmt.Sprint("failed to convert to datetime for:", this, " value:", s))
					}
				} else {
					return nil, errors.New(fmt.Sprint("value passed for:", this, ": not a datetime:", value))
				}
			case "location", "address", "anyType":
				return nil, errors.New(fmt.Sprint("unsupported type of field:", this))
			default:
				return nil, errors.New(fmt.Sprint("unknown type of field:", this))
			}
		} else if fd := rec.describe.GetRelationship(this); fd != nil {
			if value == nil {
				rec.fields[this] = describedField{value: value, describe: fd}
				return value, nil
			} else if rr, ok := value.(Record); ok {
				so := rr.SObjectType()
				for _, rn := range fd.ReferenceTo {
					if so == rn {
						rec.fields[this] = describedField{value: rr, describe: fd}
						return value, nil
					}
				}
				return nil, errors.New(fmt.Sprint("record of ", so, " not allowed for value of:", this))
			} else {
				return nil, errors.New(fmt.Sprint("value of:", this, " should be record"))
			}
		} else if cr := rec.describe.GetChildRelationship(this); cr != nil {
			if value == nil {
				rec.fields[this] = describedField{value: &EmptyQueryLocator{}, childDescribe: cr}
				return value, nil
			} else if _, ok := value.(QueryLocator); ok {
				// TODO validate type of query locator allowed?
				rec.fields[this] = describedField{value: value, childDescribe: cr}
				return value, nil
			} else {
				return nil, errors.New(fmt.Sprint("value of:", this, " should be query locator", value))
			}
		}
		return nil, errors.New(fmt.Sprint("no field named:", this))
	} else {
		if nestedRec, ok := rec.fields[this]; ok {
			if t, ok := nestedRec.value.(Record); ok {
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
func (rec *describedRecord) Fields() []string {
	result := make([]string, 0, len(rec.fields))
	for name, _ := range rec.fields {
		if describe := rec.describe.Get(name); describe != nil {
			result = append(result, describe.Name)
		} else if describe := rec.describe.GetRelationship(name); describe != nil {
			result = append(result, describe.RelationshipName)
		} else if describe := rec.describe.GetChildRelationship(name); describe != nil {
			result = append(result, describe.RelationshipName)
		}
	}
	return result
}
