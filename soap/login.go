package soap

import (
	"bytes"
	"github.com/goforce/log"
	"time"
)

func Login(host string, username string, password string) (*LoginResponse, error) {
	start := time.Now()
	defer func() { log.Println(DURATION, "Login took:", time.Since(start)) }()
	req := &XE{"tns:login", []XE{XE{"tns:username", username}, XE{"tns:password", password}}}
	var result struct {
		LoginResponse *LoginResponse `xml:"Body>loginResponse>result"`
	}
	err := Post(host+SOAP_API_SERVICES, req.write(&bytes.Buffer{}).Bytes(), &result)
	if err != nil {
		return nil, err
	}
	return result.LoginResponse, nil
}

func Logout(lr *LoginResponse) {
	start := time.Now()
	defer func() { log.Println(DURATION, "Logout took:", time.Since(start)) }()
	Call(lr, &XE{"tns:logout", ""}, nil)
}

type Connection interface {
	GetServerUrl() string
	GetToken() string
	GetUserId() string
	GetOrgId() string
}

type LoginResponse struct {
	MetadataServerUrl string             `xml:"metadataServerUrl"`
	PasswordExpired   bool               `xml:"passwordExpired"`
	Sandbox           bool               `xml:"sandbox"`
	ServerUrl         string             `xml:"serverUrl"`
	SessionId         string             `xml:"sessionId"`
	UserId            string             `xml:"userId"`
	UserInfo          *GetUserInfoResult `xml:"userInfo"`
}

func (lr *LoginResponse) GetServerUrl() string { return lr.ServerUrl }
func (lr *LoginResponse) GetToken() string     { return lr.SessionId }
func (lr *LoginResponse) GetUserId() string    { return lr.UserId }
func (lr *LoginResponse) GetOrgId() string     { return lr.UserInfo.OrganizationId }
