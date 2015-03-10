package core

import (
	"github.com/goforce/api/commons"
	"github.com/goforce/log"
	"time"
)

func DescribeGlobal(co commons.Connection) (*DescribeGlobalResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "DescribeGlobal took:", time.Since(start)) }()
	dg := struct {
		XMLName struct{}              `xml:"Envelope`
		Result  *DescribeGlobalResult `xml:"Body>describeGlobalResponse>result"`
	}{}
	err := Call(co, &XE{"tns:describeGlobal", ""}, &dg)
	if err != nil {
		return nil, err
	}
	return dg.Result, nil
}

type DescribeGlobalResult struct {
	Encoding     string                         `xml:"encoding"`
	MaxBatchSize int                            `xml:"maxBatchSize"`
	SObjects     []*DescribeGlobalSObjectResult `xml:"sobjects"`
}

type DescribeGlobalSObjectResult struct {
	Activateable        bool   `xml:"activateable"`
	Createable          bool   `xml:"createable"`
	Custom              bool   `xml:"custom"`
	CustomSetting       bool   `xml:"customSetting"`
	Deletable           bool   `xml:"deletable"`
	DeprecatedAndHidden bool   `xml:"deprecatedAndHidden"`
	FeedEnabled         bool   `xml:"feedEnabled"`
	KeyPrefix           string `xml:"keyPrefix"`
	Label               string `xml:"label"`
	LabelPlural         string `xml:"labelPlural"`
	Layoutable          bool   `xml:"layoutable"`
	Mergeable           bool   `xml:"mergeable"`
	Name                string `xml:"name"`
	Queryable           bool   `xml:"queryable"`
	Replicateable       bool   `xml:"replicateable"`
	Retrieveable        bool   `xml:"retrieveable"`
	Searchable          bool   `xml:"searchable"`
	Triggerable         bool   `xml:"triggerable"`
	Undeletable         bool   `xml:"undeletable"`
	Updateable          bool   `xml:"updateable"`
}
