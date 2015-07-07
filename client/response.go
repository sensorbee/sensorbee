package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"pfi/sensorbee/sensorbee/data"
)

// Response wraps a raw HTTP response.
type Response struct {
	// Raw has a raw HTTP response.
	Raw      *http.Response
	closed   bool
	closeErr error

	errCache *ErrorResponse
	readErr  error
}

// Close closes the body of the response.
func (r *Response) Close() error {
	if r.closed {
		return r.closeErr
	}
	r.closeErr = r.Raw.Body.Close()
	return r.closeErr
}

// IsError returns true when the response is an error.
func (r *Response) IsError() bool {
	return r.Raw.StatusCode < 200 || 300 <= r.Raw.StatusCode
}

// Error reads an error response from the server. It fails if the response
// isn't an error.
func (r *Response) Error() (*ErrorResponse, error) {
	if !r.IsError() {
		return nil, errors.New("the response isn't an error")
	}
	if r.errCache != nil || r.readErr != nil {
		return r.errCache, r.readErr
	}

	res := struct {
		Error *ErrorResponse `json:"error"`
	}{}
	if err := r.ReadJSON(&res); err != nil {
		return nil, err
	}
	r.errCache = res.Error
	return r.errCache, nil
}

// ReadJSON reads the response body as a json and unmarshals it to the given
// data structure. This method closed the response body.
func (r *Response) ReadJSON(js interface{}) error {
	if r.readErr != nil {
		return r.readErr
	}

	defer r.Close()
	body, err := ioutil.ReadAll(r.Raw.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, js); err != nil {
		r.readErr = err
		return err
	}

	r.readErr = errors.New("the response body is already read")
	return nil
}

// TODO: Add method to parse result of SELECT stmt (multipart json)

// ErrorResponse has error information of a request.
type ErrorResponse struct {
	Code      string   `json:"code"`
	Message   string   `json:"message"`
	RequestID string   `json:"request_id"`
	Meta      data.Map `json:"meta"`
}
