package core

import (
	"errors"
	"github.com/goforce/api/commons"
	"github.com/goforce/log"
	"io"
	"time"
)

type QueryResult struct {
	Done         bool
	QueryLocator string
	Recordset    commons.Recordset
	Size         int
}

func Query(co commons.Connection, query string, headers ...*XE) (*QueryResult, error) {
	var qr QueryResult
	err := QueryModel(co, query, &qr, headers)
	return &qr, err
}

// QueryModel is intended for internal use to implement different query result parsers.
// Use Query which will use core.QueryResult to parse results.
func QueryModel(co commons.Connection, query string, target interface{}, headers []*XE) error {
	start := time.Now()
	defer func() { log.Println(DURATION, "Query took:", time.Since(start)) }()
	result := struct {
		QueryResult interface{} `xml:"Body>queryResponse>result"`
	}{target}
	return Call(co, &XE{"tns:query", XE{"tns:queryString", query}}, &result, headers...)
}

func QueryMore(co commons.Connection, queryLocator string, headers ...*XE) (*QueryResult, error) {
	var qr QueryResult
	err := QueryMoreModel(co, queryLocator, &qr, headers)
	return &qr, err
}

func QueryMoreModel(co commons.Connection, queryLocator string, target interface{}, headers []*XE) error {
	start := time.Now()
	defer func() { log.Println(DURATION, "QueryMore took:", time.Since(start)) }()
	if queryLocator == "" {
		return io.EOF
	}
	result := struct {
		QueryResult interface{} `xml:"Body>queryMoreResponse>result"`
	}{target}
	return Call(co, &XE{"tns:queryMore", XE{"tns:queryLocator", queryLocator}}, &result, headers...)
}

func (qr *QueryResult) Records() commons.Recordset {
	if qr == nil || qr.Recordset == nil {
		return make(commons.Recordset, 0, 0)
	}
	return commons.Recordset(qr.Recordset)
}

func (qr *QueryResult) Next() (commons.QueryLocator, error) {
	if qr.Done {
		return nil, io.EOF
	}
	return nil, errors.New("Next() not supported by QueryResult.")
}

func (qr *QueryResult) All() (commons.Recordset, error) {
	if qr.Done {
		return qr.Records(), nil
	}
	return nil, errors.New("All() not supported by QueryResult.")
}

func (qr *QueryResult) IsLast() bool {
	return qr.Done
}

func (qr *QueryResult) TotalSize() int {
	return qr.Size
}
