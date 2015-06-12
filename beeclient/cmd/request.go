package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type requestType int

const (
	getRequst requestType = iota
	postRequest
	otherRequest
)

func request(reqType requestType, uri string) {
	if reqType == otherRequest {
		return
	}
	values := url.Values{} // value is not need because use gocraft/web package

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	fmt.Println("URI: " + uri)
	switch reqType {
	case getRequst:
		get(client, values, uri)
	case postRequest:
		post(client, values, uri)
	}
}

func get(client *http.Client, values url.Values, uri string) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.URL.RawQuery = values.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	execute(resp)
}

func post(client *http.Client, values url.Values, uri string) {
	req, err := http.NewRequest("POST", uri, strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-rulencoded")

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
