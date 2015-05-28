package data

import (
	"github.com/goforce/api/soap"
)

const MAX_BATCH_SIZE int = 200

type DmlError soap.DmlError

type DmlFailure struct {
	Record Record
	Index  int
	Errors []DmlError
}

func (co *Connection) Insert(records []Record, headers ...SoapHeader) []DmlFailure {
	xe := make([]*soap.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	failures := make([]DmlFailure, 0)
	for start, end := 0, 0; end < len(records); start = start + MAX_BATCH_SIZE {
		end = start + MAX_BATCH_SIZE
		if end > len(records) {
			end = len(records)
		}
		results, err := soap.Insert(co, convertToSoapRecords(records[start:end]), xe...)
		if err != nil {
			panic(err)
		}
		for i, r := range results {
			if r.Success {
				records[i+start].Set("Id", r.Id)
			} else {
				f := DmlFailure{Record: records[i+start], Index: i + start, Errors: make([]DmlError, len(r.Errors))}
				for j, e := range r.Errors {
					f.Errors[j] = DmlError(e)
				}
				failures = append(failures, f)
			}
		}
	}
	return failures
}

func (co *Connection) Update(records []Record, headers ...SoapHeader) ([]DmlFailure, error) {
	xe := make([]*soap.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	failures := make([]DmlFailure, 0)
	for start, end := 0, 0; end < len(records); start = start + MAX_BATCH_SIZE {
		end = start + MAX_BATCH_SIZE
		if end > len(records) {
			end = len(records)
		}
		results, err := soap.Update(co, convertToSoapRecords(records[start:end]), xe...)
		if err != nil {
			return failures, err
		}
		for i, r := range results {
			if !r.Success {
				f := DmlFailure{Record: records[i+start], Index: i + start, Errors: make([]DmlError, len(r.Errors))}
				for j, e := range r.Errors {
					f.Errors[j] = DmlError(e)
				}
				failures = append(failures, f)
			}
		}
	}
	return failures, nil
}

func (co *Connection) Upsert(records []Record, externalId string, headers ...SoapHeader) ([]DmlFailure, error) {
	xe := make([]*soap.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	failures := make([]DmlFailure, 0)
	for start, end := 0, 0; end < len(records); start = start + MAX_BATCH_SIZE {
		end = start + MAX_BATCH_SIZE
		if end > len(records) {
			end = len(records)
		}
		results, err := soap.Upsert(co, convertToSoapRecords(records[start:end]), externalId, xe...)
		if err != nil {
			return failures, err
		}
		for i, r := range results {
			if r.Success && r.Created {
				records[i+start].Set("Id", r.Id)
			} else if !r.Success {
				f := DmlFailure{Record: records[i+start], Index: i + start, Errors: make([]DmlError, len(r.Errors))}
				for j, e := range r.Errors {
					f.Errors[j] = DmlError(e)
				}
				failures = append(failures, f)
			}
		}
	}
	return failures, nil
}

func (co *Connection) Delete(records []Record, headers ...SoapHeader) ([]DmlFailure, error) {
	xe := make([]*soap.XE, 0, len(headers))
	for _, h := range headers {
		xe = append(xe, h.xe())
	}
	failures := make([]DmlFailure, 0)
	for start, end := 0, 0; end < len(records); start = start + MAX_BATCH_SIZE {
		end = start + MAX_BATCH_SIZE
		if end > len(records) {
			end = len(records)
		}
		results, err := soap.Delete(co, convertToSoapRecords(records[start:end]), xe...)
		if err != nil {
			return failures, err
		}
		for i, r := range results {
			if !r.Success {
				f := DmlFailure{Record: records[i+start], Index: i + start, Errors: make([]DmlError, len(r.Errors))}
				for j, e := range r.Errors {
					f.Errors[j] = DmlError(e)
				}
				failures = append(failures, f)
			}
		}
	}
	return failures, nil
}

func convertToSoapRecords(records []Record) []soap.Record {
	sr := make([]soap.Record, len(records))
	for i, r := range records {
		sr[i] = r
	}
	return sr
}
