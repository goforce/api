package soap

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var _ = fmt.Print

const (
	API_VERSION       string = "31.0"
	PRODUCTION        string = "https://login.salesforce.com"
	SANDBOX           string = "https://test.salesforce.com"
	SOAP_API_SERVICES        = "/services/Soap/u/" + API_VERSION
	NS_XSI                   = "http://www.w3.org/2001/XMLSchema-instance"
)

func (co *Connection) Call(parameters []KeyValue, headers ...SoapHeader) (io.ReadCloser, error) {
	return Call(co.login.ServerUrl, parameters, append(headers, &SessionHeader{co.login.SessionId})...)
}

func (co *Connection) Post(soapBody []byte, headers ...SoapHeader) (io.ReadCloser, error) {
	return Post(co.login.ServerUrl, soapBody, append(headers, &SessionHeader{co.login.SessionId})...)
}

func Call(url string, parameters []KeyValue, headers ...SoapHeader) (io.ReadCloser, error) {
	return Post(url, []byte(marshallParameters(parameters)), headers...)
}

func Post(url string, soapBody []byte, headers ...SoapHeader) (io.ReadCloser, error) {
	body := bytes.NewBuffer(make([]byte, 0, len(soapBody)+500))
	body.WriteString(xml.Header)
	body.WriteString(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" ` +
		`xmlns:tns="urn:partner.soap.sforce.com" ` +
		`xmlns:ens="urn:sobject.partner.soap.sforce.com" ` +
		`><soapenv:Header>`)
	for _, h := range headers {
		body.WriteString(h.String())
	}
	body.WriteString(`</soapenv:Header><soapenv:Body>`)
	body.Write(soapBody)
	body.WriteString(`</soapenv:Body></soapenv:Envelope>`)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Content-Type": {"text/xml; charset=utf-8"},
		"SOAPAction":   {"\"\""},
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 200 {
		return response.Body, nil
	}
	if response.StatusCode == 500 {
		defer response.Body.Close()
		var soapError struct {
			FaultCode   string `xml:"Body>Fault>faultcode"`
			FaultString string `xml:"Body>Fault>faultstring"`
		}
		err := xml.NewDecoder(response.Body).Decode(&soapError)
		if err != nil {
			b, _ := ioutil.ReadAll(response.Body)
			return nil, errors.New(fmt.Sprint("error decoding error response:", err, "\n", string(b)))
		}
		return nil, errors.New(fmt.Sprint(soapError.FaultCode, " / ", soapError.FaultString))
	}
	b, _ := ioutil.ReadAll(response.Body)
	return nil, errors.New(fmt.Sprint("soap call returned:", response.StatusCode, "/", response.Status, "\n", string(b)))
}

func marshallParameters(parameters []KeyValue) string {
	result := ""
	openTags := make([]string, 0, 5)
	open := func(tags []string) string {
		s := ""
		for _, t := range tags {
			s += "<" + t + ">"
		}
		openTags = append(openTags, tags...)
		return s
	}
	close := func(lastToClose int) string {
		if len(openTags) < lastToClose {
			return ""
		}
		s := ""
		for i := len(openTags) - 1; lastToClose <= i; i-- {
			s += "</" + openTags[i] + ">"
		}
		openTags = openTags[:lastToClose]
		return s
	}
	for _, p := range parameters {
		tags := strings.Split(p.Key, "/")
		for i, t := range tags {
			if t == "" {
				continue
			}
			if len(openTags) > i {
				if openTags[i] == t {
					if i+1 == len(tags) {
						// same tags are opened
						// close any open after this one
						result += close(i + 1)
						result += escapeString(p.Value)
						break
					}
				} else {
					// close any open including this
					result += close(i)
					// open all after and set value
					result += open(tags[i:])
					result += escapeString(p.Value)
					break
				}
			} else {
				// open all after and set value
				result += open(tags[i:])
				result += escapeString(p.Value)
				break
			}
		}
	}
	result += close(0)
	return result
}

func escapeString(s string) string {
	w := new(bytes.Buffer)
	xml.EscapeText(w, []byte(s))
	return w.String()
}
