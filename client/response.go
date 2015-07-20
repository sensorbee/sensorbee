package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"pfi/sensorbee/sensorbee/server/response"
	"strings"
)

// Response wraps a raw HTTP response. Response must only be manipulated from
// one goroutine (i.e. it isn't thread-safe).
type Response struct {
	// Raw has a raw HTTP response.
	Raw *http.Response

	closeStream  chan struct{}
	streamClosed chan struct{}
	streamErr    error

	closed   bool
	closeErr error

	bodyCache []byte
	errCache  *response.Error
	readErr   error

	mimeType     string
	mimeParams   map[string]string
	mimeParseErr error
}

// Close closes the body of the response.
func (r *Response) Close() error {
	if r.closed {
		return r.closeErr
	}
	r.closed = true

	if r.closeStream != nil {
		close(r.closeStream)
	}

	r.closeErr = r.Raw.Body.Close()

	if r.closeStream != nil {
		<-r.streamClosed // wait until stream is closed
	}
	return r.closeErr
}

// IsStream returns true when the response from the server is a stream which
// might have unbounded data.
func (r *Response) IsStream() bool {
	if err := r.parseMediaType(); err != nil {
		return false
	}

	if !strings.Contains(r.mimeType, "multipart") {
		return false
	}
	return r.mimeParams["boundary"] != ""
}

func (r *Response) parseMediaType() error {
	if r.mimeType != "" || r.mimeParseErr != nil {
		return r.mimeParseErr
	}
	t, p, err := mime.ParseMediaType(r.Raw.Header.Get("Content-Type"))
	if err != nil {
		r.mimeParseErr = err
		return err
	}
	r.mimeType = t
	r.mimeParams = p
	return nil
}

// Body returns the response body. Invoking this method will close Raw.Body.
// Don't call this method on a stream response (i.e. a response having
// unlimited length).
func (r *Response) Body() ([]byte, error) {
	if r.IsStream() {
		return nil, errors.New("cannot perform Body method on a stream response")
	}
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

// ReadStreamJSON returns multiple parsed JSONs as interface{}. The caller
// must call Response.Close when it require no more JSONs.
func (r *Response) ReadStreamJSON() (<-chan interface{}, error) {
	if !r.IsStream() {
		return nil, errors.New("the response isn't a stream")
	}

	ch := make(chan interface{})
	r.closeStream = make(chan struct{})
	r.streamClosed = make(chan struct{})

	go func() {
		defer func() {
			close(ch)
			close(r.streamClosed)
		}()

		mr := multipart.NewReader(r.Raw.Body, r.mimeParams["boundary"])
	multipartLoop:
		for {
			err := func() error {
				part, err := mr.NextPart()
				if err != nil {
					return err
				}
				defer part.Close()
				if ct := part.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
					return errors.New("the stream contains non-JSON contents")
				}

				body, err := ioutil.ReadAll(part)
				if err != nil {
					return err
				}

				var js interface{}
				if err := json.Unmarshal(body, &js); err != nil {
					return fmt.Errorf("cannot parse JSON: %v", err)
				}
				ch <- js
				return nil
			}()

			select {
			case <-r.closeStream:
				// When the response is closed first, the error doesn't have
				// to be set.
				break multipartLoop
			default:
			}

			if err != nil {
				if err != io.EOF {
					r.streamErr = err
				}
				return
			}
		}
	}()
	return ch, nil
}
