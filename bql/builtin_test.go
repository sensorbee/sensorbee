package bql

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

type testFileWriter struct {
	m   sync.Mutex
	c   *sync.Cond
	cnt int
	tss []time.Time
}

func (w *testFileWriter) Write(ctx *core.Context, t *core.Tuple) error {
	w.m.Lock()
	defer w.m.Unlock()
	w.cnt++
	w.tss = append(w.tss, t.Timestamp)
	w.c.Broadcast()
	return nil
}

func (w *testFileWriter) wait(n int) {
	w.m.Lock()
	defer w.m.Unlock()
	for w.cnt < n {
		w.c.Wait()
	}
}

func TestFileSource(t *testing.T) {
	f, err := ioutil.TempFile("", "sbtest_bql_file_source")
	if err != nil {
		t.Fatal("Cannot create a temp file:", err)
	}
	name := f.Name()
	defer func() {
		os.Remove(name)
	}()
	now := time.Now()
	nowTs := data.Timestamp(now)

	// an empty line is intentionally included
	_, err = io.WriteString(f, fmt.Sprintf(`{"int":1, "ts":%v}
 
 {"int":2, "ts":%v}
  {"int":3, "ts":%v} `, nowTs, nowTs, nowTs))
	f.Close()
	if err != nil {
		t.Fatal("Cannot write to the temp file:", err)
	}

	Convey("Given a file", t, func() {
		ctx := core.NewContext(nil)
		params := data.Map{"path": data.String(name)}
		w := &testFileWriter{}
		w.c = sync.NewCond(&w.m)

		Convey("When reading the file by file source with default params", func() {
			s, err := createFileSource(ctx, &IOParams{}, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})

			err = s.GenerateStream(ctx, w)
			So(err, ShouldBeNil)

			Convey("Then it should emit all tuples", func() {
				So(w.cnt, ShouldEqual, 3)
			})
		})

		Convey("When reading the file with custom timestamp field", func() {
			params["timestamp_field"] = data.String("ts")
			s, err := createFileSource(ctx, &IOParams{}, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})

			err = s.GenerateStream(ctx, w)
			So(err, ShouldBeNil)

			Convey("Then it should emit all tuples", func() {
				So(w.cnt, ShouldEqual, 3)
			})

			Convey("Then it should have custom timestamps", func() {
				So(w.tss, ShouldHaveLength, w.cnt)
				for _, ts := range w.tss {
					So(ts, ShouldHappenOnOrBetween, now, now)
				}
			})
		})

		Convey("When reading the file with rewindable", func() {
			params["rewindable"] = data.True
			s, err := createFileSource(ctx, &IOParams{}, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})

			ch := make(chan error, 1)
			go func() {
				ch <- s.GenerateStream(ctx, w)
			}()

			Convey("Then it should write all tuples and pause", func() {
				w.wait(3)
				So(w.cnt, ShouldEqual, 3)
				select {
				case <-ch:
					So("The source should not have stopped yet", ShouldBeNil)
				default:
				}
			})

			Convey("Then it should be able to rewind", func() {
				w.wait(3)
				rs := s.(core.RewindableSource)
				So(rs.Rewind(ctx), ShouldBeNil)
				w.wait(6)
				So(w.cnt, ShouldEqual, 6)
				select {
				case <-ch:
					So("The source should not have stopped yet", ShouldBeNil)
				default:
				}
			})

			Convey("Then it should be able to stop", func() {
				So(s.Stop(ctx), ShouldBeNil)
				err := <-ch
				So(err, ShouldBeNil)
			})
		})

		Convey("When reading the file with a repeat parameter", func() {
			params["repeat"] = data.Int(3)
			s, err := createFileSource(ctx, &IOParams{}, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})

			err = s.GenerateStream(ctx, w)
			So(err, ShouldBeNil)

			Convey("Then it should emit all tuples", func() {
				// The source emits 3 tuples for 4 times including the first run.
				So(w.cnt, ShouldEqual, 12)
			})
		})

		Convey("When reading the file with a negative repeat parameter", func() {
			params["repeat"] = data.Int(-1)
			s, err := createFileSource(ctx, &IOParams{}, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})

			ch := make(chan error, 1)
			go func() {
				ch <- s.GenerateStream(ctx, w)
			}()

			Convey("Then it should infinitely emit tuples", func() {
				w.wait(100)
				So(w.cnt, ShouldBeGreaterThanOrEqualTo, 100)
				select {
				case <-ch:
					So("The source should not have stopped yet", ShouldBeNil)
				default:
				}
			})

			Convey("Then it should be able to stop", func() {
				So(s.Stop(ctx), ShouldBeNil)
				err := <-ch
				So(err, ShouldBeNil)
			})
		})

		Convey("When reading the file with an interval parameter", func() {
			params["interval"] = data.Float(0.0001) // assuming this number is big enough
			s, err := createFileSource(ctx, &IOParams{}, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Stop(ctx)
			})

			err = s.GenerateStream(ctx, w)
			So(err, ShouldBeNil)

			Convey("Then it should emit all tuples", func() {
				So(w.cnt, ShouldEqual, 3)
			})

			Convey("Then tuples' timestamps should have proper intervals", func() {
				So(w.tss, ShouldHaveLength, w.cnt)
				for i := 1; i < len(w.tss); i++ {
					So(w.tss[i], ShouldHappenOnOrAfter, w.tss[i-1].Add(100*time.Microsecond))
				}
			})
		})

		Convey("When creating a file source with invalid parameters", func() {
			Convey("Then missing path parameter should result in an error", func() {
				delete(params, "path")
				_, err := createFileSource(ctx, &IOParams{}, params)
				So(err, ShouldNotBeNil)
			})

			Convey("Then ill-formed timestamp_path should result in an error", func() {
				params["timestamp_field"] = data.String("/this/isnt/a/xpath")
				_, err := createFileSource(ctx, &IOParams{}, params)
				So(err, ShouldNotBeNil)
			})

			Convey("Then invalid timestamp_path value should result in an error", func() {
				params["timestamp_field"] = data.True
				_, err := createFileSource(ctx, &IOParams{}, params)
				So(err, ShouldNotBeNil)
			})

			Convey("Then invalid rewindable value should result in an error", func() {
				params["rewindable"] = data.Int(1)
				_, err := createFileSource(ctx, &IOParams{}, params)
				So(err, ShouldNotBeNil)
			})

			Convey("Then invalid repeat value should result in an error", func() {
				params["repeat"] = data.Float(1.5)
				_, err := createFileSource(ctx, &IOParams{}, params)
				So(err, ShouldNotBeNil)
			})

			Convey("Then invalid interval value should result in an error", func() {
				params["interval"] = data.Map{}
				_, err := createFileSource(ctx, &IOParams{}, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
