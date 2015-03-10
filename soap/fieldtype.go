package soap

import (
	"errors"
	"fmt"
	"github.com/goforce/api/commons"
)

type FieldType int

const (
	UNKNOWN FieldType = iota + 1
	STRING
	PICKLIST
	MULTIPICKLIST
	COMBOBOX
	REFERENCE
	BASE64
	BOOLEAN
	CURRENCY
	TEXTAREA
	INT
	DOUBLE
	PERCENT
	PHONE
	ID
	DATE
	DATETIME
	TIME
	URL
	EMAIL
	ENCRYPTEDSTRING
	DATACATEGORYGROUPREFERENCE
	LOCATION
	ADDRESS
	ANYTYPE
)

func getFieldType(fieldType string) FieldType {
	switch fieldType {
	case "string":
		return STRING
	case "picklist":
		return PICKLIST
	case "multipicklist":
		return MULTIPICKLIST
	case "combobox":
		return COMBOBOX
	case "reference":
		return REFERENCE
	case "base64":
		return BASE64
	case "boolean":
		return BOOLEAN
	case "currency":
		return CURRENCY
	case "textarea":
		return TEXTAREA
	case "int":
		return INT
	case "double":
		return DOUBLE
	case "percent":
		return PERCENT
	case "phone":
		return PHONE
	case "id":
		return ID
	case "date":
		return DATE
	case "datetime":
		return DATETIME
	case "time":
		return TIME
	case "url":
		return URL
	case "email":
		return EMAIL
	case "encryptedstring":
		return ENCRYPTEDSTRING
	case "datacategorygroupreference":
		return DATACATEGORYGROUPREFERENCE
	case "location":
		return LOCATION
	case "address":
		return ADDRESS
	case "anyType":
		return ANYTYPE
	}
	return UNKNOWN
}

// Value validates and converts value to one which would match field type
func (t FieldType) Parse(value string) (interface{}, error) {
	switch t {
	case STRING, PICKLIST, MULTIPICKLIST, COMBOBOX, REFERENCE, BASE64, TEXTAREA, PHONE, ID, URL, EMAIL, ENCRYPTEDSTRING, DATACATEGORYGROUPREFERENCE:
		return commons.STRING.Parse(value)
	case BOOLEAN:
		return commons.BOOLEAN.Parse(value)
	case CURRENCY, INT, DOUBLE, PERCENT:
		return commons.NUMBER.Parse(value)
	case DATE:
		return commons.DATE.Parse(value)
	case TIME:
		return commons.TIME.Parse(value)
	case DATETIME:
		return commons.DATETIME.Parse(value)
	case LOCATION, ADDRESS, ANYTYPE:
		return nil, errors.New(fmt.Sprint("unsupported type"))
	}
	return nil, errors.New(fmt.Sprint("unknown type"))
}

// Value validates and converts value to one which would match field type
func (t FieldType) ValueOf(value interface{}) (interface{}, error) {
	switch t {
	case STRING, PICKLIST, MULTIPICKLIST, COMBOBOX, REFERENCE, BASE64, TEXTAREA, PHONE, ID, URL, EMAIL, ENCRYPTEDSTRING, DATACATEGORYGROUPREFERENCE:
		return commons.STRING.ValueOf(value)
	case BOOLEAN:
		return commons.BOOLEAN.ValueOf(value)
	case CURRENCY, INT, DOUBLE, PERCENT:
		return commons.NUMBER.ValueOf(value)
	case DATE:
		return commons.DATE.ValueOf(value)
	case TIME:
		return commons.TIME.ValueOf(value)
	case DATETIME:
		return commons.DATETIME.ValueOf(value)
	case LOCATION, ADDRESS, ANYTYPE:
		return nil, errors.New(fmt.Sprint("unsupported type"))
	}
	return nil, errors.New(fmt.Sprint("unknown type"))
}
