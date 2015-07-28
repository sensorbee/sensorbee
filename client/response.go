package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"pfi/sensorbee/sensorbee/server/response"
	"strconv"
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

	// TODO: refactor this dirty long goroutine.
	go func() {
		defer func() {
			close(ch)
			close(r.streamClosed)
		}()

		reader := bufio.NewReader(r.Raw.Body)
		nextLine := func() ([]byte, []byte, error) {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return nil, nil, err
			}
			nl := line
			line = line[:len(line)-1] // remove '\n'
			if l := len(line); l > 0 && line[l-1] == '\r' {
				line = line[:l-1]
				nl = nl[len(nl)-2:]
			} else {
				nl = nl[len(nl)-1:]
			}
			return line, nl, nil
		}
		boundary := append([]byte("--"), []byte(r.mimeParams["boundary"])...)
		isBoundaryLine := func(line []byte) (bool, bool) {
			if !bytes.HasPrefix(line, boundary) {
				return false, false
			}
			line = line[len(boundary):]
			isFinal := false
			if bytes.HasPrefix(line, []byte("--")) {
				isFinal = true
				line = line[2:]
			}
			for len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
				line = line[1:]
			}
			return len(line) == 0, isFinal
		}

		for { // find first boundary
			if line, _, err := nextLine(); err != nil {
				if err == io.EOF {
					r.streamErr = io.ErrUnexpectedEOF
				} else {
					r.streamErr = err
				}
				return

			} else if b, final := isBoundaryLine(line); b {
				if final {
					return
				}
				break
			}
		}

		body := make([]byte, 0, 4096)
		for { // read each part
			body = body[:0]
			header := http.Header{}
			for { // read header
				line, _, err := nextLine()
				if err != nil {
					if err != io.EOF {
						r.streamErr = err
					}
					return
				}
				if len(line) == 0 {
					break
				}
				l := string(line)
				i := strings.Index(l, ": ")
				if i < 0 {
					continue
				}
				header.Set(l[:i], l[i+2:])
			}

			if ct := header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
				r.streamErr = errors.New("the stream contains non-JSON contents")
				return
			}

			var contentLength int
			clStr := header.Get("Content-Length")
			if clStr != "" {
				cl, err := strconv.Atoi(clStr)
				if err != nil {
					r.streamErr = err
					return
				}
				contentLength = cl
			}

			boundaryFound := false
			finalPart := false
			for {
				line, nl, err := nextLine()
				if err != nil {
					if err != io.EOF {
						r.streamErr = err
					}
					return
				}
				if boundaryFound, finalPart = isBoundaryLine(line); boundaryFound {
					break
				}

				if len(body) != 0 {
					body = append(body, nl...)
				}
				body = append(body, line...)
				if contentLength != 0 && len(body) >= contentLength {
					break
				}
			}

			var js interface{}
			if err := json.Unmarshal(body, &js); err != nil {
				r.streamErr = fmt.Errorf("cannot parse JSON: %v", err)
				return
			}
			select {
			case <-r.closeStream:
			case ch <- js:
			}

			if !boundaryFound {
				for {
					line, _, err := nextLine()
					if err != nil {
						if err != io.EOF {
							r.streamErr = err
						}
						return
					}
					if boundaryFound, finalPart = isBoundaryLine(line); boundaryFound {
						break
					}
				}
			}
			if finalPart {
				return
			}
		}
	}()
	return ch, nil
}

// StreamError returns an error which occurred in a goroutine spawned from
// ReadStreamJSON method. Don't call this method before the channel returned
// from ReadStreamJSON is closed.
func (r *Response) StreamError() error {
	return r.streamErr
}
