package soap

import (
	"strconv"
)

type SoapHeader interface {
	String() string
}

type SessionHeader struct {
	SessionId string
}

func (h *SessionHeader) String() string {
	return `<tns:SessionHeader><tns:sessionId>` +
		h.SessionId +
		`</tns:sessionId></tns:SessionHeader>`
}

type CallOptions struct {
	Client           string
	DefaultNamespace string
}

func (h *CallOptions) String() string {
	if len(h.Client) > 0 || len(h.DefaultNamespace) > 0 {
		s := `<tns:CallOptions>`
		if len(h.Client) > 0 {
			s += `<tns:client>` + h.Client + `</tns:client>`
		}
		if len(h.DefaultNamespace) > 0 {
			s += `<tns:defaultNamespace>` + h.DefaultNamespace + `</tns:defaultNamespace>`
		}
		s += `</tns:CallOptions>`
	}
	return ""
}

type AllOrNoneHeader struct {
	AllOrNone bool
}

func (h *AllOrNoneHeader) String() string {
	return `<tns:AllOrNoneHeader><tns:allOrNone>` +
		strconv.FormatBool(h.AllOrNone) +
		`</tns:allOrNone></tns:AllOrNoneHeader>`
}

type AllowFieldTruncationHeader struct {
	AllowFieldTruncation bool
}

func (h *AllowFieldTruncationHeader) String() string {
	return `<tns:AllowFieldTruncationHeader><tns:allowFieldTruncation>` +
		strconv.FormatBool(h.AllowFieldTruncation) +
		`</tns:allowFieldTruncation></tns:AllowFieldTruncationHeader>`
}

type AssignmentRuleHeader struct {
	AssignmentRuleId string
	UseDefaultRule   *bool
}

func (h *AssignmentRuleHeader) String() string {
	s := `<tns:AssignmentRuleHeader>`
	if h.AssignmentRuleId != "" {
		s = s + `<tns:assignmentRuleId>` + h.AssignmentRuleId + `</tns:assignmentRuleId>`
	} else if h.UseDefaultRule != nil {
		s = s + `<tns:useDefaultRule>` + strconv.FormatBool(*h.UseDefaultRule) + `</tns:useDefaultRule>`
	}
	s = s + `</tns:AssignmentRuleHeader>`
	return s
}

type DisableFeedTrackingHeader struct {
	DisableFeedTracking bool
}

func (h *DisableFeedTrackingHeader) String() string {
	return `<tns:DisableFeedTrackingHeader><tns:disableFeedTracking>` +
		strconv.FormatBool(h.DisableFeedTracking) +
		`</tns:disableFeedTracking></tns:DisableFeedTrackingHeader>`
}

type EmailHeader struct {
	TriggerAutoResponseEmail bool
	TriggerOtherEmail        bool
	TriggerUserEmail         bool
}

func (h *EmailHeader) String() string {
	return `<tns:EmailHeader><tns:triggerAutoResponseEmail>` +
		strconv.FormatBool(h.TriggerAutoResponseEmail) +
		`</tns:triggerAutoResponseEmail><tns:triggerOtherEmail>` +
		strconv.FormatBool(h.TriggerOtherEmail) +
		`</tns:triggerOtherEmail><tns:triggerUserEmail>` +
		strconv.FormatBool(h.TriggerUserEmail) +
		`<tns:triggerUserEmail></tns:EmailHeader>`
}

type LocaleOptions struct {
	language string
}

func (h *LocaleOptions) String() string {
	return `<tns:LocaleOptions><tns:language>` +
		h.language +
		`</tns:language></tns:LocaleOptions>`
}

type MruHeader struct {
	UpdateMru bool
}

func (h *MruHeader) String() string {
	return `<tns:MruHeader><tns:updateMru>` +
		strconv.FormatBool(h.UpdateMru) +
		`</tns:updateMru></tns:MruHeader>`
}

type OwnerChangeOptions struct {
	TransferAttachments    bool
	TransferOpenActivities bool
}

func (h *OwnerChangeOptions) String() string {
	return `<tns:OwnerChangeOptions><tns:transferAttachments>` +
		strconv.FormatBool(h.TransferAttachments) +
		`</tns:transferAttachments><tns:transferOpenActivities>` +
		strconv.FormatBool(h.TransferOpenActivities) +
		`<tns:transferOpenActivities></tns:OwnerChangeOptions>`
}

type QueryOptions struct {
	BatchSize int
}

func (h *QueryOptions) String() string {
	return `<tns:QueryOptions><tns:batchSize>` +
		strconv.Itoa(h.BatchSize) +
		`</tns:batchSize></tns:QueryOptions>`
}

type UserTerritoryDeleteHeader struct {
	TransferToUserId string
}

func (h *UserTerritoryDeleteHeader) String() string {
	return `<tns:UserTerritoryDeleteHeader><tns:batchSize>` +
		h.TransferToUserId +
		`</tns:batchSize></tns:UserTerritoryDeleteHeader>`
}
