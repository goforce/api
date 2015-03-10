package core

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/goforce/api/commons"
	"io"
	"strconv"
	"strings"
)

var _ = fmt.Println

const (
	QUERY_SIZE = 1000
)

func (qr *QueryResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Space == NS_PARTNER {
			if se.Name.Local == "records" {
				qr.Recordset = qr.Recordset.GrowIfFull()
				r := make(Record)
				qr.Recordset = append(qr.Recordset, r)
				d.DecodeElement(&r, &se)
			} else {
				text, err := ParseTextElementComplete(d, se.Name)
				if err != nil {
					return err
				}
				switch se.Name.Local {
				case "done":
					qr.Done = text == "true"
				case "queryLocator":
					qr.QueryLocator = text
				case "size":
					qr.Size, _ = strconv.Atoi(text)
				}
			}
		}
	}
	return nil
}

var NameXsiType xml.Name = xml.Name{NS_XSI, "type"}
var NameXsiNil xml.Name = xml.Name{NS_XSI, "nil"}

func (r Record) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	if r[sobjecttype], err = ParseTextElement(d, xml.Name{NS_OBJECT, "type"}); err != nil {
		return err
	}
	if _, err = ParseTextElement(d, xml.Name{NS_OBJECT, "Id"}); err != nil {
		return err
	}
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			xsitype := ""
			xsinil := false
			for _, attr := range se.Attr {
				if attr.Name == NameXsiType {
					xsitype = attr.Value
				} else if attr.Name == NameXsiNil {
					xsinil = attr.Value == "true"
				}
			}
			name := strings.ToLower(se.Name.Local)
			if xsinil {
				r[name] = nil
			} else {
				switch xsitype {
				case "QueryResult":
					var q QueryResult
					if err := d.DecodeElement(&q, &se); err != nil {
						return err
					}
					r[name] = &q
				case "sf:sObject":
					nested := make(Record)
					if err := d.DecodeElement(&nested, &se); err != nil {
						return err
					}
					r[name] = nested
				case "location":
					d.Skip()
				case "address":
					d.Skip()
				default:
					text, err := ParseTextElementComplete(d, se.Name)
					if err != nil {
						return err
					}
					switch xsitype {
					case "xsd:int", "xsd:double":
						r[name], err = commons.NUMBER.Parse(text)
					case "xsd:dateTime":
						r[name], err = commons.DATETIME.Parse(text)
					case "xsd:date":
						r[name], err = commons.DATE.Parse(text)
					case "xsd:time":
						r[name], err = commons.TIME.Parse(text)
					default:
						r[name] = text
					}
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func ParseTextElement(d *xml.Decoder, ele xml.Name) (string, error) {
	t, err := d.Token()
	if v, ok := t.(xml.StartElement); err != nil || !ok || v.End().Name != ele {
		return "", xmlError(ele.Local)
	}
	return ParseTextElementComplete(d, ele)
}

func ParseTextElementComplete(d *xml.Decoder, ele xml.Name) (string, error) {
	t, err := d.Token()
	if err != nil {
		return "", xmlError(ele.Local)
	}
	if v, ok := t.(xml.EndElement); ok && v.Name == ele {
		return "", nil
	}
	re, ok := t.(xml.CharData)
	if !ok {
		return "", xmlError(ele.Local)
	}
	re = re.Copy()
	t, err = d.Token()
	if v, ok := t.(xml.EndElement); err != nil || !ok || v.Name != ele {
		return "", xmlError(ele.Local)
	}
	return string(re), nil
}

func xmlError(name string) error {
	return errors.New("unexpected error parsing xml at " + name)
}
