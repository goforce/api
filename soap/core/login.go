package core

import (
	"bytes"
	"github.com/goforce/log"
	"time"
)

type LoginResponse struct {
	MetadataServerUrl string             `xml:"metadataServerUrl"`
	PasswordExpired   bool               `xml:"passwordExpired"`
	Sandbox           bool               `xml:"sandbox"`
	ServerUrl         string             `xml:"serverUrl"`
	SessionId         string             `xml:"sessionId"`
	UserId            string             `xml:"userId"`
	UserInfo          *GetUserInfoResult `xml:"userInfo"`
}

func (lr *LoginResponse) Server() string { return lr.ServerUrl }
func (lr *LoginResponse) Token() string  { return lr.SessionId }
func (lr *LoginResponse) User() string   { return lr.UserId }

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
