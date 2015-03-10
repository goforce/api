package soap

import (
//"encoding/xml"
//"github.com/goforce/api/commons"
// 	"strings"
)

// type recordset commons.Recordset

// func (rsp *recordset) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	log.Println(ERRORS, "unmarshalling array of records: ", len(*rsp), cap(*rsp))
// 	r := make(record)
// 	rs := *rsp
// 	if len(rs) == cap(rs) {
// 		rs = make(recordset, 0, cap(rs)+queryRecordsSize)
// 		copy(rs, *rsp)
// 	}
// 	rs = append(rs, r)
// 	*rsp = rs
// 	return r.UnmarshalXML(d, start)
// }

// func (rsp *recordset) Records() commons.Recordset {
// 	return commons.Recordset(*rsp)
// }

// func (rsp *Recordset) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	log.Println(ERRORS, "unmarshalling array of records: ", len(*rsp), cap(*rsp))
// 	r := make(Record)
// 	rs := *rsp
// 	if len(rs) == cap(rs) {
// 		rs = make(Recordset, 0, cap(rs)+1000)
// 		copy(rs, *rsp)
// 	}
// 	rs = append(rs, r)
// 	*rsp = rs
// 	return r.UnmarshalXML(d, start)
// }

// func (r Record) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	log.Println(ERRORS, "unmarshalling Record")
// 	var text string
// 	var xsitype string
// 	for {
// 		t, err := d.Token()
// 		if err == io.EOF {
// 			break
// 		}
// 		switch t.(type) {
// 		case xml.StartElement:
// 			se := t.(xml.StartElement)
// 			text = ""
// 			xsitype = ""
// 			for _, attr := range se.Attr {
// 				if attr.Name.Space == NS_XSI {
// 					if attr.Name.Local == "type" {
// 						if attr.Value == "QueryResult" {
// 							q := &QueryResult{}
// 							if err := d.DecodeElement(q, &se); err != nil {
// 								return err
// 							}
// 							r[strings.ToLower(se.Name.Local)] = q
// 						} else if attr.Value == "sf:sObject" {
// 							var nr Record = make(Record)
// 							if err := d.DecodeElement(&nr, &se); err != nil {
// 								return err
// 							}
// 							r[strings.ToLower(se.Name.Local)] = nr
// 						} else {
// 							xsitype = attr.Value
// 						}
// 					} else if attr.Name.Local == "nil" {
// 						r[strings.ToLower(se.Name.Local)] = nil
// 					}
// 				}
// 			}
// 		case xml.EndElement:
// 			ee := t.(xml.EndElement)
// 			if ee.Name.Space == "urn:sobject.partner.soap.sforce.com" {
// 				name := strings.ToLower(ee.Name.Local)
// 				// if value added then do not overwrite it, there could be record, query result or nil
// 				if _, ok := r[name]; !ok {
// 					if len(xsitype) > 0 {
// 						var err error
// 						if xsitype == "xsd:int" || xsitype == "xsd:double" {
// 							r[name], err = commons.NUMBER.Parse(text)
// 						} else if xsitype == "xsd:dateTime" {
// 							r[name], err = commons.DATETIME.Parse(text)
// 						} else if xsitype == "xsd:date" {
// 							r[name], err = commons.DATE.Parse(text)
// 						} else if xsitype == "xsd:time" {
// 							r[name], err = commons.TIME.Parse(text)
// 						} else {
// 							r[name] = text
// 						}
// 						if err != nil {
// 							return err
// 						}
// 					} else {
// 						r[name] = text
// 					}
// 				}
// 			}
// 			text = ""
// 			xsitype = ""
// 		case xml.CharData:
// 			text += string(t.(xml.CharData))
// 		default:
// 			log.Println(ERRORS, "parsing xml record: unknown type: ", t)
// 		}
// 	}
// 	return nil
// }
