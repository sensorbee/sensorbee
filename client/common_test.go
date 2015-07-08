package client

import (
	"pfi/sensorbee/sensorbee/server/testutil"
)

func newTestRequester(s *testutil.Server) *Requester {
	r, err := NewRequesterWithClient(s.URL(), "v1", s.HTTPClient())
	if err != nil {
		panic(err)
	}
	return r
}

func do(r *Requester, m Method, path string, body interface{}) (*Response, map[string]interface{}, error) {
	res, err := r.Do(m, path, body)
	if err != nil {
		return nil, nil, err
	}
	var js map[string]interface{}
	if err := res.ReadJSON(&js); err != nil {
		return nil, nil, err
	}
	return res, js, nil
}
