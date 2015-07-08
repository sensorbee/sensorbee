package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"pfi/sensorbee/sensorbee/server/response"
)

// Response wraps a raw HTTP response.
type Response struct {
	// Raw has a raw HTTP response.
	Raw      *http.Response
	closed   bool
	closeErr error

	bodyCache []byte
	errCache  *response.Error
	readErr   error
}

// Close closes the body of the response.
func (r *Response) Close() error {
	if r.closed {
		return r.closeErr
	}
	r.closeErr = r.Raw.Body.Close()
	return r.closeErr
}

// Body returns the response body. Invoking this method will close Raw.Body.
// Don't call this method on a stream response (i.e. a response having
// unlimited length).
func (r *Response) Body() ([]byte, error) {
	if r.bodyCache != nil || r.readErr != nil {
		return r.bodyCache, r.readErr
	}

	defer r.Close()
	r.bodyCache, r.readErr = ioutil.ReadAll(r.Raw.Body)
	return r.bodyCache, r.readErr
}

// IsError returns true when the response is an error.
func (r *Response) IsError() bool {
	return r.Raw.StatusCode < 200 || 300 <= r.Raw.StatusCode
}

// Error reads an error response from the server. It fails if the response
// isn't an error.
func (r *Response) Error() (*response.Error, error) {
	if !r.IsError() {
		return nil, errors.New("the response isn't an error")
	}
	if r.errCache != nil || r.readErr != nil {
		return r.errCache, r.readErr
	}

	res := struct {
		Error *response.Error `json:"error"`
	}{}
	if err := r.ReadJSON(&res); err != nil {
		return nil, err
	}
	if res.Error == nil {
		return nil, errors.New("invalid error response")
	}
	r.errCache = res.Error
	return r.errCache, nil
}

// ReadJSON reads the response body as a json and unmarshals it to the given
// data structure. This method closed the response body. This method can be
// called multiple times with different types of arguments. Don't call this
// method on a stream response (i.e. a response having unlimited length).
func (r *Response) ReadJSON(js interface{}) error {
	body, err := r.Body()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, js); err != nil {
		r.readErr = err
		return err
	}
	return nil
}

// TODO: Add method to parse result of SELECT stmt (multipart json)
