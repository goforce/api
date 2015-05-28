package data

import (
	"github.com/goforce/api/soap"
	"strings"
)

const sobject_aggregate_result = "aggregateresult"

type SObjectInfo struct {
	dr                 *soap.DescribeSObjectResult
	actionOverrides    *ActionOverrides
	childRelationships *ChildRelationships
	fields             *Fields
	namedLayoutInfos   *NamedLayoutInfos
	recordTypeInfos    *RecordTypeInfos
}

func (co *Connection) DescribeSObject(name string) (*SObjectInfo, error) {
	if info, ok := co.sobjectInfos[strings.ToLower(name)]; ok {
		return info, nil
	}
	dr, err := soap.DescribeSObject(co, name)
	if err != nil {
		return nil, err
	}
	info := &SObjectInfo{dr: dr}
	co.sobjectInfos[strings.ToLower(dr.Name)] = info
	return info, nil
}

func (co *Connection) DescribeSObjects(names ...string) (map[string]*SObjectInfo, error) {
	ns := make([]string, 0, len(names))
	for _, n := range names {
		n = strings.ToLower(n)
		if _, ok := co.sobjectInfos[n]; !ok {
			ns = append(ns, n)
		}
	}
	if len(ns) > 0 {
		drs, err := soap.DescribeSObjects(co, ns)
		if err != nil {
			return nil, err
		}
		for _, dr := range drs {
			co.sobjectInfos[strings.ToLower(dr.Name)] = &SObjectInfo{dr: dr}
		}
	}
	m := make(map[string]*SObjectInfo)
	for _, n := range names {
		m[n], _ = co.sobjectInfos[strings.ToLower(n)]
	}
	return m, nil
}

func (co *Connection) ResetSObjectInfos(names ...string) error {
	if len(names) == 0 {
		co.sobjectInfos = make(map[string]*SObjectInfo)
		return nil
	}
	ns := make([]string, 0, len(names))
	for _, n := range names {
		n = strings.ToLower(n)
		delete(co.sobjectInfos, n)
		ns = append(ns, n)
	}
	drs, err := soap.DescribeSObjects(co, ns)
	if err != nil {
		return err
	}
	for _, dr := range drs {
		co.sobjectInfos[strings.ToLower(dr.Name)] = &SObjectInfo{dr: dr}
	}
	return nil
}

func (oi *SObjectInfo) ActionOverrides() *ActionOverrides {
	if oi.actionOverrides == nil {
		oi.actionOverrides = newActionOverrides(oi.dr.ActionOverrides)
	}
	return oi.actionOverrides
}

func (oi *SObjectInfo) Activateable() bool { return oi.dr.Activateable }

func (oi *SObjectInfo) ChildRelationships() *ChildRelationships {
	if oi.childRelationships == nil {
		oi.childRelationships = newChildRelationships(oi.dr.ChildRelationships)
	}
	return oi.childRelationships
}

func (oi *SObjectInfo) CompactLayoutable() bool   { return oi.dr.CompactLayoutable }
func (oi *SObjectInfo) Createable() bool          { return oi.dr.Createable }
func (oi *SObjectInfo) Custom() bool              { return oi.dr.Custom }
func (oi *SObjectInfo) CustomSetting() bool       { return oi.dr.CustomSetting }
func (oi *SObjectInfo) Deletable() bool           { return oi.dr.Deletable }
func (oi *SObjectInfo) DeprecatedAndHidden() bool { return oi.dr.DeprecatedAndHidden }
func (oi *SObjectInfo) FeedEnabled() bool         { return oi.dr.FeedEnabled }

func (oi *SObjectInfo) Fields() *Fields {
	if oi.fields == nil {
		oi.fields = newFields(oi.dr.Fields)
	}
	return oi.fields
}

