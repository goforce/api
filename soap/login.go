package soap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
)

type Connection struct {
	login      *LoginResponse
	cache      map[string]interface{}
	UseStrings bool
}

type LoginResponse struct {
	MetadataServerUrl string `xml:"metadataServerUrl"`
	PasswordExpired   bool   `xml:"passwordExpired"`
	Sandbox           bool   `xml:"sandbox"`
	ServerUrl         string `xml:"serverUrl"`
	SessionId         string `xml:"sessionId"`
	UserId            string `xml:"userId"`
	UserInfo          struct {
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
	} `xml:"userInfo"`
}

type KeyValue struct {
	Key   string
	Value string
}

func Login(host string, username string, password string) (*Connection, error) {
	req := []KeyValue{
		KeyValue{"tns:login/tns:username", username},
		KeyValue{"tns:login/tns:password", password},
	}
	response, err := Call(host+SOAP_API_SERVICES, req)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return nil, err
	}
	var result struct {
		LoginResponse *LoginResponse `xml:"Body>loginResponse>result"`
	}
	err = xml.NewDecoder(response).Decode(&result)
	if err != nil {
		b, _ := ioutil.ReadAll(response)
		return nil, errors.New(fmt.Sprint("error decoding login response:", err, "\n", string(b)))
	}
	return &Connection{login: result.LoginResponse, cache: make(map[string]interface{})}, nil
}

func (co *Connection) Logout() error {
	req := []KeyValue{
		KeyValue{"tns:logout", ""},
	}
	response, err := co.Call(req)
	if response != nil {
		defer response.Close()
	}
	if err != nil {
		return err
	}
	return nil
}
