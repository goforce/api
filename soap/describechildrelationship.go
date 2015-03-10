package soap

import (
	"github.com/goforce/api/soap/core"
)

type ChildRelationship struct {
	describe *core.ChildRelationship
}

func newChildRelationship(describe *core.ChildRelationship) *ChildRelationship {
	return &ChildRelationship{describe: describe}
}

func (describe *ChildRelationship) CascadeDelete() bool  { return describe.describe.CascadeDelete }
func (describe *ChildRelationship) ChildSObject() string { return describe.describe.ChildSObject }
func (describe *ChildRelationship) DeprecatedAndHidden() bool {
	return describe.describe.DeprecatedAndHidden
}
func (describe *ChildRelationship) Field() string { return describe.describe.Field }
func (describe *ChildRelationship) JunctionIdListName() string {
	return describe.describe.JunctionIdListName
}
func (describe *ChildRelationship) JunctionReferenceTo() []string {
	r := make([]string, len(describe.describe.JunctionReferenceTo))
	copy(r, describe.describe.JunctionReferenceTo)
	return r
}
func (describe *ChildRelationship) RelationshipName() string {
	return describe.describe.RelationshipName
}
func (describe *ChildRelationship) RestrictedDelete() bool {
	return describe.describe.RestrictedDelete
}
