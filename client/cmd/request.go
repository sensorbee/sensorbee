package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestType int

const (
	GetRequest RequestType = iota
	PostRequest
	PutRequest
	DeleteRequest
	otherRequest
)

func (r RequestType) String() string {
	switch r {
	case GetRequest:
		return "GET"
	case PostRequest:
		return "POST"
	case PutRequest:
		return "PUT"
	case DeleteRequest:
		return "DELETE"
	case otherRequest:
		return "OTHER"
	default:
		return "unknown"
	}
}

func Request(reqType RequestType, uri string, bodyJSON interface{}) []byte {
	if reqType == otherRequest {
		return []byte{}
	}

	req, err := CreateRequest(reqType, uri, bodyJSON)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	fmt.Println("URI: " + uri)
	return rawRequest(req)
}

func CreateRequest(reqType RequestType, uri string, bodyJSON interface{}) (*http.Request, error) {
	var body io.Reader
	if bodyJSON == nil {
		body = nil
	} else {
		bd, err := json.Marshal(bodyJSON)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bd)
	}

	req, err := http.NewRequest(reqType.String(), uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func rawRequest(req *http.Request) []byte {
	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	defer resp.Body.Close()

	return execute(resp)
}

func execute(resp *http.Response) []byte {
	// get response body // TODO get JSON
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println(err)
		fmt.Println(string(b))
	}
	fmt.Println(m) // for debug
	return b
}
