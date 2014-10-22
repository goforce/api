package commons

import (
	"strings"
)

type DescribeGlobalResult struct {
	Encoding     string                         `xml:"encoding" json:"encoding"`
	MaxBatchSize int                            `xml:"maxBatchSize" json:"maxBatchSize"`
	SObjects     []*DescribeGlobalSObjectResult `xml:"sobjects" json:"sobjects"`
	sObjectsMap  map[string]*DescribeGlobalSObjectResult
}

func (dgr *DescribeGlobalResult) Get(sObject string) *DescribeGlobalSObjectResult {
	if dgr.sObjectsMap == nil {
		dgr.sObjectsMap = make(map[string]*DescribeGlobalSObjectResult)
		for _, d := range dgr.SObjects {
			dgr.sObjectsMap[strings.ToLower(d.Name)] = d
		}
	}
	return dgr.sObjectsMap[strings.ToLower(sObject)]
}

type DescribeGlobalSObjectResult struct {
	Activateable        bool        `xml:"activateable" json:"activateable"`
	Createable          bool        `xml:"createable" json:"createable"`
	Custom              bool        `xml:"custom" json:"custom"`
	CustomSetting       bool        `xml:"customSetting" json:"customSetting"`
	Deletable           bool        `xml:"deletable" json:"deletable"`
	DeprecatedAndHidden bool        `xml:"deprecatedAndHidden" json:"deprecatedAndHidden"`
	FeedEnabled         bool        `xml:"feedEnabled" json:"feedEnabled"`
	KeyPrefix           string      `xml:"keyPrefix" json:"keyPrefix"`
	Label               string      `xml:"label" json:"label"`
	LabelPlural         string      `xml:"labelPlural" json:"labelPlural"`
	Layoutable          bool        `xml:"layoutable" json:"layoutable"`
	Mergeable           bool        `xml:"mergeable" json:"mergeable"`
	Name                string      `xml:"name" json:"name"`
	Queryable           bool        `xml:"queryable" json:"queryable"`
	Replicateable       bool        `xml:"replicateable" json:"replicateable"`
	Retrieveable        bool        `xml:"retrieveable" json:"retrieveable"`
	Searchable          bool        `xml:"searchable" json:"searchable"`
	Triggerable         bool        `xml:"triggerable" json:"triggerable"`
	Undeletable         bool        `xml:"undeletable" json:"undeletable"`
	Updateable          bool        `xml:"updateable" json:"updateable"`
	Urls                SObjectUrls `json:"urls"`
}

type SObjectUrls struct {
	SObject         string `json:"sobject"`
	ApprovalLayouts string `json:"approvalLayouts"`
	QuickActions    string `json:"quickActions"`
	Describe        string `json:"describe"`
	RowTemplate     string `json:"rowTemplate"`
	Layouts         string `json:"layouts"`
	CompactLayouts  string `json:"compactLayouts"`
}

