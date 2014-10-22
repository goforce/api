package commons

import (
	"fmt"
	"github.com/goforce/log"
	"io"
	"time"
)

var _ = fmt.Print

// QueryLocator encapsulates query results and provide additional methods to get All records or iterate over results using Next.
// Subselects are fields of type QueryLocator.
type QueryLocator interface {
	TotalSize() int
	QueryMore() (QueryLocator, error)
	All() ([]Record, error)
	Records() []Record
}

// QueryReader provides simple iterator over query results using Read.
type QueryReader struct {
	locator QueryLocator
	index   int
}

// NewReader creates new QueryReader starting with records in QueryLocator passed.
func NewReader(locator QueryLocator, errs ...error) (*QueryReader, error) {
	if len(errs) != 0 && errs[0] != nil {
		return nil, errs[0]
	}
	return &QueryReader{locator: locator, index: -1}, nil
}

// Read returns one record from QueryReader and advances to the next.
// Error could be returned when QueryMore should be called to get next QueryLocator.
func (reader *QueryReader) Read() (Record, error) {
	reader.index += 1
	records := reader.locator.Records()
	if reader.index < len(records) {
		return records[reader.index], nil
	}
	locator, err := reader.locator.QueryMore()
	if err != nil {
		return nil, err
	}
	reader.locator = locator
	reader.index = 0
	return reader.locator.Records()[reader.index], nil
}

// ReadAll returns all records queried
// Optional errs is passed through. Meant to be used as ReadAll(co.Query("..."))
func ReadAll(locator QueryLocator, errs ...error) ([]Record, error) {
	start := time.Now()
	defer func() { log.Println(DURATIONS, "ReadAll took:", time.Since(start)) }()
	if len(errs) != 0 && errs[0] != nil {
		return nil, errs[0]
	}
	return locator.All()
}

// EmptyQueryLocator is convinience type. Returned for subselects which render null result.
type EmptyQueryLocator struct{}

func (ql *EmptyQueryLocator) TotalSize() int                   { return 0 }
func (ql *EmptyQueryLocator) QueryMore() (QueryLocator, error) { return nil, io.EOF }
func (ql *EmptyQueryLocator) All() ([]Record, error)           { return make([]Record, 0, 0), nil }
func (ql *EmptyQueryLocator) Records() []Record                { return make([]Record, 0, 0) }
