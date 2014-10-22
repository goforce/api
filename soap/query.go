package soap

import (
	"encoding/xml"
	"fmt"
	"github.com/goforce/api/commons"
	"github.com/goforce/log"
	"io"
	"strconv"
	"time"
)

var _ = fmt.Print

func (co *Connection) NewRecord(sObjectType string) (commons.Record, error) {
	if co.UseStrings {
		return commons.NewUndescribedRecord(sObjectType), nil
	} else {
		return commons.NewDescribedRecord(co.DescribeSObject(sObjectType))
	}
}

type QueryResult struct {
	done         bool
	queryLocator string
	records      []commons.Record
	size         int
	co           *Connection
	headers      []SoapHeader
	useStrings   bool
}

func (co *Connection) Query(query string, headers ...SoapHeader) (commons.QueryLocator, error) {
	start := time.Now()
	defer func() { log.Println(commons.DURATIONS, "Query took:", time.Since(start)) }()
	req := []KeyValue{
		KeyValue{"tns:query/tns:queryString", query},
	}
	response, err := co.Call(req, headers...)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	return parseXmlQueryResponse(response, &QueryResult{co: co, headers: headers, useStrings: co.UseStrings})
}

func (qr *QueryResult) TotalSize() int {
	if qr == nil {
		return 0
	}
	return qr.size
}

func (qr *QueryResult) Done() bool {
	if qr == nil {
		return true
	}
	return qr.done
}

func (qr *QueryResult) QueryMore() (commons.QueryLocator, error) {
	start := time.Now()
	defer func() { log.Println(commons.DURATIONS, "QueryMore took:", time.Since(start)) }()
	if qr == nil || qr.queryLocator == "" {
		return nil, io.EOF
	}
	req := []KeyValue{
		KeyValue{"tns:queryMore/tns:queryLocator", qr.queryLocator},
	}
	response, err := qr.co.Call(req, qr.headers...)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	l, err := parseXmlQueryResponse(response, qr.clone())
	return l, err
	//	return parseXmlQueryResponse(response, qr.clone())
}

// All() returns all the records from QueryLocator until QueryMore encounters end of data.
func (qr *QueryResult) All() ([]commons.Record, error) {
	if qr == nil {
		return make([]commons.Record, 0), nil
	}
	results := qr.records
	var ql commons.QueryLocator
	var err error = nil
	for err == nil && !qr.done {
		ql, err = qr.QueryMore()
		qr = ql.(*QueryResult)
		results = append(results, qr.records...)
	}
	if err != nil && err != io.EOF {
		return nil, err
	}
	return results, nil
}

func (qr *QueryResult) Records() []commons.Record {
	if qr == nil {
		return make([]commons.Record, 0)
	}
	return qr.records
}

func (qr *QueryResult) clone() *QueryResult {
	return &QueryResult{co: qr.co, headers: qr.headers, useStrings: qr.useStrings}
}

func parseXmlQueryResponse(reader io.Reader, qr *QueryResult) (*QueryResult, error) {
	xqr := xmlQueryResult{context: qr}
	xqr.newQueryResult = func() *QueryResult {
		return qr.clone()
	}
	xqr.newRecord = func(sObjectType string) (commons.Record, error) {
		if qr.useStrings {
			return commons.NewUndescribedRecord(sObjectType), nil
		} else {
			return commons.NewDescribedRecord(qr.co.DescribeSObject(sObjectType))
		}
	}
	response := struct {
		QueryResult     *xmlQueryResult `xml:"Body>queryResponse>result"`
		QueryMoreResult *xmlQueryResult `xml:"Body>queryMoreResponse>result"`
	}{&xqr, &xqr}
	err := xml.NewDecoder(reader).Decode(&response)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

type xmlPassThrough struct {
	newRecord      func(string) (commons.Record, error)
	newQueryResult func() *QueryResult
}

type xmlQueryResult struct {
	context *QueryResult
	xmlPassThrough
}

func (xqr *xmlQueryResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var se xml.StartElement
	var text string
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		switch t.(type) {
		case xml.StartElement:
			se = t.(xml.StartElement)
			text = ""
			if se.Name.Space == "urn:partner.soap.sforce.com" {
				if se.Name.Local == "records" {
					rec := xmlRecord{xmlPassThrough: xqr.xmlPassThrough}
					if err = d.DecodeElement(&rec, &se); err != nil {
						return err
					}
					xqr.context.records = append(xqr.context.records, rec.context)
				}
			}
		case xml.EndElement:
			if se.Name.Space == "urn:partner.soap.sforce.com" {
				if se.Name.Local == "done" {
					xqr.context.done = text == "true"
				} else if se.Name.Local == "queryLocator" {
					xqr.context.queryLocator = text
				} else if se.Name.Local == "size" {
					if xqr.context.size, err = strconv.Atoi(text); err != nil {
						return err
					}
				}
			}
			se = xml.StartElement{}
		case xml.CharData:
			text += string(t.(xml.CharData))
		default:
			log.Println(commons.ERRORS, "parsing xml query result: unknown type: ", t)
		}
	}
	return nil
}

type xmlRecord struct {
	context commons.Record
	xmlPassThrough
}

func (xr *xmlRecord) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var se xml.StartElement
	var text string
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		switch t.(type) {
		case xml.StartElement:
			se = t.(xml.StartElement)
			text = ""
			for _, attr := range se.Attr {
				if attr.Name.Space == NS_XSI {
					if attr.Name.Local == "type" {
						if attr.Value == "QueryResult" {
							xqr := xmlQueryResult{context: xr.newQueryResult(), xmlPassThrough: xr.xmlPassThrough}
							if err := d.DecodeElement(&xqr, &se); err != nil {
								return err
							}
							if _, err := xr.context.Set(se.Name.Local, xqr.context); err != nil {
								return err
							}
						} else if attr.Value == "sf:sObject" {
							rec := xmlRecord{xmlPassThrough: xr.xmlPassThrough}
							if err := d.DecodeElement(&rec, &se); err != nil {
								return err
							}
							if _, err := xr.context.Set(se.Name.Local, rec.context); err != nil {
								return err
							}
						}
					} else if attr.Name.Local == "nil" {
						if _, err := xr.context.Set(se.Name.Local, nil); err != nil {
							return err
						}
					}
					// set se to emtpy to ignore any text within start end tags
					se = xml.StartElement{}
				}
			}
		case xml.EndElement:
			if se.Name.Space == "urn:sobject.partner.soap.sforce.com" {
				if se.Name.Local == "type" {
					if xr.context, err = xr.newRecord(text); err != nil {
						return err
					}
				} else {
					if _, err := xr.context.Set(se.Name.Local, text); err != nil {
						return err
					}
				}
			}
			se = xml.StartElement{}
		case xml.CharData:
			text += string(t.(xml.CharData))
		default:
			log.Println(commons.ERRORS, "parsing xml record: unknown type: ", t)
		}
	}
	return nil
}
