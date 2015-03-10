package soap

import (
	"github.com/goforce/api/soap/core"
	"strings"
)

type SObject struct {
	describe *core.DescribeGlobalSObjectResult
}

func (co *Connection) SObjects() ([]*SObject, error) {
	if co.sobjects == nil {
		err := co.describeGlobal()
		if err != nil {
			return nil, err
		}
	}
	r := make([]*SObject, 0, len(co.sobjects))
	for _, so := range co.sobjects {
		r = append(r, so)
	}
	return r, nil
}

func (co *Connection) SObject(name string) (*SObject, bool) {
	n := strings.ToLower(name)
	if co.sobjects == nil {
		err := co.describeGlobal()
		if err != nil {
			return nil, false
		}
	}
	return co.sobjects[n], true
}

func (co *Connection) describeGlobal() error {
	gd, err := core.DescribeGlobal(co)
	if err != nil {
		return err
	}
	co.sobjects = make(map[string]*SObject)
	for _, d := range gd.SObjects {
		co.sobjects[strings.ToLower(d.Name)] = &SObject{d}
	}
	return nil
}

func (o *SObject) Name() string { return o.describe.Name }
