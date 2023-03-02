// (c) Copyright 2022 Hewlett Packard Enterprise Development LP

package restclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/hpe-hcss/lh-cdc-singularity/constants"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/hpe-storage/common-host-libs/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	StatusCodeCreated         = 201
	StatusCodeOk              = 200
	StatusCodeTooManyRequests = 429
	StatusCodeNotFound        = 404
)

type ServerInfo struct {
	RestServer string
	UserInfo   UserInfo
}

type UserInfo struct {
	UserName   string
	UserPwd    string
	UserTicket string
}

type Error struct {
	ErrorId          int    `json:"id"`
	ErrorDescription string `json:"desc"`
}

type RestClient struct {
}

func NewRestClient() *RestClient {
	return &RestClient{}
}

func (c *RestClient) constructQuery(restserver string, userinfo UserInfo, path string, rawquery string) url.URL {
	u, err := url.Parse(restserver)
	if err != nil {
		return *u
	}
	// only host name without scheme and port does not parse well
	if u.Host == "" {
		u.Host = u.Path
	}

	u.Path = path
	u.RawQuery = rawquery
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	return *u
}

func (c *RestClient) ExecuteRestRequest(requestMethod string, baseurl string, resource string, rawQuery string, userInfo UserInfo, restServer string, asUser string, timeoutsecs int, body []byte, header map[string]string, token string) (int, []byte, error) {
	log.Infof("Inside ExecuteRestRequest")
	var querystring string
	if resource != "" {
		querystring = baseurl + "/" + resource
	} else {
		querystring = baseurl
	}

	query := c.constructQuery(restServer, userInfo, querystring, rawQuery)
	log.Infof("request method = %v action = %v base url = %v raw query = %v rest server = %v user = %v, time = %v", requestMethod, resource, baseurl, rawQuery, restServer, asUser, timeoutsecs)

	log.Infof("query = %v", query)
	resp, err := c.executeQueryWithTimeout(requestMethod, query, asUser, userInfo, timeoutsecs, body, header, token)
	if err != nil {
		return 0, []byte{}, status.Errorf(codes.Internal, fmt.Sprintf("Error in REST call, error: %v", err))
	} else {
		log.Infof("rest call successful, code = %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, []byte{}, status.Errorf(codes.Internal, fmt.Sprintf("Invalid REST response, error: %v", err))
	}

	return resp.StatusCode, body, nil
}

func (c *RestClient) executeQueryWithTimeout(requestMethod string, query url.URL, asUser string, userInfo UserInfo, timeoutsecs int, body []byte, header map[string]string, token string) (*http.Response, error) {
	w := os.Stdout
	if timeoutsecs < 60 {
		timeoutsecs = 60
	}
	// TODO: Take a param to verify via cert authority
	client := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout:       time.Duration(timeoutsecs) * time.Second,
			ResponseHeaderTimeout: time.Duration(timeoutsecs) * time.Second,
			TLSClientConfig: &tls.Config{
				KeyLogWriter:       w,
				InsecureSkipVerify: true,
			},
		},
	}
	qstring := query.String()
	log.Infof("query string = %v", qstring)
	var reqbody io.Reader
	if body != nil {
		reqbody = bytes.NewBuffer(body)
	} else {
		log.Infof("request body not set")
	}
	req, err := http.NewRequest(requestMethod, qstring, reqbody)
	if err != nil {
		return nil, err
	}
	if asUser != "" {
		req.Header.Set("x-mapr-impersonated-user", asUser)
	}
	if reqbody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if userInfo.UserTicket != "" {
		req.Header.Set("Authorization", "MAPR-Negotiate "+userInfo.UserTicket)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Add("API-Version", constants.GLM_CLIENT_VERSION)

	for key, value := range header {
		req.Header.Set(key, value)
	}
	log.Infof("headers = %v", req.Header)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
