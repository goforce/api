package data

import (
	"github.com/goforce/api/soap"
	"io"
)

type QueryResult interface {
	Records() Records
	Reader() Reader
	ExpectedLength() int
}

type Reader interface {
	Read() (Record, error)
}

type query struct {
	co      *Connection
	qr      *soap.QueryResult
	headers []*soap.XE
	ois     map[string]*SObjectInfo
}

type reader struct {
	q        query
	location int
}

func (q *query) Reader() Reader {
	q.getAllDetails()
	return &reader{q: *q, location: 0}
}

func (r *reader) Read() (rec Record, err error) {
	if r.location < len(r.q.qr.Records) {
		rec, err = r.q.wrapRow(r.q.qr.Records[r.location].(soap.Row))
		r.location++
		return
	} else if r.q.qr.Done {
		return nil, io.EOF
	} else {
		qr := r.q.qr
		qr.Records = qr.Records[0:0]
		r.q.qr, err = soap.QueryNext(r.q.co, qr, r.q.headers...)
		if err != nil {
			panic(err)
		}
		r.q.getAllDetails()
		r.location = 0
		if len(r.q.qr.Records) > 0 {
			rec, err = r.q.wrapRow(r.q.qr.Records[r.location].(soap.Row))
			r.location++
			return
		} else {
			return nil, io.EOF
		}
	}
}

func (co *Connection) Query(soql string, headers ...SoapHeader) (QueryResult, error) {
	xe := makeQueryXE(headers)
	qr, err := soap.Query(co, soql, xe...)
	if err != nil {
		return nil, err
	}
	return &query{co: co, qr: qr, headers: xe}, nil
}

func (co *Connection) QueryAll(soql string, headers ...SoapHeader) (Records, error) {
	qr, err := co.Query(soql, headers...)
	if err != nil {
		return nil, err
	}
	return qr.Records(), nil
}

func (q *query) Records() Records {
	var err error
	qr := q.qr
	for !qr.Done {
		qr, err = soap.QueryNext(q.co, q.qr, q.headers...)
		if err != nil {
			panic(err)
		}
	}
	q.getAllDetails()
	return q.wrapRecords(q.qr.Records)
}

func (q *query) getAllDetails() {
	var err error
	q.ois, err = q.co.DescribeSObjects(q.qr.SObjectTypes()...)
	if err != nil {
		panic(err)
	}
	for _, oi := range q.ois {
		_ = oi.Fields()
	}
	return
}

func (q *query) ExpectedLength() int { return q.qr.Size }

func (q *query) wrapRecords(rs []soap.Record) Records {
	var name string
	var oi *SObjectInfo
	var err error
	records := make(Records, len(rs))
	for i, r := range rs {
		if name != r.SObjectType() {
			name = r.SObjectType()
			oi, _ = q.ois[name]
		}
		row := r.(soap.Row)
		if records[i], err = q.wrapRowWithDetails(row, oi); err != nil {
			panic(err)
		}
	}
	return records
}

func (q *query) wrapQueryResult(qr *soap.QueryResult) (*query, error) {
	return &query{co: q.co, qr: qr, headers: q.headers, ois: q.ois}, nil
}

func (q *query) wrapRow(row soap.Row) (*record, error) {
	oi, _ := q.ois[row.SObjectType()]
	return q.wrapRowWithDetails(row, oi)
}

func (q *query) wrapRowWithDetails(row soap.Row, oi *SObjectInfo) (*record, error) {
	var err error
	for k, v := range row {
		switch v.(type) {
		case soap.Row:
			row[k], err = q.wrapRow(v.(soap.Row))
		case *soap.QueryResult:
			row[k], err = q.wrapQueryResult(v.(*soap.QueryResult))
		default:
			if oi != nil && v != nil {
				if f, ok := oi.fields.byName[k]; ok {
					row[k], err = f.Type().Parse(v.(string))
				}
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return &record{values: row, info: oi}, nil
}

func makeQueryXE(headers []SoapHeader) []*soap.XE {
	xe := make([]*soap.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	return xe
}
