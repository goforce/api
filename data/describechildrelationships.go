package data

import (
	"github.com/goforce/api/soap"
	"strings"
)

type ChildRelationships struct {
	byName map[string]*ChildRelationship
}

func newChildRelationships(crs []*soap.ChildRelationship) *ChildRelationships {
	m := make(map[string]*ChildRelationship)
	for _, cr := range crs {
		m[strings.ToLower(cr.RelationshipName)] = &ChildRelationship{cr}
	}
	return &ChildRelationships{m}
}

func (crs *ChildRelationships) List() []*ChildRelationship {
	l := make([]*ChildRelationship, 0, len(crs.byName))
	for _, cr := range crs.byName {
		l = append(l, cr)
	}
	return l
}

func (crs *ChildRelationships) Get(name string) (*ChildRelationship, bool) {
	v, ok := crs.byName[strings.ToLower(name)]
	return v, ok
}

type ChildRelationship struct {
	dr *soap.ChildRelationship
}

func (cr *ChildRelationship) CascadeDelete() bool        { return cr.dr.CascadeDelete }
func (cr *ChildRelationship) ChildSObject() string       { return cr.dr.ChildSObject }
func (cr *ChildRelationship) DeprecatedAndHidden() bool  { return cr.dr.DeprecatedAndHidden }
func (cr *ChildRelationship) Field() string              { return cr.dr.Field }
func (cr *ChildRelationship) JunctionIdListName() string { return cr.dr.JunctionIdListName }
func (cr *ChildRelationship) JunctionReferenceTo() []string {
	r := make([]string, len(cr.dr.JunctionReferenceTo))
	copy(r, cr.dr.JunctionReferenceTo)
	return r
}
func (cr *ChildRelationship) RelationshipName() string { return cr.dr.RelationshipName }
func (cr *ChildRelationship) RestrictedDelete() bool   { return cr.dr.RestrictedDelete }
