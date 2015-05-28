package data

import (
	"github.com/goforce/api/soap"
	"strings"
)

type Fields struct {
	byName             map[string]*Field
	byRelationshipName map[string]*Field
}

func newFields(fs []*soap.Field) *Fields {
	n := make(map[string]*Field)
	r := make(map[string]*Field)
	for _, f := range fs {
		d := &Field{dr: f}
		n[strings.ToLower(f.Name)] = d
		if len(f.RelationshipName) > 0 {
			r[strings.ToLower(f.RelationshipName)] = d
		}
	}
	return &Fields{byName: n, byRelationshipName: r}
}

func (fs *Fields) List() []*Field {
	l := make([]*Field, 0, len(fs.byName))
	for _, f := range fs.byName {
		l = append(l, f)
	}
	return l
}

func (fs *Fields) Get(name string) (*Field, bool) {
	v, ok := fs.byName[strings.ToLower(name)]
	return v, ok
}

func (fs *Fields) GetByRelationshipName(name string) (*Field, bool) {
	v, ok := fs.byRelationshipName[strings.ToLower(name)]
	return v, ok
}

type Field struct {
	dr                 *soap.Field
	fieldType          FieldType
	filteredLookupInfo *FilteredLookupInfo
	picklistValues     *PicklistEntries
}

func (df *Field) AutoNumber() bool               { return df.dr.AutoNumber }
func (df *Field) ByteLength() int                { return df.dr.ByteLength }
func (df *Field) Calculated() bool               { return df.dr.Calculated }
func (df *Field) CalculatedFormula() string      { return df.dr.CalculatedFormula }
func (df *Field) CascadeDelete() bool            { return df.dr.CascadeDelete }
func (df *Field) CaseSensitive() bool            { return df.dr.CaseSensitive }
func (df *Field) ControllerName() string         { return df.dr.ControllerName }
func (df *Field) Createable() bool               { return df.dr.Createable }
func (df *Field) Custom() bool                   { return df.dr.Custom }
func (df *Field) DefaultValueFormula() string    { return df.dr.DefaultValueFormula }
func (df *Field) DefaultedOnCreate() bool        { return df.dr.DefaultedOnCreate }
func (df *Field) DependentPicklist() bool        { return df.dr.DependentPicklist }
func (df *Field) DeprecatedAndHidden() bool      { return df.dr.DeprecatedAndHidden }
func (df *Field) Digits() int                    { return df.dr.Digits }
func (df *Field) DisplayLocationInDecimal() bool { return df.dr.DisplayLocationInDecimal }
func (df *Field) ExternalId() bool               { return df.dr.ExternalId }
func (df *Field) ExtraTypeInfo() string          { return df.dr.ExtraTypeInfo }
func (df *Field) Filterable() bool               { return df.dr.Filterable }

func (df *Field) FilteredLookupInfo() *FilteredLookupInfo {
	if df.filteredLookupInfo == nil {
		df.filteredLookupInfo = &FilteredLookupInfo{df.dr.FilteredLookupInfo}
	}
	return df.filteredLookupInfo
}

func (df *Field) Groupable() bool        { return df.dr.Groupable }
func (df *Field) HighScaleNumber() bool  { return df.dr.HighScaleNumber }
func (df *Field) HtmlFormatted() bool    { return df.dr.HtmlFormatted }
func (df *Field) IdLookup() bool         { return df.dr.IdLookup }
func (df *Field) InlineHelpText() string { return df.dr.InlineHelpText }
func (df *Field) Label() string          { return df.dr.Label }
func (df *Field) Length() int            { return df.dr.Length }
func (df *Field) Mask() string           { return df.dr.Mask }
func (df *Field) MaskType() string       { return df.dr.MaskType }
func (df *Field) Name() string           { return df.dr.Name }
func (df *Field) NameField() bool        { return df.dr.NameField }
func (df *Field) NamePointing() bool     { return df.dr.NamePointing }
func (df *Field) Nillable() bool         { return df.dr.Nillable }
func (df *Field) Permissionable() bool   { return df.dr.Permissionable }

func (df *Field) PicklistValues() *PicklistEntries {
	if df.picklistValues == nil {
		df.picklistValues = newPicklistEntries(df.dr.PicklistValues)
	}
	return df.picklistValues
}

func (df *Field) Precision() int               { return df.dr.Precision }
func (df *Field) QueryByDistance() bool        { return df.dr.QueryByDistance }
func (df *Field) ReferenceTargetField() string { return df.dr.ReferenceTargetField }

func (df *Field) ReferenceTo() []string {
	r := make([]string, len(df.dr.ReferenceTo))
	copy(r, df.dr.ReferenceTo)
	return r
}

func (df *Field) RelationshipName() string { return df.dr.RelationshipName }
func (df *Field) RelationshipOrder() int   { return df.dr.RelationshipOrder }
func (df *Field) RestrictedDelete() bool   { return df.dr.RestrictedDelete }
func (df *Field) RestrictedPicklist() bool { return df.dr.RestrictedPicklist }
func (df *Field) Scale() int               { return df.dr.Scale }
func (df *Field) SoapType() string         { return df.dr.SoapType }
func (df *Field) Sortable() bool           { return df.dr.Sortable }

func (df *Field) Type() FieldType {
	if df.fieldType == 0 {
		df.fieldType = getFieldType(df.dr.Type)
	}
	return df.fieldType
}

func (df *Field) Unique() bool                  { return df.dr.Unique }
func (df *Field) Updateable() bool              { return df.dr.Updateable }
func (df *Field) WriteRequiresMasterRead() bool { return df.dr.WriteRequiresMasterRead }

type FilteredLookupInfo struct {
	dr *soap.FilteredLookupInfo
}

func (fli *FilteredLookupInfo) ControllingFields() []string {
	r := make([]string, len(fli.dr.ControllingFields))
	copy(r, fli.dr.ControllingFields)
	return r
}
func (fli *FilteredLookupInfo) Dependent() bool      { return fli.dr.Dependent }
func (fli *FilteredLookupInfo) OptionalFilter() bool { return fli.dr.OptionalFilter }

type PicklistEntries struct {
	list []*PicklistEntry
}

func newPicklistEntries(pes []*soap.PicklistEntry) *PicklistEntries {
	l := make([]*PicklistEntry, 0, len(pes))
	for _, pe := range pes {
		l = append(l, &PicklistEntry{pe})
	}
	return &PicklistEntries{l}
}

func (pes *PicklistEntries) List() []*PicklistEntry {
	r := make([]*PicklistEntry, len(pes.list))
	copy(r, pes.list)
	return r
}

type PicklistEntry struct {
	dr *soap.PicklistEntry
}

func (pe *PicklistEntry) Active() bool       { return pe.dr.Active }
func (pe *PicklistEntry) DefaultValue() bool { return pe.dr.DefaultValue }
func (pe *PicklistEntry) Label() string      { return pe.dr.Label }
func (pe *PicklistEntry) ValidFor() string   { return pe.dr.ValidFor }
func (pe *PicklistEntry) Value() string      { return pe.dr.Value }
