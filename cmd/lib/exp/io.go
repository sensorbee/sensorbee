package exp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"time"
)

type jsonlRecord struct {
	Timestamp data.Timestamp `json:"t"`
	Data      data.Map       `json:"d"`
}

type jsonlSource struct {
	filename string
}

func newJSONLSource(filename string) (*jsonlSource, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	return &jsonlSource{
		filename: filename,
	}, nil
}

func (s *jsonlSource) GenerateStream(ctx *core.Context, w core.Writer) error {
	f, err := os.Open(s.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			if err == io.EOF {
				return nil
			}
			continue
		}

		var jr jsonlRecord
		if err := json.Unmarshal(line, &jr); err != nil {
			ctx.ErrLog(err).WithField("body", string(line)).Error("Cannot parse a json object")
			continue
		}

		t := &core.Tuple{
			Timestamp:     time.Time(jr.Timestamp),
			ProcTimestamp: time.Now(),
			Data:          jr.Data,
		}
		w.Write(ctx, t)
	}
	return nil
}

func (s *jsonlSource) Stop(ctx *core.Context) error {
	// TODO: support graceful stop on interruption
	return nil
}

type jsonlSink struct {
	f   *os.File
	buf *bufio.Writer
}

func newJSONLSink(filename string) (*jsonlSink, error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return &jsonlSink{
		f:   f,
		buf: bufio.NewWriter(f),
	}, nil
}

func (s *jsonlSink) Write(ctx *core.Context, t *core.Tuple) error {
	data, err := json.Marshal(&jsonlRecord{
		Timestamp: data.Timestamp(t.Timestamp),
		Data:      t.Data,
	})
	if err != nil {
		return err
	}
	_, err = s.buf.Write(data)
	if err != nil {
		return err
	}
	_, err = s.buf.Write([]byte("\n"))
	return err
}

func (s *jsonlSink) Close(ctx *core.Context) error {
	if s.f == nil {
		return nil
	}
	s.buf.Flush()
	err := s.f.Close()
	s.f = nil
	return err
}
