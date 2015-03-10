package soap

import (
	"github.com/goforce/api/commons"
	"github.com/goforce/api/soap/core"
	"io"
)

type XE core.XE

type QueryLocator struct {
	co           *Connection
	isLast       bool
	queryLocator string
	totalSize    int
	records      commons.Recordset
	headers      []SoapHeader
}

func (co *Connection) Query(query string, headers ...SoapHeader) (*QueryLocator, error) {
	ql := QueryLocator{co: co, headers: headers}
	err := core.QueryModel(co, query, &ql, soapHeadersToXe(ql.headers))
	return &ql, err
}

func (co *Connection) QueryAll(query string, headers ...SoapHeader) (commons.Recordset, error) {
	ql, err := co.Query(query, headers...)
	if err != nil {
		return nil, err
	}
	return ql.All()
}

func (ql *QueryLocator) Records() commons.Recordset {
	return ql.records
}

func (ql *QueryLocator) TotalSize() int {
	return ql.totalSize
}

func (ql *QueryLocator) IsLast() bool {
	return ql.isLast
}

func (ql *QueryLocator) Next() (commons.QueryLocator, error) {
	if ql == nil || ql.isLast {
		return nil, io.EOF
	}
	next := QueryLocator{co: ql.co, headers: ql.headers}
	err := core.QueryMoreModel(ql.co, ql.queryLocator, &next, soapHeadersToXe(ql.headers))
	return &next, err
}

func (ql *QueryLocator) All() (commons.Recordset, error) {
	if ql == nil {
		return make(commons.Recordset, 0, 0), nil
	}
	for !ql.isLast {
		err := core.QueryMoreModel(ql.co, ql.queryLocator, &ql, soapHeadersToXe(ql.headers))
		if err != nil {
			return nil, err
		}
	}
	return ql.records, nil
}

func NewQueryLocator() *QueryLocator {
	return &QueryLocator{isLast: true, records: make(commons.Recordset, 0, 100)}
}

func (ql *QueryLocator) Add(record commons.Record) {
	ql.records = append(ql.records, record)
}

func soapHeadersToXe(headers []SoapHeader) []*core.XE {
	xe := make([]*core.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	return xe
}
