package soap

import (
	"errors"
	"github.com/goforce/api/conv"
	"github.com/goforce/log"
	"math/big"
	"strings"
	"time"
)

type DmlResult struct {
	Created bool       `xml:"created"`
	Errors  []DmlError `xml:"errors"`
	Id      string     `xml:"id"`
	Success bool       `xml:"success"`
}
type DmlError struct {
	Fields     []string `xml:"fields"`
	Message    string   `xml:"message"`
	StatusCode string   `xml:"statusCode"`
}

const fieldsLen int = 100

func Insert(co Connection, records []Record, headers ...*XE) ([]DmlResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "Insert took:", time.Since(start)) }()
	if len(records) == 0 {
		return make([]DmlResult, 0, 0), nil
	}
	fields := make([]string, 0, fieldsLen)
	fieldsToNull := make([]string, 0, fieldsLen)
	var req Buffer
	req.openTag("tns:create")
	for _, r := range records {
		if r == nil {
			return nil, errors.New("No nil Records allowed for Insert")
		}
		fields = fields[:0]
		fieldsToNull = fieldsToNull[:0]
		for _, f := range r.Fields() {
			v, _ := r.Get(f)
			if conv.IsBlank(v) {
				fieldsToNull = append(fieldsToNull, f)
			} else {
				fields = append(fields, f)
			}
		}
		err := writeRecord(&req, r, fields, fieldsToNull)
		if err != nil {
			return nil, err
		}
	}
	req.closeTag("tns:create")
	var result struct {
		DmlResults []DmlResult `xml:"Body>createResponse>result"`
	}
	err := Post(co.GetServerUrl(), req.Bytes(), &result, append(headers, &XE{"tns:SessionHeader", XE{"tns:sessionId", co.GetToken()}})...)
	if err != nil {
		return nil, err
	}
	return result.DmlResults, nil
}

func Update(co Connection, records []Record, headers ...*XE) ([]DmlResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "Update took:", time.Since(start)) }()
	if len(records) == 0 {
		return make([]DmlResult, 0, 0), nil
	}
	fields := make([]string, 0, fieldsLen)
	fieldsToNull := make([]string, 0, fieldsLen)
	var req Buffer
	req.openTag("tns:update")
	for _, r := range records {
		if r == nil {
			return nil, errors.New("No nil Records allowed for Update")
		}
		fields = fields[:0]
		fieldsToNull = fieldsToNull[:0]
		for _, f := range r.Fields() {
			v, _ := r.Get(f)
			if conv.IsBlank(v) {
				fieldsToNull = append(fieldsToNull, f)
			} else {
				fields = append(fields, f)
			}
			err := writeRecord(&req, r, fields, fieldsToNull)
			if err != nil {
				return nil, err
			}
		}
	}
	req.closeTag("tns:update")
	var result struct {
		DmlResults []DmlResult `xml:"Body>updateResponse>result"`
	}
	err := Post(co.GetServerUrl(), req.Bytes(), &result, append(headers, &XE{"tns:SessionHeader", XE{"tns:sessionId", co.GetToken()}})...)
	if err != nil {
		return nil, err
	}
	return result.DmlResults, nil
}

func Upsert(co Connection, records []Record, externalId string, headers ...*XE) ([]DmlResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "Upsert took:", time.Since(start)) }()
	if len(records) == 0 {
		return make([]DmlResult, 0, 0), nil
	}
	if externalId == "" {
		externalId = "Id"
	}
	fields := make([]string, 0, fieldsLen)
	fieldsToNull := make([]string, 0, fieldsLen)
	var req Buffer
	req.openTag("tns:upsert")
	req.writeElement("tns:externalIDFieldName", externalId)
	for _, r := range records {
		if r == nil {
			return nil, errors.New("No nil Records allowed for Upsert")
		}
		fields = fields[:0]
		fieldsToNull = fieldsToNull[:0]
		for _, f := range r.Fields() {
			v, _ := r.Get(f)
			if conv.IsBlank(v) && strings.ToLower(f) != strings.ToLower(externalId) {
				fieldsToNull = append(fieldsToNull, f)
			} else {
				fields = append(fields, f)
			}
			err := writeRecord(&req, r, fields, fieldsToNull)
			if err != nil {
				return nil, err
			}
		}
	}
	req.closeTag("tns:upsert")
	var result struct {
		DmlResults []DmlResult `xml:"Body>upsertResponse>result"`
	}
	err := Post(co.GetServerUrl(), req.Bytes(), &result, append(headers, &XE{"tns:SessionHeader", XE{"tns:sessionId", co.GetToken()}})...)
	if err != nil {
		return nil, err
	}
	return result.DmlResults, nil
}

func Delete(co Connection, records []Record, headers ...*XE) ([]DmlResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "Delete took:", time.Since(start)) }()
	if len(records) == 0 {
		return make([]DmlResult, 0, 0), nil
	}
	var req Buffer
	req.openTag("tns:delete")
	for _, rec := range records {
		if rec == nil {
			return nil, errors.New("No nil Records allowed for Delete")
		}
		if id, ok := rec.Get("Id"); ok {
			req.writeValue("tns:ids", id)
		} else {
			return nil, errors.New("All records for delete should have Id field specified.")
		}
	}
	req.closeTag("tns:delete")
	var result struct {
		DmlResults []DmlResult `xml:"Body>deleteResponse>result"`
	}
	err := Post(co.GetServerUrl(), req.Bytes(), &result, append(headers, &XE{"tns:SessionHeader", XE{"tns:sessionId", co.GetToken()}})...)
	if err != nil {
		return nil, err
	}
	return result.DmlResults, nil
}

func writeRecord(buf *Buffer, record Record, fields []string, fieldsToNull []string) error {
	buf.openTag("tns:sObjects")
	buf.writeElement("ens:type", record.SObjectType())
	for _, f := range fieldsToNull {
		buf.writeElement("ens:fieldsToNull", f)
	}
	for _, f := range fields {
		v, _ := record.Get(f)
		switch v.(type) {
		case *QueryResult:
			// do nothing, just skip nested subselects
		case Record:
			nested := v.(Record)
			hasNested := false
			for _, nf := range nested.Fields() {
				v, _ := nested.Get(nf)
				if !conv.IsBlank(v) {
					if !hasNested {
						buf.openTag(f)
						buf.writeElement("type", nested.SObjectType())
						hasNested = true
					}
					// ignore error, only simple fields should be written
					buf.writeValue(nf, v)
				}
			}
			if hasNested {
				buf.closeTag(f)
			}
		case bool, string, *big.Rat, time.Time:
			v, _ := record.Get(f)
			err := buf.writeValue(f, v)
			if err != nil {
				return err
			}
		default:
			// all unsupported types are skipped
			// TODO this is needed to skip QueryLocator from data. Should return error.
		}
	}
	buf.closeTag("tns:sObjects")
	return nil
}
