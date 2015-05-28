package soap

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/goforce/log"
	"io/ioutil"
	"net/http"
)

const (
	API_VERSION       string = "33.0"
	SOAP_API_SERVICES        = "/services/Soap/u/" + API_VERSION
	NS_XSI                   = "http://www.w3.org/2001/XMLSchema-instance"
	NS_PARTNER               = "urn:partner.soap.sforce.com"
	NS_OBJECT                = "urn:sobject.partner.soap.sforce.com"
)

func Call(co Connection, reqp *XE, result interface{}, headers ...*XE) error {
	return Post(co.GetServerUrl(),
		reqp.write(&bytes.Buffer{}).Bytes(),
		result, append(headers, &XE{"tns:SessionHeader", XE{"tns:sessionId", co.GetToken()}})...,
	)
}

func Post(url string, soapBody []byte, result interface{}, headers ...*XE) error {
	body := bytes.NewBuffer(make([]byte, 0, len(soapBody)+500))
	body.WriteString(xml.Header)
	body.WriteString(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" ` +
		`xmlns:tns="` + NS_PARTNER + `" ` +
		`xmlns:ens="` + NS_OBJECT + `" ` +
		`><soapenv:Header>`)
	for _, h := range headers {
		h.write(body)
	}
	body.WriteString(`</soapenv:Header><soapenv:Body>`)
	body.Write(soapBody)
	body.WriteString(`</soapenv:Body></soapenv:Envelope>`)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header = http.Header{
		"Content-Type": {"text/xml; charset=utf-8"},
		"SOAPAction":   {"\"\""},
	}
	log.Println(REQUEST, req)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Println(RESPONSE, string(responseBody))
	if response.StatusCode == 200 {
		err := xml.NewDecoder(bytes.NewReader(responseBody)).Decode(result)
		if err != nil {
			return errors.New(fmt.Sprint("error decoding response:", err, "\n", string(responseBody)))
		}
		return nil
	}
	if response.StatusCode == 500 {
		var soapError struct {
			FaultCode   string `xml:"Body>Fault>faultcode"`
			FaultString string `xml:"Body>Fault>faultstring"`
		}
		err := xml.NewDecoder(bytes.NewReader(responseBody)).Decode(&soapError)
		if err != nil {
			return errors.New(fmt.Sprint("error decoding error response:", err, "\n", string(responseBody)))
		}
		return errors.New(fmt.Sprint(soapError.FaultCode, " / ", soapError.FaultString))
	}
	return errors.New(fmt.Sprint("soap call returned:", response.StatusCode, "/", response.Status, "\n", string(responseBody)))
}
