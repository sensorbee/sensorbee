package client

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"gopkg.in/sensorbee/sensorbee.v0/server/testutil"
	"time"
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

type dummySource struct {
}

func (d *dummySource) GenerateStream(ctx *core.Context, w core.Writer) error {
	// This is a dummy and implementation is very naive. It doesn't stop until
	// it generates all tuples.
	for i := 0; i < 4; i++ {
		now := time.Now()
		if err := w.Write(ctx, &core.Tuple{
			Data: data.Map{
				"int": data.Int(i),
			},
			Timestamp:     now,
			ProcTimestamp: now,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (d *dummySource) Stop(ctx *core.Context) error {
	return nil
}

func createDummySource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Source, error) {
	return &dummySource{}, nil
}

func createRewindableDummySource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (core.Source, error) {
	return core.NewRewindableSource(&dummySource{}), nil
}

func init() {
	bql.MustRegisterGlobalSourceCreator("dummy", bql.SourceCreatorFunc(createDummySource))
	bql.MustRegisterGlobalSourceCreator("rewindable_dummy", bql.SourceCreatorFunc(createRewindableDummySource))
}
