package soap

import (
	"github.com/goforce/api/commons"
	"github.com/goforce/api/soap/core"
)

type DmlError struct {
	err    *core.DmlError
	record commons.Record
	index  int
}

const MAX_BATCH_SIZE int = 200

func (e *DmlError) Record() commons.Record { return e.record }
func (e *DmlError) Fields() []string       { return e.err.Fields }
func (e *DmlError) StatusCode() string     { return e.err.StatusCode }
func (e *DmlError) Message() string        { return e.err.Message }
func (e *DmlError) Index() int             { return e.index }

func (co *Connection) Insert(records []commons.Record, headers ...SoapHeader) ([]DmlError, error) {
	xe := make([]*core.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	results, err := core.Insert(co, records, xe...)
	if err != nil {
		return nil, err
	}
	errs := make([]DmlError, 0, MAX_BATCH_SIZE)
	for i, r := range results {
		if !r.Success {
			for _, err := range r.Errors {
				errs = append(errs, DmlError{err: &err, record: records[i], index: i})
			}
		}
	}
	return errs, nil
}

func (co *Connection) Update(records []commons.Record, headers ...SoapHeader) ([]DmlError, error) {
	xe := make([]*core.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	results, err := core.Update(co, records, xe...)
	if err != nil {
		return nil, err
	}
	errs := make([]DmlError, 0, MAX_BATCH_SIZE)
	for i, r := range results {
		if !r.Success {
			for _, err := range r.Errors {
				errs = append(errs, DmlError{err: &err, record: records[i], index: i})
			}
		}
	}
	return errs, nil
}

func (co *Connection) Upsert(records []commons.Record, externalId string, headers ...SoapHeader) ([]DmlError, error) {
	xe := make([]*core.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	results, err := core.Upsert(co, records, externalId, xe...)
	if err != nil {
		return nil, err
	}
	errs := make([]DmlError, 0, MAX_BATCH_SIZE)
	for i, r := range results {
		if !r.Success {
			for _, err := range r.Errors {
				errs = append(errs, DmlError{err: &err, record: records[i], index: i})
			}
		}
	}
	return errs, nil
}

func (co *Connection) Delete(records []commons.Record, headers ...SoapHeader) ([]DmlError, error) {
	xe := make([]*core.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	results, err := core.Delete(co, records, xe...)
	if err != nil {
		return nil, err
	}
	errs := make([]DmlError, 0, MAX_BATCH_SIZE)
	for i, r := range results {
		if !r.Success {
			for _, err := range r.Errors {
				errs = append(errs, DmlError{err: &err, record: records[i], index: i})
			}
		}
	}
	return errs, nil
}
