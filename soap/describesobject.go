package soap

import (
	"github.com/goforce/log"
	"time"
)

func DescribeSObject(co Connection, name string) (*DescribeSObjectResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "DescribeSObject took:", time.Since(start)) }()
	dso := struct {
		XMLName struct{}               `xml:"Envelope`
		Result  *DescribeSObjectResult `xml:"Body>describeSObjectResponse>result"`
	}{}
	err := Call(co, &XE{"tns:describeSObject", XE{"tns:sObjectType", name}}, &dso)
	if err != nil {
		return nil, err
	}
	return dso.Result, nil
}

func DescribeSObjects(co Connection, names []string) ([]*DescribeSObjectResult, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "DescribeSObjects took:", time.Since(start)) }()
	ta := make([]XE, 0, len(names))
	for _, name := range names {
		ta = append(ta, XE{"tns:sObjectType", name})
	}
	dsos := struct {
		XMLName struct{}                 `xml:"Envelope`
		Result  []*DescribeSObjectResult `xml:"Body>describeSObjectsResponse>result"`
	}{}
	err := Call(co, &XE{"tns:describeSObjects", ta}, &dsos)
	if err != nil {
		return nil, err
	}
	return dsos.Result, nil
}

type DescribeSObjectResult struct {
	ActionOverrides     []*ActionOverride    `xml:"actionOverrides"`
	Activateable        bool                `xml:"activateable"`
	ChildRelationships  []*ChildRelationship `xml:"childRelationships"`
	CompactLayoutable   bool                 `xml:"compactLayoutable"`
	Createable          bool                 `xml:"createable"`
	Custom              bool                 `xml:"custom"`
	CustomSetting       bool                 `xml:"customSetting"`
	Deletable           bool                 `xml:"deletable"`
	DeprecatedAndHidden bool                 `xml:"deprecatedAndHidden"`
	FeedEnabled         bool                 `xml:"feedEnabled"`
	Fields              []*Field             `xml:"fields"`
	KeyPrefix           string               `xml:"keyPrefix"`
	Label               string               `xml:"label"`
	LabelPlural         string               `xml:"labelPlural"`
	Layoutable          bool                 `xml:"layoutable"`
	Mergeable           bool                 `xml:"mergeable"`
	Name                string               `xml:"name"`
	NamedLayoutInfos    []*NamedLayoutInfo   `xml:"namedLayoutInfos"`
	Queryable           bool                 `xml:"queryable"`
	RecordTypeInfos     []*RecordTypeInfo    `xml:"recordTypeInfos"`
	Replicateable       bool                 `xml:"replicateable"`
	Retrieveable        bool                 `xml:"retrieveable"`
	SearchLayoutable    bool                 `xml:"searchLayoutable"`
	Searchable          bool                 `xml:"searchable"`
	Triggerable         bool                 `xml:"triggerable"`
	Undeletable         bool                 `xml:"undeletable"`
	Updateable          bool                 `xml:"updateable"`
	UrlDetail           string               `xml:"urlDetail"`
	UrlEdit             string               `xml:"urlEdit"`
	UrlNew              string               `xml:"urlNew"`
}

type ActionOverride struct {
	IsAvailableInTouch bool   `xml:"isAvailableInTouch"`
	Name               string `xml:"name"`
	PageId             string `xml:"pageId"`
	Url                string `xml:"url"`
}

type ChildRelationship struct {
	CascadeDelete       bool     `xml:"cascadeDelete"`
	ChildSObject        string   `xml:"childSObject"`
	DeprecatedAndHidden bool     `xml:"deprecatedAndHidden"`
	Field               string   `xml:"field"`
	JunctionIdListName  string   `xml:"junctionIdListName"`
	JunctionReferenceTo []string `xml:"junctionReferenceTo"`
	RelationshipName    string   `xml:"relationshipName"`
	RestrictedDelete    bool     `xml:"restrictedDelete"`
}

type NamedLayoutInfo struct {
	Name string `xml:"name"`
}

type RecordTypeInfo struct {
	Available                bool   `xml:"available"`
	DefaultRecordTypeMapping bool   `xml:"defaultRecordTypeMapping"`
	Name                     string `xml:"name"`
	RecordTypeId             string `xml:"recordTypeId"`
}

type Field struct {
	AutoNumber               bool                `xml:"autoNumber"`
	ByteLength               int                 `xml:"byteLength"`
	Calculated               bool                `xml:"calculated"`
	CalculatedFormula        string              `xml:"calculatedFormula"`
	CascadeDelete            bool                `xml:"cascadeDelete"`
	CaseSensitive            bool                `xml:"caseSensitive"`
	ControllerName           string              `xml:"controllerName"`
	Createable               bool                `xml:"createable"`
	Custom                   bool                `xml:"custom"`
	DefaultValueFormula      string              `xml:"defaultValueFormula"`
	DefaultedOnCreate        bool                `xml:"defaultedOnCreate"`
	DependentPicklist        bool                `xml:"dependentPicklist"`
	DeprecatedAndHidden      bool                `xml:"deprecatedAndHidden"`
	Digits                   int                 `xml:"digits"`
	DisplayLocationInDecimal bool                `xml:"displayLocationInDecimal"`
	ExternalId               bool                `xml:"externalId"`
	ExtraTypeInfo            string              `xml:"extraTypeInfo"`
	Filterable               bool                `xml:"filterable"`
	FilteredLookupInfo       *FilteredLookupInfo `xml:"filteredLookupInfo"`
	Groupable                bool                `xml:"groupable"`
	HighScaleNumber          bool                `xml:"highScaleNumber"`
	HtmlFormatted            bool                `xml:"htmlFormatted"`
	IdLookup                 bool                `xml:"idLookup"`
	InlineHelpText           string              `xml:"inlineHelpText"`
	Label                    string              `xml:"label"`
	Length                   int                 `xml:"length"`
	Mask                     string              `xml:"mask"`
	MaskType                 string              `xml:"maskType"`
	Name                     string              `xml:"name"`
	NameField                bool                `xml:"nameField"`
	NamePointing             bool                `xml:"namePointing"`
	Nillable                 bool                `xml:"nillable"`
	Permissionable           bool                `xml:"permissionable"`
	PicklistValues           []*PicklistEntry    `xml:"picklistValues"`
	Precision                int                 `xml:"precision"`
	QueryByDistance          bool                `xml:"queryByDistance"`
	ReferenceTargetField     string              `xml:"referenceTargetField"`
	ReferenceTo              []string            `xml:"referenceTo"`
	RelationshipName         string              `xml:"relationshipName"`
	RelationshipOrder        int                 `xml:"relationshipOrder"`
	RestrictedDelete         bool                `xml:"restrictedDelete"`
	RestrictedPicklist       bool                `xml:"restrictedPicklist"`
	Scale                    int                 `xml:"scale"`
	SoapType                 string              `xml:"soapType"`
	Sortable                 bool                `xml:"sortable"`
	Type                     string              `xml:"type"`
	Unique                   bool                `xml:"unique"`
	Updateable               bool                `xml:"updateable"`
	WriteRequiresMasterRead  bool                `xml:"writeRequiresMasterRead"`
}

type FilteredLookupInfo struct {
	ControllingFields []string `xml:"controllingFields"`
	Dependent         bool     `xml:"dependent"`
	OptionalFilter    bool     `xml:"optionalFilter"`
}

type PicklistEntry struct {
	Active       bool   `xml:"active"`
	DefaultValue bool   `xml:"defaultValue"`
	Label        string `xml:"label"`
	ValidFor     string `xml:"validFor"`
	Value        string `xml:"value"`
}
