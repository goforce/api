package core

type GetUserInfoResult struct {
	AccessibilityMode          bool   `xml:"accessibilityMode"`
	CurrencySymbol             string `xml:"currencySymbol"`
	OrgAttachmentFileSizeLimit int    `xml:"orgAttachmentFileSizeLimit"`
	OrgDefaultCurrencyIsoCode  string `xml:"orgDefaultCurrencyIsoCode"`
	OrgDisallowHtmlAttachments bool   `xml:"orgDisallowHtmlAttachments"`
	OrgHasPersonAccounts       bool   `xml:"orgHasPersonAccounts"`
	OrganizationId             string `xml:"organizationId"`
	OrganizationMultiCurrency  bool   `xml:"organizationMultiCurrency"`
	OrganizationName           string `xml:"organizationName"`
	ProfileId                  string `xml:"profileId"`
	RoleId                     string `xml:"roleId"`
	SessionSecondsValid        int    `xml:"sessionSecondsValid"`
	UserDefaultCurrencyIsoCode string `xml:"userDefaultCurrencyIsoCode"`
	UserEmail                  string `xml:"userEmail"`
	UserFullName               string `xml:"userFullName"`
	UserId                     string `xml:"userId"`
	UserLanguage               string `xml:"userLanguage"`
	UserLocale                 string `xml:"userLocale"`
	UserName                   string `xml:"userName"`
	UserTimeZone               string `xml:"userTimeZone"`
	UserType                   string `xml:"userType"`
	UserUiSkin                 string `xml:"userUiSkin"`
}
