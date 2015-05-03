package soap

import (
	"github.com/goforce/log"
	"io"
	"time"
)

// QueryResult is parsed result in queryResponse. SObjectTypes for upper level QueryResult only will contain
// all the SObject types parsed.
type QueryResult struct {
	Done         bool
	QueryLocator string
	Records      []Record
	Size         int
	sObjectTypes map[string]struct{}
}

func Query(co Connection, query string, headers ...*XE) (*QueryResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "Query took:", time.Since(start)) }()
	qr := &QueryResult{sObjectTypes: make(map[string]struct{})}
	result := struct {
		QueryResult interface{} `xml:"Body>queryResponse>result"`
	}{qr}
	err := Call(co, &XE{"tns:query", XE{"tns:queryString", query}}, &result, headers...)
	return qr, err
}

func QueryMore(co Connection, queryLocator string, headers ...*XE) (*QueryResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "QueryMore took:", time.Since(start)) }()
	if queryLocator == "" {
		return nil, io.EOF
	}
	qr := &QueryResult{sObjectTypes: make(map[string]struct{})}
	result := struct {
		QueryResult interface{} `xml:"Body>queryMoreResponse>result"`
	}{qr}
	err := Call(co, &XE{"tns:queryMore", XE{"tns:queryLocator", queryLocator}}, &result, headers...)
	return qr, err
}

// QueryNext does QueryMore API call and appends results to QueryResult passed.
func QueryNext(co Connection, qr *QueryResult, headers ...*XE) (*QueryResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "QueryMore took:", time.Since(start)) }()
	if qr.QueryLocator == "" {
		return qr, io.EOF
	}
	result := struct {
		QueryResult interface{} `xml:"Body>queryMoreResponse>result"`
	}{qr}
	err := Call(co, &XE{"tns:queryMore", XE{"tns:queryLocator", qr.QueryLocator}}, &result, headers...)
	return qr, err
}

func (qr *QueryResult) SObjectTypes() []string {
	re := make([]string, 0, len(qr.sObjectTypes))
	for k, _ := range qr.sObjectTypes {
		re = append(re, k)
	}
	return re
}

// Append growths QueryResult Records linearly by 1000 records
func (qr *QueryResult) Append(r Record) {
	if len(qr.Records) == cap(qr.Records) {
		cp := make([]Record, len(qr.Records), cap(qr.Records)+1000)
		copy(cp, qr.Records)
		qr.Records = cp
	}
	qr.Records = append(qr.Records, r)
}
