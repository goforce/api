package data

import (
	"errors"
	"fmt"
	"github.com/goforce/api/conv"
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

func (ft FieldType) String() string {
	switch ft {
	case STRING:
		return "string"
	case PICKLIST:
		return "picklist"
	case MULTIPICKLIST:
		return "multipicklist"
	case COMBOBOX:
		return "combobox"
	case REFERENCE:
		return "reference"
	case BASE64:
		return "base64"
	case BOOLEAN:
		return "boolean"
	case CURRENCY:
		return "currency"
	case TEXTAREA:
		return "textarea"
	case INT:
		return "int"
	case DOUBLE:
		return "double"
	case PERCENT:
		return "percent"
	case PHONE:
		return "phone"
	case ID:
		return "id"
	case DATE:
		return "date"
	case DATETIME:
		return "datetime"
	case TIME:
		return "time"
	case URL:
		return "url"
	case EMAIL:
		return "email"
	case ENCRYPTEDSTRING:
		return "encryptedstring"
	case DATACATEGORYGROUPREFERENCE:
		return "datacategorygroupreference"
	case LOCATION:
		return "location"
	case ADDRESS:
		return "address"
	case ANYTYPE:
		return "anyType"
	}
	return "unknown"
}

// Value validates and converts value to one which would match field type
func (t FieldType) Parse(value string) (interface{}, error) {
	switch t {
	case STRING, PICKLIST, MULTIPICKLIST, COMBOBOX, REFERENCE, BASE64, TEXTAREA, PHONE, ID, URL, EMAIL, ENCRYPTEDSTRING, DATACATEGORYGROUPREFERENCE:
		return conv.STRING.Parse(value)
	case BOOLEAN:
		return conv.BOOLEAN.Parse(value)
	case CURRENCY, INT, DOUBLE, PERCENT:
		return conv.NUMBER.Parse(value)
	case DATE:
		return conv.DATE.Parse(value)
	case TIME:
		return conv.TIME.Parse(value)
	case DATETIME:
		return conv.DATETIME.Parse(value)
	case LOCATION, ADDRESS, ANYTYPE:
		return nil, errors.New(fmt.Sprint("unsupported type"))
	}
	return nil, errors.New(fmt.Sprint("unknown type"))
}

// Value validates and converts value to one which would match field type
func (t FieldType) ValueOf(value interface{}) (interface{}, error) {
	switch t {
	case STRING, PICKLIST, MULTIPICKLIST, COMBOBOX, REFERENCE, BASE64, TEXTAREA, PHONE, ID, URL, EMAIL, ENCRYPTEDSTRING, DATACATEGORYGROUPREFERENCE:
		return conv.STRING.ValueOf(value)
	case BOOLEAN:
		return conv.BOOLEAN.ValueOf(value)
	case CURRENCY, INT, DOUBLE, PERCENT:
		return conv.NUMBER.ValueOf(value)
	case DATE:
		return conv.DATE.ValueOf(value)
	case TIME:
		return conv.TIME.ValueOf(value)
	case DATETIME:
		return conv.DATETIME.ValueOf(value)
	case LOCATION, ADDRESS, ANYTYPE:
		return nil, errors.New(fmt.Sprint("unsupported type"))
	}
	return nil, errors.New(fmt.Sprint("unknown type"))
}
