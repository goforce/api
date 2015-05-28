package data

import (
	"github.com/goforce/api/soap"
	"strings"
)

type SObject struct {
	dr *soap.DescribeGlobalSObjectResult
}

func (co *Connection) SObjects() map[string]*SObject {
	if co.sobjects == nil {
		co.describeGlobal()
	}
	m := make(map[string]*SObject)
	for _, so := range co.sobjects {
		m[so.dr.Name] = so
	}
	return m
}

func (co *Connection) SObject(name string) (*SObject, bool) {
	n := strings.ToLower(name)
	if co.sobjects == nil {
		co.describeGlobal()
	}
	so, ok := co.sobjects[n]
	return so, ok
}

func (co *Connection) describeGlobal() {
	re, err := soap.DescribeGlobal(co)
	if err != nil {
		panic(err)
	}
	co.sobjects = make(map[string]*SObject)
	for _, d := range re.SObjects {
		co.sobjects[strings.ToLower(d.Name)] = &SObject{d}
	}
}

func (o *SObject) Activateable() bool        { return o.dr.Activateable }
func (o *SObject) Createable() bool          { return o.dr.Createable }
func (o *SObject) Custom() bool              { return o.dr.Custom }
func (o *SObject) CustomSetting() bool       { return o.dr.CustomSetting }
func (o *SObject) Deletable() bool           { return o.dr.Deletable }
func (o *SObject) DeprecatedAndHidden() bool { return o.dr.DeprecatedAndHidden }
func (o *SObject) FeedEnabled() bool         { return o.dr.FeedEnabled }
func (o *SObject) KeyPrefix() string         { return o.dr.KeyPrefix }
func (o *SObject) Label() string             { return o.dr.Label }
func (o *SObject) LabelPlural() string       { return o.dr.LabelPlural }
func (o *SObject) Layoutable() bool          { return o.dr.Layoutable }
func (o *SObject) Mergeable() bool           { return o.dr.Mergeable }
func (o *SObject) Name() string              { return o.dr.Name }
func (o *SObject) Queryable() bool           { return o.dr.Queryable }
func (o *SObject) Replicateable() bool       { return o.dr.Replicateable }
func (o *SObject) Retrieveable() bool        { return o.dr.Retrieveable }
func (o *SObject) Searchable() bool          { return o.dr.Searchable }
func (o *SObject) Triggerable() bool         { return o.dr.Triggerable }
func (o *SObject) Undeletable() bool         { return o.dr.Undeletable }
func (o *SObject) Updateable() bool          { return o.dr.Updateable }
