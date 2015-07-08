package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Method int

const (
	Get Method = iota
	Post
	Put
	Delete
	OtherMethod
)

func (r Method) String() string {
	switch r {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	case OtherMethod:
		return "OTHER"
	default:
		return "unknown"
	}
}

func Request(method Method, uri string, bodyJSON interface{}) []byte {
	if method == OtherMethod {
		return []byte{}
	}

	req, err := CreateRequest(method, uri, bodyJSON)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	fmt.Println("URI: " + uri)
	return rawRequest(req)
}

func CreateRequest(method Method, uri string, bodyJSON interface{}) (*http.Request, error) {
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

	req, err := http.NewRequest(method.String(), uri, body)
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