type DescribeSObjectResult struct {
	ActionOverride []struct {
		IsAvailableInTouch bool   `xml:"isAvailableInTouch"`
		Name               string `xml:"name"`
		PageId             string `xml:"pageId"`
		Url                string `xml:"url"`
	} `xml:"actionOverrides"`
	Activateable        bool                          `xml:"activateable" json:"activateable"`
	ChildRelationships  []*ChildRelationship          `xml:"childRelationships" json:"childRelationships"`
	CompactLayoutable   bool                          `xml:"compactLayoutable" json:"compactLayoutable"`
	Createable          bool                          `xml:"createable" json:"createable"`
	Custom              bool                          `xml:"custom" json:"custom"`
	CustomSetting       bool                          `xml:"customSetting" json:"customSetting"`
	Deletable           bool                          `xml:"deletable" json:"deletable"`
	DeprecatedAndHidden bool                          `xml:"deprecatedAndHidden" json:"deprecatedAndHidden"`
	FeedEnabled         bool                          `xml:"feedEnabled" json:"feedEnabled"`
	Fields              []*DescribeSObjectFieldResult `xml:"fields" json:"fields"`
	KeyPrefix           string                        `xml:"keyPrefix" json:"keyPrefix"`
	Label               string                        `xml:"label" json:"label"`
	LabelPlural         string                        `xml:"labelPlural" json:"labelPlural"`
	Layoutable          bool                          `xmk:"layoutable" json:"layoutable"`
	LookupLayoutable    string                        `json:"lookupLayoutable"`
	Listviewable        string                        `json:"listviewable"`
	Mergeable           bool                          `xml:"mergeable" json:"mergeable"`
	Name                string                        `xml:"name" json:"name"`
	NamedLayoutInfos    []struct {
		Name string `xml:"name" json:"name"`
	} `xml:"namedLayoutInfos" json:"namedLayoutInfos"`
	Queryable       bool `xml:"queryable" json:"queryable"`
	RecordTypeInfos []*struct {
		Available                bool   `xml:"available" json:"available"`
		DefaultRecordTypeMapping bool   `xml:"defaultRecordTypeMapping" json:"defaultRecordTypeMapping"`
		Name                     string `xml:"name" json:"name"`
		RecordTypeId             string `xml:"recordTypeId" json:"recordTypeId"`
		Urls                     struct {
			Layout string `json:"layout"`
		} `json:"urls"`
	} `xml:"recordTypeInfos" json:"recordTypeInfos"`
	Replicateable    bool   `xml:"replicateable" json:"replicateable"`
	Retrieveable     bool   `xml:"retrieveable" json:"retrieveable"`
	SearchLayoutable bool   `xml:"searchLayoutable" json:"searchLayoutable"`
	Searchable       bool   `xml:"searchable" json:"searchable"`
	Triggerable      bool   `xml:"triggerable" json:"triggerable"`
	Undeletable      bool   `xml:"undeletable" json:"undeletable"`
	Updateable       bool   `xml:"updateable" json:"updateable"`
	UrlDetail        string `xml:"urlDetail"`
	UrlEdit          string `xml:"urlEdit"`
	UrlNew           string `xml:"urlNew"`
	Urls             struct {
		UiEditTemplate   string `json:"uiEditTemplate"`
		SObject          string `json:"sobject"`
		QuickActions     string `json:"quickActions"`
		UiDetailTemplate string `json:"uiDetailTemplate"`
		Describe         string `json:"describe"`
		RowTemplate      string `json:"rowTemplate"`
		Layouts          string `json:"layouts"`
		cCompactLayouts  string `json:"compactLayouts"`
		UiNewRecord      string `json:"uiNewRecord"`
	} `json:"urls"`
	fields             map[string]*DescribeSObjectFieldResult
	relationships      map[string]*DescribeSObjectFieldResult
	childRelationships map[string]*ChildRelationship
}

func (dsor *DescribeSObjectResult) Get(field string) *DescribeSObjectFieldResult {
	if dsor.fields == nil {
		dsor.fields = make(map[string]*DescribeSObjectFieldResult)
		for _, d := range dsor.Fields {
			dsor.fields[strings.ToLower(d.Name)] = d
		}
	}
	return dsor.fields[strings.ToLower(field)]
}

func (dsor *DescribeSObjectResult) GetRelationship(name string) *DescribeSObjectFieldResult {
	if dsor.relationships == nil {
		dsor.relationships = make(map[string]*DescribeSObjectFieldResult)
		for _, d := range dsor.Fields {
			if d.RelationshipName != "" {
				dsor.relationships[strings.ToLower(d.RelationshipName)] = d
			}
		}
	}
	return dsor.relationships[strings.ToLower(name)]
}

func (dsor *DescribeSObjectResult) GetChildRelationship(name string) *ChildRelationship {
	if dsor.childRelationships == nil {
		dsor.childRelationships = make(map[string]*ChildRelationship)
		for _, d := range dsor.ChildRelationships {
			dsor.childRelationships[strings.ToLower(d.RelationshipName)] = d
		}
	}
	return dsor.childRelationships[strings.ToLower(name)]
}

type ChildRelationship struct {
	CascadeDelete       bool   `xml:"cascadeDelete" json:"cascadeDelete"`
	ChildSObject        string `xml:"childSObject" json:"childSObject"`
	DeprecatedAndHidden bool   `xml:"deprecatedAndHidden" json:"deprecatedAndHidden"`
	Field               string `xml:"field" json:"field"`
	RelationshipName    string `xml:"relationshipName" json:"relationshipName"`
	RestrictedDelete    bool   `xml:"restrictedDelete" json:"restrictedDelete"`
}

