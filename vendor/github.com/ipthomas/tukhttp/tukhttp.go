package tukhttp

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukutil"
)

var DebugMode = true

type CGLRequest struct {
	Request      string
	X_Api_Key    string
	X_Api_Secret string
	StatusCode   int
	Response     []byte
}
type PIXmRequest struct {
	URL        string
	PID_OID    string
	PID        string
	Timeout    int64
	StatusCode int
	Response   []byte
}
type SOAPRequest struct {
	URL        string
	SOAPAction string
	Timeout    int64
	StatusCode int
	Body       []byte
	Response   []byte
}
type AWS_APIRequest struct {
	URL         string
	Version     int
	Act         string
	Resource    string
	ContentType string
	Timeout     int64
	StatusCode  int
	Body        []byte
	Response    []byte
}
type ClientRequest struct {
	HttpRequest  *http.Request
	ServerURL    string `json:"serverurl"`
	Act          string `json:"act"`
	User         string `json:"user"`
	Org          string `json:"org"`
	Orgoid       string `json:"orgoid"`
	Role         string `json:"role"`
	NHS          string `json:"nhs"`
	PID          string `json:"pid"`
	PIDOrg       string `json:"pidorg"`
	PIDOID       string `json:"pidoid"`
	FamilyName   string `json:"familyname"`
	GivenName    string `json:"givenname"`
	DOB          string `json:"dob"`
	Gender       string `json:"gender"`
	ZIP          string `json:"zip"`
	Status       string `json:"status"`
	XDWKey       string `json:"xdwkey"`
	ID           int    `json:"id"`
	Task         string `json:"task"`
	Pathway      string `json:"pathway"`
	Version      int    `json:"version"`
	ReturnFormat string `json:"returnformat"`
}
type TukHTTPInterface interface {
	newRequest() error
}

var debugMode bool = false

