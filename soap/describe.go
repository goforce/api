package soap

import (
	"encoding/xml"
	"fmt"
	"github.com/goforce/api/commons"
	"github.com/goforce/log"
	"strings"
	"time"
)

var _ = fmt.Print

func (co *Connection) DescribeGlobal() (*commons.DescribeGlobalResult, error) {
	name := ""
	var r *commons.DescribeGlobalResult
	if v, ok := co.GetFromCache(name, r); ok {
		return v.(*commons.DescribeGlobalResult), nil
	}
	start := time.Now()
	defer func() { log.Println(commons.DURATIONS, "DescribeGlobal took:", time.Since(start)) }()
	req := []KeyValue{
		KeyValue{"tns:describeGlobal", ""},
	}
	reader, err := co.Call(req)
	if reader != nil {
		defer reader.Close()
	}
	if err != nil {
		return nil, err
	}
	dg := struct {
		XMLName struct{}                      `xml:"Envelope`
		Result  *commons.DescribeGlobalResult `xml:"Body>describeGlobalResponse>result"`
	}{}
	err = xml.NewDecoder(reader).Decode(&dg)
	if err != nil {
		return nil, err
	}
	co.AddToCache(name, dg.Result)
	return dg.Result, nil
}

func (co *Connection) DescribeSObject(name string) (*commons.DescribeSObjectResult, error) {
	var r *commons.DescribeSObjectResult
	if v, ok := co.GetFromCache(strings.ToLower(name), r); ok {
		return v.(*commons.DescribeSObjectResult), nil
	}
	start := time.Now()
	defer func() { log.Println(commons.DURATIONS, "DescribeSObject took:", time.Since(start)) }()
	req := []KeyValue{
		KeyValue{"tns:describeSObject/tns:sObjectType", name},
	}
	response, err := co.Call(req)
	if err != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	dso := struct {
		XMLName struct{}                       `xml:"Envelope`
		Result  *commons.DescribeSObjectResult `xml:"Body>describeSObjectResponse>result"`
	}{}
	err = xml.NewDecoder(response).Decode(&dso)
	if err != nil {
		return nil, err
	}
	co.AddToCache(strings.ToLower(name), dso.Result)
	return dso.Result, nil
}

func (co *Connection) DescribeSObjects(names []string) ([]*commons.DescribeSObjectResult, error) {
	ar := make([]*commons.DescribeSObjectResult, 0, len(names))
	qn := make([]string, 0, len(names))
	for _, name := range names {
		var r *commons.DescribeSObjectResult
		if v, ok := co.GetFromCache(strings.ToLower(name), r); ok {
			ar = append(ar, v.(*commons.DescribeSObjectResult))
		} else {
			qn = append(qn, name)
		}
	}
	if len(qn) == 0 {
		return ar, nil
	}
	start := time.Now()
	defer func() { log.Println(commons.DURATIONS, "DescribeSObjects took:", time.Since(start)) }()
	req := make([]KeyValue, 0, len(names))
	for _, name := range qn {
		req = append(req, KeyValue{"tns:describeSObjects/tns:sObjectType", name})
	}
	response, err := co.Call(req)
	if err != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	dsos := struct {
		XMLName struct{}                         `xml:"Envelope`
		Result  []*commons.DescribeSObjectResult `xml:"Body>describeSObjectsResponse>result"`
	}{}
	err = xml.NewDecoder(response).Decode(&dsos)
	if err != nil {
		return nil, err
	}
	for _, r := range dsos.Result {
		co.AddToCache(strings.ToLower(r.Name), r)
	}
	ar = append(ar, dsos.Result...)
	return ar, nil
}
