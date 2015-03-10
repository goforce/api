package soap

import (
	"github.com/goforce/api/soap/core"
)

type Field struct {
	describe  *core.Field
	fieldType FieldType
}

func newField(describe *core.Field) *Field {
	return &Field{describe: describe}
}

func (describe *Field) Name() string          { return describe.describe.Name }
func (describe *Field) ReferenceTo() []string { return describe.describe.ReferenceTo }
func (describe *Field) FieldType() FieldType {
	if describe.fieldType == 0 {
		switch describe.describe.Type {
		case "string":
			describe.fieldType = STRING
		case "picklist":
			describe.fieldType = PICKLIST
		case "multipicklist":
			describe.fieldType = MULTIPICKLIST
		case "combobox":
			describe.fieldType = COMBOBOX
		case "reference":
			describe.fieldType = REFERENCE
		case "base64":
			describe.fieldType = BASE64
		case "boolean":
			describe.fieldType = BOOLEAN
		case "currency":
			describe.fieldType = CURRENCY
		case "textarea":
			describe.fieldType = TEXTAREA
		case "int":
			describe.fieldType = INT
		case "double":
			describe.fieldType = DOUBLE
		case "percent":
			describe.fieldType = PERCENT
		case "phone":
			describe.fieldType = PHONE
		case "id":
			describe.fieldType = ID
		case "date":
			describe.fieldType = DATE
		case "datetime":
			describe.fieldType = DATETIME
		case "time":
			describe.fieldType = TIME
		case "url":
			describe.fieldType = URL
		case "email":
			describe.fieldType = EMAIL
		case "encryptedstring":
			describe.fieldType = ENCRYPTEDSTRING
		case "datacategorygroupreference":
			describe.fieldType = DATACATEGORYGROUPREFERENCE
		case "location":
			describe.fieldType = LOCATION
		case "address":
			describe.fieldType = ADDRESS
		case "anyType":
			describe.fieldType = ANYTYPE
		default:
			describe.fieldType = UNKNOWN
		}
	}
	return describe.fieldType
}
