package data

import (
	"github.com/goforce/api/soap"
	"strconv"
)

type SoapHeader interface {
	xe() *soap.XE
}

type CallOptions struct {
	Client           string
	DefaultNamespace string
}

func (h *CallOptions) xe() *soap.XE {
	if len(h.Client) == 0 && len(h.DefaultNamespace) == 0 {
		return nil
	}
	x := make([]soap.XE, 0, 2)
	if len(h.Client) > 0 {
		x = append(x, soap.XE{"tns:client", h.Client})
	}
	if len(h.DefaultNamespace) > 0 {
		x = append(x, soap.XE{"tns:defaultNamespace", h.DefaultNamespace})
	}
	return &soap.XE{"tns:CallOptions", x}
}

type AllOrNoneHeader struct {
	AllOrNone bool
}

func (h *AllOrNoneHeader) xe() *soap.XE {
	return &soap.XE{"tns:AllOrNoneHeader", soap.XE{"tns:allOrNone", strconv.FormatBool(h.AllOrNone)}}
}

type AllowFieldTruncationHeader struct {
	AllowFieldTruncation bool
}

func (h *AllowFieldTruncationHeader) xe() *soap.XE {
	return &soap.XE{"tns:AllowFieldTruncationHeader", soap.XE{"tns:allowFieldTruncation", strconv.FormatBool(h.AllowFieldTruncation)}}
}

type AssignmentRuleHeader struct {
	AssignmentRuleId string
	UseDefaultRule   *bool
}

func (h *AssignmentRuleHeader) xe() *soap.XE {
	if len(h.AssignmentRuleId) == 0 && h.UseDefaultRule == nil {
		return nil
	}
	x := make([]soap.XE, 0, 2)
	if len(h.AssignmentRuleId) > 0 {
		x = append(x, soap.XE{"tns:assignmentRuleId", h.AssignmentRuleId})
	}
	if h.UseDefaultRule != nil {
		x = append(x, soap.XE{"tns:useDefaultRule", strconv.FormatBool(*h.UseDefaultRule)})
	}
	return &soap.XE{"tns:AssignmentRuleHeader", x}
}

type DisableFeedTrackingHeader struct {
	DisableFeedTracking bool
}

func (h *DisableFeedTrackingHeader) xe() *soap.XE {
	return &soap.XE{"tns:DisableFeedTrackingHeader", soap.XE{"tns:disableFeedTracking", strconv.FormatBool(h.DisableFeedTracking)}}
}

type EmailHeader struct {
	TriggerAutoResponseEmail bool
	TriggerOtherEmail        bool
	TriggerUserEmail         bool
}

func (h *EmailHeader) xe() *soap.XE {
	x := make([]soap.XE, 0, 3)
	x = append(x, soap.XE{"tns:triggerAutoResponseEmail", strconv.FormatBool(h.TriggerAutoResponseEmail)})
	x = append(x, soap.XE{"tns:triggerOtherEmail", strconv.FormatBool(h.TriggerOtherEmail)})
	x = append(x, soap.XE{"tns:triggerUserEmail", strconv.FormatBool(h.TriggerUserEmail)})
	return &soap.XE{"tns:EmailHeader", x}
}

type LocaleOptions struct {
	Language string
}

func (h *LocaleOptions) xe() *soap.XE {
	return &soap.XE{"tns:LocaleOptions", soap.XE{"tns:language", h.Language}}
}

type MruHeader struct {
	UpdateMru bool
}

func (h *MruHeader) xe() *soap.XE {
	return &soap.XE{"tns:MruHeader", soap.XE{"tns:updateMru", strconv.FormatBool(h.UpdateMru)}}
}

type OwnerChangeOptions struct {
	TransferAttachments    bool
	TransferOpenActivities bool
}

func (h *OwnerChangeOptions) xe() *soap.XE {
	x := make([]soap.XE, 0, 2)
	x = append(x, soap.XE{"tns:transferAttachments", strconv.FormatBool(h.TransferAttachments)})
	x = append(x, soap.XE{"tns:transferOpenActivities", strconv.FormatBool(h.TransferOpenActivities)})
	return &soap.XE{"tns:OwnerChangeOptions", x}
}

type QueryOptions struct {
	BatchSize int
}

func (h *QueryOptions) xe() *soap.XE {
	return &soap.XE{"tns:QueryOptions", soap.XE{"tns:batchSize", strconv.Itoa(h.BatchSize)}}
}

type UserTerritoryDeleteHeader struct {
	TransferToUserId string
}

func (h *UserTerritoryDeleteHeader) xe() *soap.XE {
	return &soap.XE{"tns:UserTerritoryDeleteHeader", soap.XE{"tns:transferToUserId", h.TransferToUserId}}
}
