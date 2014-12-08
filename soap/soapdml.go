package soap

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/goforce/api/commons"
	"github.com/goforce/log"
	"strings"
)

var _ = fmt.Print

type DmlResult struct {
	Created bool `xml:"created"`
	Errors  struct {
		Message    string `xml:"message"`
		StatusCode string `xml:"statusCode"`
	} `xml:"errors"`
	Id      string `xml:"id"`
	Success bool   `xml:"success"`
}

func (co *Connection) Insert(records []commons.Record, headers ...SoapHeader) (result []DmlResult, err error) {
	_ = log.IsOn(commons.CALLS) && log.Println(commons.CALLS, "insert ", len(records), " records")
	var req bytes.Buffer
	req.WriteString(`<tns:create>`)
	err = writeRecords(&req, records, "")
	if err != nil {
		return nil, err
	}
	req.WriteString(`</tns:create>`)
	_ = log.IsOn(commons.MESSAGES) && log.Println(commons.MESSAGES, string(req.Bytes()))
	response, err := co.Post(req.Bytes(), headers...)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	var createResponse struct {
		Results []DmlResult `xml:"Body>createResponse>result"`
	}
	xml.NewDecoder(response).Decode(&createResponse)
	for i, _ := range createResponse.Results {
		if createResponse.Results[i].Success {
			records[i].Set("Id", createResponse.Results[i].Id)
		}
	}
	return createResponse.Results, nil
}

func (co *Connection) Update(records []commons.Record, headers ...SoapHeader) (result []DmlResult, err error) {
	_ = log.IsOn(commons.CALLS) && log.Println(commons.CALLS, "update ", len(records), " records")
	var req bytes.Buffer
	req.WriteString(`<tns:update>`)
	err = writeRecords(&req, records, "")
	if err != nil {
		return nil, err
	}
	req.WriteString(`</tns:update>`)
	_ = log.IsOn(commons.MESSAGES) && log.Println(commons.MESSAGES, string(req.Bytes()))
	response, err := co.Post(req.Bytes(), headers...)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	var updateResponse struct {
		Results []DmlResult `xml:"Body>updateResponse>result"`
	}
	xml.NewDecoder(response).Decode(&updateResponse)
	return updateResponse.Results, nil
}

func (co *Connection) Upsert(records []commons.Record, externalId string, headers ...SoapHeader) (result []DmlResult, err error) {
	_ = log.IsOn(commons.CALLS) && log.Println(commons.CALLS, "upsert ", len(records), " records")
	var req bytes.Buffer
	req.WriteString(`<tns:upsert>`)
	if externalId == "" {
		externalId = "Id"
	}
	req.WriteString(`<tns:externalIDFieldName>`)
	req.WriteString(externalId)
	req.WriteString(`</tns:externalIDFieldName>`)
	err = writeRecords(&req, records, externalId)
	if err != nil {
		return nil, err
	}
	req.WriteString(`</tns:upsert>`)
	_ = log.IsOn(commons.MESSAGES) && log.Println(commons.MESSAGES, string(req.Bytes()))
	response, err := co.Post(req.Bytes(), headers...)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	var upsertResponse struct {
		Results []DmlResult `xml:"Body>upsertResponse>result"`
	}
	xml.NewDecoder(response).Decode(&upsertResponse)
	for i, _ := range upsertResponse.Results {
		if upsertResponse.Results[i].Success && upsertResponse.Results[i].Created {
			records[i].Set("Id", upsertResponse.Results[i].Id)
		}
	}
	return upsertResponse.Results, nil
}

func (co *Connection) Delete(records []commons.Record, headers ...SoapHeader) (result []DmlResult, err error) {
	_ = log.IsOn(commons.CALLS) && log.Println(commons.CALLS, "delete ", len(records), " records")
	var req bytes.Buffer
	req.WriteString(`<tns:delete>`)
	for _, rec := range records {
		if id, ok := rec.Get("Id"); ok {
			writeValue(&req, "tns:ids", id)
		} else {
			return nil, errors.New("Only records with Id field specified can be deleted.")
		}
	}
	req.WriteString(`</tns:delete>`)
	_ = log.IsOn(commons.MESSAGES) && log.Println(commons.MESSAGES, string(req.Bytes()))
	response, err := co.Post(req.Bytes(), headers...)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	var deleteResponse struct {
		Results []DmlResult `xml:"Body>deleteResponse>result"`
	}
	xml.NewDecoder(response).Decode(&deleteResponse)
	for i, _ := range deleteResponse.Results {
		if deleteResponse.Results[i].Success && deleteResponse.Results[i].Created {
			records[i].Set("Id", deleteResponse.Results[i].Id)
		}
	}
	return deleteResponse.Results, nil
}

func writeRecords(buf *bytes.Buffer, records []commons.Record, externalId string) error {
	var fieldsToWrite bytes.Buffer
	var fieldsToNull bytes.Buffer
	var nestedFieldsToWrite bytes.Buffer
	for _, rec := range records {
		fieldsToNull.Reset()
		fieldsToWrite.Reset()
		for _, f := range rec.Fields() {
			v := commons.Must(rec.Get(f))
			if commons.IsBlank(v) && strings.ToLower(f) != strings.ToLower(externalId) {
				fieldsToNull.WriteString(`<ens:fieldsToNull>`)
				fieldsToNull.WriteString(f)
				fieldsToNull.WriteString(`</ens:fieldsToNull>`)
			} else {
				switch v.(type) {
				case commons.QueryLocator:
					// do nothing, just skip nested subselects
				case commons.Record:
					nested := v.(commons.Record)
					nestedFieldsToWrite.Reset()
					for _, nf := range nested.Fields() {
						v := commons.Must(nested.Get(nf))
						if !commons.IsBlank(v) {
							writeValue(&nestedFieldsToWrite, nf, v)
						}
					}
					if nestedFieldsToWrite.Len() > 0 {
						fieldsToWrite.WriteByte('<')
						fieldsToWrite.WriteString(f)
						fieldsToWrite.WriteByte('>')
						fieldsToWrite.WriteString(`<type>`)
						fieldsToWrite.WriteString(nested.SObjectType())
						fieldsToWrite.WriteString(`</type>`)
						fieldsToWrite.Write(nestedFieldsToWrite.Bytes())
						fieldsToWrite.WriteString(`</`)
						fieldsToWrite.WriteString(f)
						fieldsToWrite.WriteByte('>')
					}
				default:
					writeValue(&fieldsToWrite, f, commons.Must(rec.Get(f)))
				}
			}
		}
		buf.WriteString(`<tns:sObjects><ens:type>`)
		buf.WriteString(string(rec.SObjectType()))
		buf.WriteString(`</ens:type>`)
		buf.Write(fieldsToNull.Bytes())
		buf.Write(fieldsToWrite.Bytes())
		buf.WriteString(`</tns:sObjects>`)
	}
	return nil
}

func writeValue(buf *bytes.Buffer, tag string, value interface{}) {
	buf.WriteByte('<')
	buf.WriteString(tag)
	buf.WriteByte('>')
	err := xml.EscapeText(buf, []byte(commons.String(value)))
	if err != nil {
		panic(err)
	}
	buf.WriteString("</")
	buf.WriteString(tag)
	buf.WriteByte('>')
}
