package soap

import (
	"encoding/xml"
	"fmt"
	"github.com/goforce/api/commons"
	"github.com/goforce/api/soap/core"
	"io"
	"strconv"
	"strings"
)

const (
	QUERY_SIZE = 1000
)

type recordDecoder struct {
	ql     *QueryLocator
	record *Record
}

func (ql *QueryLocator) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if se, ok := t.(xml.StartElement); ok && se.Name.Space == core.NS_PARTNER {
			if se.Name.Local == "records" {
				ql.records = ql.records.GrowIfFull()
				rd := recordDecoder{ql, &Record{}}
				ql.records = append(ql.records, rd.record)
				d.DecodeElement(&rd, &se)
			} else {
				text, err := core.ParseTextElementComplete(d, se.Name)
				if err != nil {
					return err
				}
				switch se.Name.Local {
				case "done":
					ql.isLast = text == "true"
				case "queryLocator":
					ql.queryLocator = text
				case "size":
					ql.totalSize, _ = strconv.Atoi(text)
				}
			}
		}
	}
	return nil
}

func (rd recordDecoder) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	rd.record.values = make(map[string]interface{})
	if text, err := core.ParseTextElement(d, xml.Name{core.NS_OBJECT, "type"}); err == nil {
		// AggregateResult is known type without describe.
		if text != "AggregateResult" {
			rd.record.describe, err = rd.ql.co.DescribeSObject(text)
			if err != nil {
				fmt.Println("failed to describe: ", text, " error returned:", err)
				return err
			}
		}
	} else {
		return err
	}
	if _, err := core.ParseTextElement(d, xml.Name{core.NS_OBJECT, "Id"}); err != nil {
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
				if attr.Name == core.NameXsiType {
					xsitype = attr.Value
				} else if attr.Name == core.NameXsiNil {
					xsinil = attr.Value == "true"
				}
			}
			name := strings.ToLower(se.Name.Local)
			if xsinil {
				rd.record.values[name] = nil
			} else {
				switch xsitype {
				case "QueryResult":
					ql := QueryLocator{co: rd.ql.co, headers: rd.ql.headers}
					if err := d.DecodeElement(&ql, &se); err != nil {
						return err
					}
					rd.record.values[name] = &ql
				case "sf:sObject":
					nested := recordDecoder{ql: rd.ql, record: &Record{}}
					if err := d.DecodeElement(&nested, &se); err != nil {
						return err
					}
					rd.record.values[name] = nested.record
				case "location":
					d.Skip()
				case "address":
					d.Skip()
				default:
					text, err := core.ParseTextElementComplete(d, se.Name)
					if err != nil {
						return err
					}
					switch xsitype {
					case "xsd:int", "xsd:double":
						rd.record.values[name], err = commons.NUMBER.Parse(text)
					case "xsd:dateTime":
						rd.record.values[name], err = commons.DATETIME.Parse(text)
					case "xsd:date":
						rd.record.values[name], err = commons.DATE.Parse(text)
					case "xsd:time":
						rd.record.values[name], err = commons.TIME.Parse(text)
					default:
						if rd.record.describe != nil {
							if d, ok := rd.record.describe.FieldByName(name); ok {
								v, err := d.FieldType().Parse(text)
								if err != nil {
									return err
								}
								rd.record.values[name] = v
							} else {
								rd.record.values[name] = text
							}
						} else {
							rd.record.values[name] = text
						}
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