type DescribeSObjectFieldResult struct {
	AutoNumber               bool   `xml:"autoNumber" json:"autoNumber"`
	ByteLength               int    `xml:"byteLength" json:"byteLength"`
	Calculated               bool   `xml:"calculated" json:"calculated"`
	CalculatedFormula        string `xml:"calculatedFormula" json:"calculatedFormula"`
	CascadeDelete            bool   `xml:"cascadeDelete" json:"cascadeDelete"`
	CaseSensitive            bool   `xml:"caseSensitive" json:"caseSensitive"`
	ControllerName           string `xml:"controllerName" json:"controllerName"`
	Createable               bool   `xml:"createable" json:"createable"`
	Custom                   bool   `xml:"custom" json:"custom"`
	DefaultValue             string `json:"defaultValue"` // TODO what's the equivalent in xml
	DefaultValueFormula      string `xml:"defaultValueFormula" json:"defaultValueFormula"`
	DefaultedOnCreate        bool   `xml:"defaultedOnCreate" json:"defaultedOnCreate"`
	DependentPicklist        bool   `xml:"dependentPicklist" json:"dependentPicklist"`
	DeprecatedAndHidden      bool   `xml:"deprecatedAndHidden" json:"deprecatedAndHidden"`
	Digits                   int    `xml:"digits" json:"digits"`
	DisplayLocationInDecimal bool   `xml:"displayLocationInDecimal" json:"displayLocationInDecimal"`
	ExternalId               bool   `xml:"externalId" json:"externalId"`
	ExtraTypeInfo            string `xml:"extraTypeInfo"`
	Filterable               bool   `xml:"filterable" json:"filterable"`
	FilteredLookupInfo       struct {
		ControllingFields []string `xml:"controllingFields" json:"controllingFields"`
		Dependent         bool     `xml:"dependent" json:"dependent"`
		OptionalFilter    bool     `xml:"optionalFilter" json:"optionalFilter"`
	} `xml:"filteredLookupInfo" json:"filteredLookupInfo"`
	Groupable               bool             `xml:"groupable" json:"groupable"`
	HtmlFormatted           bool             `xml:"htmlFormatted" json:"htmlFormatted"`
	IdLookup                bool             `xml:"idLookup" json:"idLookup"`
	InlineHelpText          string           `xml:"inlineHelpText" json:"inlineHelpText"`
	Label                   string           `xml:"label" json:"label"`
	Length                  int              `xml:"length" json:"length"`
	Mask                    string           `xml:"mask" json:"mask"`
	MaskType                string           `xml:"maskType" json:"maskType"`
	Name                    string           `xml:"name" json:"name"`
	NameField               bool             `xml:"nameField" json:"nameField"`
	NamePointing            bool             `xml:"namePointing" json:"namePointing"`
	Nillable                bool             `xml:"nillable" json:"nillable"`
	Permissionable          bool             `xml:"permissionable" json:"permissionable"`
	PicklistValues          []*PicklistEntry `xml:"picklistValues" json:"picklistValues"`
	Precision               int              `xml:"precision" json:"precision"`
	QueryByDistance         bool             `xml:"queryByDistance"`
	ReferenceTargetField    string           `xml:"referenceTargetField"`
	ReferenceTo             []string         `xml:"referenceTo" json:"referenceTo"`
	RelationshipName        string           `xml:"relationshipName" json:"relationshipName"`
	RelationshipOrder       int              `xml:"relationshipOrder" json:"relationshipOrder"`
	RestrictedDelete        bool             `xml:"restrictedDelete" json:"restrictedDelete"`
	RestrictedPicklist      bool             `xml:"restrictedPicklist" json:"restrictedPicklist"`
	Scale                   int              `xml:"scale" json:"scale"`
	SoapType                string           `xml:"soapType" json:"soapType"`
	Sortable                bool             `xml:"sortable" json:"sortable"`
	Type                    string           `xml:"type" json:"type"`
	Unique                  bool             `xml:"unique" json:"unique"`
	Updateable              bool             `xml:"updateable" json:"updateable"`
	WriteRequiresMasterRead bool             `xml:"writeRequiresMasterRead" json:"writeRequiresMasterRead"`
}

type PicklistEntry struct {
	Active       bool   `xml:"active" json:"active"`
	DefaultValue bool   `xml:"defaultValue" json:"defaultValue"`
	Label        string `xml:"label" json:"label"`
	ValidFor     string `xml:"validFor" json:"validFor"`
	Value        string `xml:"value" json:"value"`
}
