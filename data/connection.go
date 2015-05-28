package data

import (
	"github.com/goforce/api/soap"
)

const (
	PRODUCTION string = "https://login.salesforce.com"
	SANDBOX    string = "https://test.salesforce.com"
)

type Connection struct {
	login        *soap.LoginResponse
	sobjects     map[string]*SObject
	sobjectInfos map[string]*SObjectInfo
}

func (co *Connection) GetServerUrl() string { return co.login.ServerUrl }
func (co *Connection) GetToken() string     { return co.login.SessionId }
func (co *Connection) GetUserId() string    { return co.login.UserId }
func (co *Connection) GetOrgId() string     { return co.login.UserInfo.OrganizationId }

func (co *Connection) UserInfo() UserInfo { return UserInfo{co.login.UserInfo} }

func Login(host string, username string, password string) (*Connection, error) {
	login, err := soap.Login(host, username, password)
	if err != nil {
		return nil, err
	}
	return &Connection{login: login, sobjectInfos: make(map[string]*SObjectInfo)}, nil
}

func (co *Connection) Logout() {
	soap.Logout(co.login)
	co.login = nil
	co.sobjects = nil
	co.sobjectInfos = nil
}
