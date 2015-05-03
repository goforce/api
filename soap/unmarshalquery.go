package soap

import (
	"encoding/xml"
	"errors"
	"github.com/goforce/api/conv"
	"io"
	"strconv"
	"strings"
)

type xmlRecord struct {
	r  Row
	qr *QueryResult
}

func (qr *QueryResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var (
		s   string
		err error
	)
	if s, err = ParseTextElement(d, xml.Name{NS_PARTNER, "done"}); err != nil {
		return err
	}
	qr.Done = s == "true"
	if qr.QueryLocator, err = ParseTextElement(d, xml.Name{NS_PARTNER, "queryLocator"}); err != nil {
		return err
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Space == NS_PARTNER && se.Name.Local == "records" {
			x := xmlRecord{qr: qr}
			err := d.DecodeElement(&x, &se)
			if err != nil {
				return err
			}
			qr.Append(x.r)
		} else {
			break
		}
	}
	if s, err = ParseTextElementComplete(d, xml.Name{NS_PARTNER, "size"}); err != nil {
		return err
	}
	qr.Size, err = strconv.Atoi(s)
	// skip till EOF
	for ; err == nil; _, err = d.Token() {
	}
	if err == io.EOF {
		return nil
	}
	return err
}

var nameXsiType xml.Name = xml.Name{NS_XSI, "type"}
var nameXsiNil xml.Name = xml.Name{NS_XSI, "nil"}

func (x *xmlRecord) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var (
		s   string
		err error
	)
	if s, err = ParseTextElement(d, xml.Name{NS_OBJECT, "type"}); err != nil {
		return err
	}
	x.r = NewRow(s)
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
				if attr.Name == nameXsiType {
					xsitype = attr.Value
				} else if attr.Name == nameXsiNil {
					xsinil = attr.Value == "true"
				}
			}
			name := strings.ToLower(se.Name.Local)
			if xsinil {
				x.r[name] = nil
			} else {
				switch xsitype {
				case "QueryResult":
					q := &QueryResult{sObjectTypes: x.qr.sObjectTypes}
					if err := d.DecodeElement(q, &se); err != nil {
						return err
					}
					x.r[name] = q
				case "sf:sObject":
					nx := xmlRecord{qr: x.qr}
					if err := d.DecodeElement(&nx, &se); err != nil {
						return err
					}
					x.r[name] = nx.r
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
						x.r[name], err = conv.NUMBER.Parse(text)
					case "xsd:dateTime":
						x.r[name], err = conv.DATETIME.Parse(text)
					case "xsd:date":
						x.r[name], err = conv.DATE.Parse(text)
					case "xsd:time":
						x.r[name], err = conv.TIME.Parse(text)
					default:
						x.r[name] = text
					}
					if err != nil {
						return err
					}
				}
			}
		}
	}
	x.qr.sObjectTypes[x.r.SObjectType()] = struct{}{}
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