func SetDebugMode(debug bool) {
	debugMode = debug
}
func NewRequest(i TukHTTPInterface) error {
	return i.newRequest()
}
func (i *ClientRequest) newRequest() error {
	req := i.HttpRequest
	req.ParseForm()
	i.Act = req.FormValue(tukcnst.ACT)
	i.User = req.FormValue(tukcnst.QUERY_PARAM_USER)
	i.Org = req.FormValue(tukcnst.QUERY_PARAM_ORG)
	i.Orgoid = tukutil.GetCodeSystemVal(req.FormValue(tukcnst.QUERY_PARAM_ORG))
	i.Role = req.FormValue(tukcnst.QUERY_PARAM_ROLE)
	i.NHS = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_NHS)
	i.PID = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_PID)
	i.PIDOrg = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_PID_ORG)
	i.PIDOID = tukutil.GetCodeSystemVal(req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_PID_ORG))
	i.FamilyName = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_FAMILY_NAME)
	i.GivenName = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_GIVEN_NAME)
	i.DOB = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_DOB)
	i.Gender = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_GENDER)
	i.ZIP = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_ZIP)
	i.Status = req.FormValue(tukcnst.TUK_EVENT_QUERY_PARAM_STATUS)
	i.ID = tukutil.GetIntFromString(req.FormValue(tukcnst.QUERY_PARAM_ID))
	i.Task = req.FormValue(tukcnst.QUERY_PARAM_TASK)
	i.Pathway = req.FormValue(tukcnst.QUERY_PARAM_PATHWAY)
	i.Version = tukutil.GetIntFromString(req.FormValue(tukcnst.QUERY_PARAM_VERSION))
	i.XDWKey = req.FormValue("xdwkey")
	i.ReturnFormat = req.Header.Get(tukcnst.ACCEPT)
	if len(i.XDWKey) > 12 {
		i.Pathway, i.NHS = tukutil.SplitXDWKey(i.XDWKey)
	}
	return nil
}
func (i *SOAPRequest) newRequest() error {
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequest(http.MethodPost, i.URL, strings.NewReader(string(i.Body)))
	if err != nil {
		return err
	}
	if i.SOAPAction != "" {
		req.Header.Set(tukcnst.SOAP_ACTION, i.SOAPAction)
	}
	req.Header.Set(tukcnst.CONTENT_TYPE, tukcnst.SOAP_XML)
	req.Header.Set(tukcnst.ACCEPT, tukcnst.ALL)
	req.Header.Set(tukcnst.CONNECTION, tukcnst.KEEP_ALIVE)
	i.logRequest(req.Header)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	i.StatusCode = resp.StatusCode
	i.Response, err = io.ReadAll(resp.Body)
	i.logResponse()
	return err
}
func (i *PIXmRequest) newRequest() error {
	var err error
	var req *http.Request
	if i.Timeout == 0 {
		i.Timeout = 15
	}
	i.URL = i.URL + "?identifier=" + i.PID_OID + "%7C" + i.PID + tukcnst.FORMAT_JSON_PRETTY
	if req, err = http.NewRequest(tukcnst.HTTP_GET, i.URL, nil); err == nil {
		req.Header.Set(tukcnst.CONTENT_TYPE, tukcnst.APPLICATION_JSON)
		req.Header.Set(tukcnst.ACCEPT, tukcnst.ALL)
		req.Header.Set(tukcnst.CONNECTION, tukcnst.KEEP_ALIVE)
		i.logRequest(req.Header)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
		defer cancel()
		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}
		i.StatusCode = resp.StatusCode
		if i.Response, err = io.ReadAll(resp.Body); err != nil {
			log.Println(err.Error())
		}
		defer resp.Body.Close()
		i.logResponse()
		return nil
	}
	return err
}
func (i *CGLRequest) newRequest() error {
	req, _ := http.NewRequest(tukcnst.HTTP_GET, i.Request, nil)
	req.Header.Set(tukcnst.ACCEPT, tukcnst.APPLICATION_JSON)
	req.Header.Set("X-API-KEY", i.X_Api_Key)
	req.Header.Set("X-API-SECRET", i.X_Api_Secret)
	i.logRequest(req.Header)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	i.StatusCode = resp.StatusCode
	i.Response, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	i.logResponse()
	return err
}
func (i *AWS_APIRequest) newRequest() error {
	if i.Timeout == 0 {
		i.Timeout = 5
	}
	var err error
	var req *http.Request
	var resp *http.Response
	client := &http.Client{}
	if req, err = http.NewRequest(http.MethodPost, i.URL+i.Resource, bytes.NewBuffer(i.Body)); err == nil {
		req.Header.Add(tukcnst.CONTENT_TYPE, i.ContentType)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i.Timeout)*time.Second)
		defer cancel()
		i.logRequest(req.Header)
		log.Printf("Sending HTTP Request - %s", req.URL)
		if resp, err = client.Do(req.WithContext(ctx)); err == nil {
			if resp.StatusCode == http.StatusOK {
				i.Response, err = io.ReadAll(resp.Body)
				defer resp.Body.Close()
				i.StatusCode = resp.StatusCode
				i.logResponse()
			}
		}
	}
	if err != nil {
		log.Println(err.Error())
	}
	return err
}
func (i *AWS_APIRequest) logRequest(headers http.Header) {
	if debugMode {
		log.Printf("HTTP POST Request Headers")
		tukutil.Log(headers)
		log.Printf("HTTP Request\nURL = %s\nTimeout = %v\nMessage body\n%s", i.URL, i.Timeout, string(i.Body))
	}
}
func (i *AWS_APIRequest) logResponse() {
	if debugMode {
		log.Printf("HTML Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
	}
}
func (i *SOAPRequest) logRequest(headers http.Header) {
	if debugMode {
		log.Println("SOAP Request Headers")
		tukutil.Log(headers)
		log.Printf("SOAP Request\nURL = %s\nAction = %s\nTimeout = %v\n\n%s", i.URL, i.SOAPAction, i.Timeout, string(i.Body))
	}
}
func (i *SOAPRequest) logResponse() {
	if debugMode {
		log.Printf("SOAP Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
	}
}
func (i *PIXmRequest) logRequest(headers http.Header) {
	log.Println("HTTP GET Request Headers")
	tukutil.Log(headers)
	log.Printf("HTTP Request\nURL = %s\nTimeout = %v", i.URL, i.Timeout)
}
func (i *PIXmRequest) logResponse() {
	log.Printf("HTML Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
}
func (i *CGLRequest) logRequest(headers http.Header) {
	log.Printf("HTTP GET Request Headers")
	tukutil.Log(headers)
	log.Printf("HTTP Request\nURL = %s - Timeout = %v", i.Request, 5)
}
func (i *CGLRequest) logResponse() {
	log.Printf("HTML Response - Status Code = %v\n%s", i.StatusCode, string(i.Response))
}
