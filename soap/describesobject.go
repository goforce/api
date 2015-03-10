package soap

import (
	"github.com/goforce/api/soap/core"
	"strings"
)

type DescribeSObject struct {
	describe                 *core.DescribeSObjectResult
	childRelationships       []*ChildRelationship
	childRelationshipsByName map[string]*ChildRelationship
	fieldsByName             map[string]*Field
	relationshipsByName      map[string]*Field
}

func (describe *DescribeSObject) ChildRelationshipByName(name string) (*ChildRelationship, bool) {
	if describe.childRelationshipsByName == nil {
		describe.childRelationshipsByName = make(map[string]*ChildRelationship)
		for _, d := range describe.describe.ChildRelationships {
			describe.childRelationshipsByName[strings.ToLower(d.RelationshipName)] = &ChildRelationship{describe: d}
		}
	}
	d, ok := describe.childRelationshipsByName[strings.ToLower(name)]
	return d, ok
}

func (describe *DescribeSObject) FieldByName(name string) (*Field, bool) {
	if describe.fieldsByName == nil {
		describe.fieldsByName = make(map[string]*Field)
		for _, d := range describe.describe.Fields {
			describe.fieldsByName[strings.ToLower(d.Name)] = &Field{describe: d}
		}
	}
	d, ok := describe.fieldsByName[strings.ToLower(name)]
	return d, ok
}

func (describe *DescribeSObject) FieldByRelationshipName(name string) (*Field, bool) {
	if describe.relationshipsByName == nil {
		describe.relationshipsByName = make(map[string]*Field)
		for _, d := range describe.describe.Fields {
			if len(d.RelationshipName) > 0 {
				describe.relationshipsByName[strings.ToLower(d.RelationshipName)] = &Field{describe: d}
			}
		}
	}
	d, ok := describe.relationshipsByName[strings.ToLower(name)]
	return d, ok
}

func (co *Connection) DescribeSObject(name string) (*DescribeSObject, error) {
	n := strings.ToLower(name)
	if d, ok := co.describe[n]; ok {
		return d, nil
	}
	d, err := core.DescribeSObject(co, name)
	if err != nil {
		return nil, err
	}
	r := &DescribeSObject{describe: d}
	co.describe[strings.ToLower(d.Name)] = r
	return r, nil
}

func (co *Connection) WillUseSObjects(names []string) error {
	da, err := core.DescribeSObjects(co, names)
	if err != nil {
		return err
	}
	for _, d := range da {
		co.describe[strings.ToLower(d.Name)] = &DescribeSObject{describe: d}
	}
	return nil
}
