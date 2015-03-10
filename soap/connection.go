package soap

import (
	"github.com/goforce/api/soap/core"
)

const (
	PRODUCTION string = "https://login.salesforce.com"
	SANDBOX    string = "https://test.salesforce.com"
)

type Connection struct {
	login    *core.LoginResponse
	sobjects map[string]*SObject
	describe map[string]*DescribeSObject
}

func (co *Connection) Server() string { return co.login.ServerUrl }
func (co *Connection) Token() string  { return co.login.SessionId }
func (co *Connection) User() string   { return co.login.UserId }

func Login(host string, username string, password string) (*Connection, error) {
	login, err := core.Login(host, username, password)
	if err != nil {
		return nil, err
	}
	return &Connection{login: login, describe: make(map[string]*DescribeSObject)}, nil
}

func (co *Connection) Logout() {
	core.Logout(co.login)
	co.login = nil
	co.sobjects = nil
	co.describe = nil
}
