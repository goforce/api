package soap

import (
	"io"
)

type Reader interface {
	Read() (Record, error)
	Unread() bool
	ExpectedLength() int
	SObjectTypes() []string
}

type reader struct {
	co       Connection
	headers  []*XE
	qr       *QueryResult
	location int
}

func NewReader(co Connection, query string, headers ...*XE) (Reader, error) {
	qr, err := Query(co, query, headers...)
	if err != nil {
		return nil, err
	}
	return &reader{co: co, headers: headers, qr: qr, location: 0}, nil
}

func (reader *reader) SObjectTypes() []string {
	return reader.qr.SObjectTypes()
}

func (reader *reader) ExpectedLength() int {
	return reader.qr.Size
}

func (reader *reader) Unread() bool {
	if reader.location > 0 {
		reader.location--
		return true
	}
	return false
}

func (reader *reader) Read() (rec Record, err error) {
	if reader.location < len(reader.qr.Records) {
		rec = reader.qr.Records[reader.location]
		reader.location++
		return rec, nil
	} else if reader.qr.Done {
		return nil, io.EOF
	} else {
		reader.qr.Records = reader.qr.Records[0:0]
		reader.qr, err = QueryNext(reader.co, reader.qr, reader.headers...)
		if err != nil {
			return nil, err
		}
		reader.location = 0
		if len(reader.qr.Records) > 0 {
			rec = reader.qr.Records[0]
			reader.location++
			return
		} else {
			return nil, io.EOF
		}
	}
}
