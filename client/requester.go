package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Requester sends raw HTTP requests to the server. Requester doesn't have
// a state, so it can be used concurrently.
type Requester struct {
	cli    *http.Client
	url    string
	prefix string
}

// NewRequester creates a new requester
func NewRequester(url, version string) (*Requester, error) {
	return NewRequesterWithClient(url, version, http.DefaultClient)
}

// NewRequesterWithClient creates a new requester with a custom HTTP client.
func NewRequesterWithClient(url, version string, cli *http.Client) (*Requester, error) {
	if err := ValidateURL(url); err != nil {
		return nil, err
	}
	if err := ValidateAPIVersion(version); err != nil {
		return nil, err
	}

	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	return &Requester{
		cli:    http.DefaultClient,
		url:    url,
		prefix: "api/" + version,
	}, nil
}

// Do sends a JSON request to server. The caller has to close the body of
// the response.
func (r *Requester) Do(method Method, path string, body interface{}) (*Response, error) {
	req, err := r.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	return r.DoWithRequest(req)
}

// NewRequest creates a new HTTP request having a JSON content. The caller has
// to close the body of the response.
func (r *Requester) NewRequest(method Method, apiPath string, bodyJSON interface{}) (*http.Request, error) {
	var body io.Reader // this is necessary because body has to have io.Reader "type".
	if bodyJSON == nil {
		body = nil
	} else {
		bd, err := json.Marshal(bodyJSON)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bd)
	}

	req, err := http.NewRequest(method.String(), r.url+path.Join(r.prefix, apiPath), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// DoWithRequest sends a custom HTTP request to server.
func (r *Requester) DoWithRequest(req *http.Request) (*Response, error) {
	res, err := r.cli.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{
		Raw: res,
	}, nil
}

// ValidateURL validates if the given URL is valid for the SensorBee API server.
func ValidateURL(u string) error {
	_, err := url.Parse(u)
	if err != nil {
		return err
	}

	// TODO: check other parameters

	// uri can contain path like http://host:port/prefix because the SensorBee
	// server might have "prefix".
	return nil
}

// ValidateAPIVersion checks if the given API version is valid.
func ValidateAPIVersion(v string) error {
	if v != "v1" {
		return fmt.Errorf("not supported version: %v", v)
	}
	return nil
}
