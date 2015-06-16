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

type requestType int

const (
	getRequst requestType = iota
	postRequest
	putRequest
	otherRequest
)

func (r requestType) String() string {
	switch r {
	case getRequst:
		return "GET"
	case postRequest:
		return "POST"
	case putRequest:
		return "PUT"
	case otherRequest:
		return "OTHER"
	default:
		return "unknown"
	}
}

func request(reqType requestType, uri string, bodyJSON interface{}) {
	if reqType == otherRequest {
		return
	}

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	fmt.Println("URI: " + uri)
	rawRequest(reqType, client, uri, bodyJSON)
}

func rawRequest(reqType requestType, client *http.Client, uri string,
	bodyJSON interface{}) {
	var body io.Reader
	if bodyJSON == nil {
		body = nil
	} else {
		bd, err := json.Marshal(bodyJSON)
		if err != nil {
			fmt.Println(err)
			return
		}
		body = bytes.NewReader(bd)
	}

	req, err := http.NewRequest(reqType.String(), uri, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	execute(resp)
}

func execute(resp *http.Response) {
	// get response body // TODO get JSON
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