func (oi *SObjectInfo) KeyPrefix() string   { return oi.dr.KeyPrefix }
func (oi *SObjectInfo) Label() string       { return oi.dr.Label }
func (oi *SObjectInfo) LabelPlural() string { return oi.dr.LabelPlural }
func (oi *SObjectInfo) Layoutable() bool    { return oi.dr.Layoutable }
func (oi *SObjectInfo) Mergeable() bool     { return oi.dr.Mergeable }
func (oi *SObjectInfo) Name() string        { return oi.dr.Name }

func (oi *SObjectInfo) NamedLayoutInfos() *NamedLayoutInfos {
	if oi.namedLayoutInfos == nil {
		oi.namedLayoutInfos = newNamedLayoutInfos(oi.dr.NamedLayoutInfos)
	}
	return oi.namedLayoutInfos
}

func (oi *SObjectInfo) Queryable() bool { return oi.dr.Queryable }

func (oi *SObjectInfo) RecordTypeInfos() *RecordTypeInfos {
	if oi.recordTypeInfos == nil {
		oi.recordTypeInfos = newRecordTypeInfos(oi.dr.RecordTypeInfos)
	}
	return oi.recordTypeInfos
}

func (oi *SObjectInfo) Replicateable() bool    { return oi.dr.Replicateable }
func (oi *SObjectInfo) Retrieveable() bool     { return oi.dr.Retrieveable }
func (oi *SObjectInfo) SearchLayoutable() bool { return oi.dr.SearchLayoutable }
func (oi *SObjectInfo) Searchable() bool       { return oi.dr.Searchable }
func (oi *SObjectInfo) Triggerable() bool      { return oi.dr.Triggerable }
func (oi *SObjectInfo) Undeletable() bool      { return oi.dr.Undeletable }
func (oi *SObjectInfo) Updateable() bool       { return oi.dr.Updateable }
func (oi *SObjectInfo) UrlDetail() string      { return oi.dr.UrlDetail }
func (oi *SObjectInfo) UrlEdit() string        { return oi.dr.UrlEdit }
func (oi *SObjectInfo) UrlNew() string         { return oi.dr.UrlNew }

type NamedLayoutInfos struct {
	list []*soap.NamedLayoutInfo
}

func newNamedLayoutInfos(nlis []*soap.NamedLayoutInfo) *NamedLayoutInfos {
	return &NamedLayoutInfos{nlis}
}

func (nlis *NamedLayoutInfos) List() []string {
	re := make([]string, 0, len(nlis.list))
	for _, v := range nlis.list {
		re = append(re, v.Name)
	}
	return re
}

type RecordTypeInfos struct {
	byName map[string]*RecordTypeInfo
	byId   map[string]*RecordTypeInfo
}

func newRecordTypeInfos(rtis []*soap.RecordTypeInfo) *RecordTypeInfos {
	n := make(map[string]*RecordTypeInfo)
	i := make(map[string]*RecordTypeInfo)
	for _, v := range rtis {
		t := &RecordTypeInfo{v}
		n[strings.ToLower(v.Name)] = t
		i[v.RecordTypeId] = t
	}
	return &RecordTypeInfos{n, i}
}

func (rtis *RecordTypeInfos) List() []*RecordTypeInfo {
	re := make([]*RecordTypeInfo, 0, len(rtis.byName))
	for _, v := range rtis.byName {
		re = append(re, v)
	}
	return re
}

func (rtis *RecordTypeInfos) Get(name string) (*RecordTypeInfo, bool) {
	v, ok := rtis.byName[strings.ToLower(name)]
	return v, ok
}

func (rtis *RecordTypeInfos) GetById(id string) (*RecordTypeInfo, bool) {
	v, ok := rtis.byId[id]
	return v, ok
}

type RecordTypeInfo struct {
	dr *soap.RecordTypeInfo
}

func (rti *RecordTypeInfo) Available() bool { return rti.dr.Available }
func (rti *RecordTypeInfo) DefaultRecordTypeMapping() bool {
	return rti.dr.DefaultRecordTypeMapping
}
func (rti *RecordTypeInfo) Name() string         { return rti.dr.Name }
func (rti *RecordTypeInfo) RecordTypeId() string { return rti.dr.RecordTypeId }
