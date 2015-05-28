package data

import (
	"github.com/goforce/api/soap"
)

type UserInfo struct {
	userInfo *soap.GetUserInfoResult
}

func (ui UserInfo) AccessibilityMode() bool            { return ui.userInfo.AccessibilityMode }
func (ui UserInfo) CurrencySymbol() string             { return ui.userInfo.CurrencySymbol }
func (ui UserInfo) OrgAttachmentFileSizeLimit() int    { return ui.userInfo.OrgAttachmentFileSizeLimit }
func (ui UserInfo) OrgDefaultCurrencyIsoCode() string  { return ui.userInfo.OrgDefaultCurrencyIsoCode }
func (ui UserInfo) OrgDisallowHtmlAttachments() bool   { return ui.userInfo.OrgDisallowHtmlAttachments }
func (ui UserInfo) OrgHasPersonAccounts() bool         { return ui.userInfo.OrgHasPersonAccounts }
func (ui UserInfo) OrganizationId() string             { return ui.userInfo.OrganizationId }
func (ui UserInfo) OrganizationMultiCurrency() bool    { return ui.userInfo.OrganizationMultiCurrency }
func (ui UserInfo) OrganizationName() string           { return ui.userInfo.OrganizationName }
func (ui UserInfo) ProfileId() string                  { return ui.userInfo.ProfileId }
func (ui UserInfo) RoleId() string                     { return ui.userInfo.RoleId }
func (ui UserInfo) SessionSecondsValid() int           { return ui.userInfo.SessionSecondsValid }
func (ui UserInfo) UserDefaultCurrencyIsoCode() string { return ui.userInfo.UserDefaultCurrencyIsoCode }
func (ui UserInfo) UserEmail() string                  { return ui.userInfo.UserEmail }
func (ui UserInfo) UserFullName() string               { return ui.userInfo.UserFullName }
func (ui UserInfo) UserId() string                     { return ui.userInfo.UserId }
func (ui UserInfo) UserLanguage() string               { return ui.userInfo.UserLanguage }
func (ui UserInfo) UserLocale() string                 { return ui.userInfo.UserLocale }
func (ui UserInfo) UserName() string                   { return ui.userInfo.UserName }
func (ui UserInfo) UserTimeZone() string               { return ui.userInfo.UserTimeZone }
func (ui UserInfo) UserType() string                   { return ui.userInfo.UserType }
func (ui UserInfo) UserUiSkin() string                 { return ui.userInfo.UserUiSkin }
