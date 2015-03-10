package soap

import (
	"github.com/goforce/api/soap/core"
	"strconv"
)

type SoapHeader interface {
	xe() *core.XE
}

type CallOptions struct {
	Client           string
	DefaultNamespace string
}

func (h *CallOptions) xe() *core.XE {
	if len(h.Client) == 0 && len(h.DefaultNamespace) == 0 {
		return nil
	}
	x := make([]XE, 0, 2)
	if len(h.Client) > 0 {
		x = append(x, XE{"tns:client", h.Client})
	}
	if len(h.DefaultNamespace) > 0 {
		x = append(x, XE{"tns:defaultNamespace", h.DefaultNamespace})
	}
	return &core.XE{"tns:CallOptions", x}
}

type AllOrNoneHeader struct {
	AllOrNone bool
}

func (h *AllOrNoneHeader) xe() *core.XE {
	return &core.XE{"tns:AllOrNoneHeader", XE{"tns:allOrNone", strconv.FormatBool(h.AllOrNone)}}
}

type AllowFieldTruncationHeader struct {
	AllowFieldTruncation bool
}

func (h *AllowFieldTruncationHeader) xe() *core.XE {
	return &core.XE{"tns:AllowFieldTruncationHeader", XE{"tns:allowFieldTruncation", strconv.FormatBool(h.AllowFieldTruncation)}}
}

type AssignmentRuleHeader struct {
	AssignmentRuleId string
	UseDefaultRule   *bool
}

func (h *AssignmentRuleHeader) xe() *core.XE {
	if len(h.AssignmentRuleId) == 0 && h.UseDefaultRule == nil {
		return nil
	}
	x := make([]XE, 0, 2)
	if len(h.AssignmentRuleId) > 0 {
		x = append(x, XE{"tns:assignmentRuleId", h.AssignmentRuleId})
	}
	if h.UseDefaultRule != nil {
		x = append(x, XE{"tns:useDefaultRule", strconv.FormatBool(*h.UseDefaultRule)})
	}
	return &core.XE{"tns:AssignmentRuleHeader", x}
}

type DisableFeedTrackingHeader struct {
	DisableFeedTracking bool
}

func (h *DisableFeedTrackingHeader) xe() *core.XE {
	return &core.XE{"tns:DisableFeedTrackingHeader", XE{"tns:disableFeedTracking", strconv.FormatBool(h.DisableFeedTracking)}}
}

type EmailHeader struct {
	TriggerAutoResponseEmail bool
	TriggerOtherEmail        bool
	TriggerUserEmail         bool
}

func (h *EmailHeader) xe() *core.XE {
	x := make([]XE, 0, 3)
	x = append(x, XE{"tns:triggerAutoResponseEmail", strconv.FormatBool(h.TriggerAutoResponseEmail)})
	x = append(x, XE{"tns:triggerOtherEmail", strconv.FormatBool(h.TriggerOtherEmail)})
	x = append(x, XE{"tns:triggerUserEmail", strconv.FormatBool(h.TriggerUserEmail)})
	return &core.XE{"tns:EmailHeader", x}
}

type LocaleOptions struct {
	Language string
}

func (h *LocaleOptions) xe() *core.XE {
	return &core.XE{"tns:LocaleOptions", XE{"tns:language", h.Language}}
}

type MruHeader struct {
	UpdateMru bool
}

func (h *MruHeader) xe() *core.XE {
	return &core.XE{"tns:MruHeader", XE{"tns:updateMru", strconv.FormatBool(h.UpdateMru)}}
}

type OwnerChangeOptions struct {
	TransferAttachments    bool
	TransferOpenActivities bool
}

func (h *OwnerChangeOptions) xe() *core.XE {
	x := make([]XE, 0, 2)
	x = append(x, XE{"tns:transferAttachments", strconv.FormatBool(h.TransferAttachments)})
	x = append(x, XE{"tns:transferOpenActivities", strconv.FormatBool(h.TransferOpenActivities)})
	return &core.XE{"tns:OwnerChangeOptions", x}
}

type QueryOptions struct {
	BatchSize int
}

func (h *QueryOptions) xe() *core.XE {
	return &core.XE{"tns:QueryOptions", XE{"tns:batchSize", strconv.Itoa(h.BatchSize)}}
}

type UserTerritoryDeleteHeader struct {
	TransferToUserId string
}

func (h *UserTerritoryDeleteHeader) xe() *core.XE {
	return &core.XE{"tns:UserTerritoryDeleteHeader", XE{"tns:transferToUserId", h.TransferToUserId}}
}
