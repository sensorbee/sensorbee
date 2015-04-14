package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
	"testing"
	"time"
)

func TestDefaultTopologyBuilderInterface(t *testing.T) {
	Convey("When creating a default topology builder", t, func() {
		var tb StaticTopologyBuilder = NewDefaultStaticTopologyBuilder()
		So(tb, ShouldNotBeNil)
	})

	Convey("Given a default topology builder", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DefaultSource{}
		b := &DefaultBox{}
		si := &DefaultSink{}
		var err DeclarerError

		Convey("when using a source name twice", func() {
			err = tb.AddSource("mySource", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("mySource", s)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a source name with a box name", func() {
			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a source name with a sink name", func() {
			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name twice", func() {
			err = tb.AddBox("myBox", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("myBox", b)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name with a source name", func() {
			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a box name with a sink name", func() {
			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name twice", func() {
			err = tb.AddSink("mySink", si)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("mySink", si)
			So(err, ShouldNotBeNil)
			Convey("the second time should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name with a source name", func() {
			err = tb.AddSource("someName", s)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when using a sink name with a box name", func() {
			err = tb.AddBox("someName", b)
			So(err, ShouldNotBeNil)
			So(err.Err(), ShouldBeNil)

			err = tb.AddSink("someName", si)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})
	})

	Convey("Given a default topology builder with a source", t, func() {
		tb := NewDefaultStaticTopologyBuilder()
		s := &DefaultSource{}
		tb.AddSource("aSource", s)
		b := &DefaultBox{}
		tb.AddBox("aBox", b)
		si := &DefaultSink{}
		var err DeclarerError

		Convey("when a new box references a non-existing item", func() {
			err = tb.AddBox("otherBox", b).
				Input("something", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing source", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing box", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references multiple items", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new box references an existing source twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aSource", nil).
				Input("aSource", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new box references an existing box twice", func() {
			err = tb.AddBox("otherBox", b).
				Input("aBox", nil).
				Input("aBox", nil)
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references a non-existing item", func() {
			err = tb.AddSink("aSink", si).
				Input("something")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references an existing source", func() {
			err = tb.AddSink("aSink", si).
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references an existing box", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references multiple items", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should work", func() {
				So(err.Err(), ShouldBeNil)
			})
		})

		Convey("when a new sink references an existing source twice", func() {
			err = tb.AddSink("aSink", si).
				Input("aSource").
				Input("aSource")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})

		Convey("when a new sink references an existing box twice", func() {
			err = tb.AddSink("aSink", si).
				Input("aBox").
				Input("aBox")
			So(err, ShouldNotBeNil)
			Convey("adding should fail", func() {
				So(err.Err(), ShouldNotBeNil)
			})
		})
	})
}

func TestBasicDefaultTopologyTransport(t *testing.T) {
	Convey("Given basic topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		b1 := BoxFunc(dummyToUpperBoxFunc)
		tb.AddBox("aBox", &b1).Input("source1", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run()
			So(si.results[0], ShouldEqual, "VALUE")
		})
	})

	Convey("Given 2 sources topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		s2 := &DummyDefaultSource{"hoge"}
		tb.AddSource("source2", s2)
		b1 := BoxFunc(dummyToUpperBoxFunc)
		tb.AddBox("aBox", &b1).
			Input("source1", nil).
			Input("source2", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			start := time.Now()
			t.Run()
			So(len(si.results), ShouldEqual, 2)
			So(si.results, ShouldContain, "VALUE")
			So(si.results, ShouldContain, "HOGE")
			So(start, ShouldHappenWithin, 600*time.Millisecond, time.Now())
		})
	})

	Convey("Given 2 tuples in 1source topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s := &DummyDefaultSource2{"value", "hoge"}
		tb.AddSource("source", s)
		b1 := BoxFunc(dummyToUpperBoxFunc)
		tb.AddBox("aBox", &b1).
			Input("source", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run()
			So(si.results, ShouldResemble, []string{"VALUE", "HOGE"})
		})
	})

	Convey("Given 2 boxes topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		b1 := BoxFunc(dummyToUpperBoxFunc)
		tb.AddBox("aBox", &b1).Input("source1", nil)
		b2 := BoxFunc(dummyAddSuffixBoxFunc)
		tb.AddBox("bBox", &b2).Input("source1", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox").Input("bBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run()
			So(si.results[0], ShouldEqual, "VALUE")
			So(si.results2[0], ShouldEqual, "value_1")
		})
	})

	Convey("Given 2 sinks topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("source1", s1)
		b1 := BoxFunc(dummyToUpperBoxFunc)
		tb.AddBox("aBox", &b1).Input("source1", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		si2 := &DummyDefaultSink{}
		tb.AddSink("si2", si2).Input("aBox")
		t := tb.Build()
		Convey("Run topology with ToUpperBox", func() {
			t.Run()
			So(si.results[0], ShouldEqual, "VALUE")
			So(si2.results[0], ShouldEqual, "VALUE")
		})
	})

}

type DummyDefaultSource struct{ initial string }

func (s *DummyDefaultSource) GenerateStream(w Writer) error {
	time.Sleep(0.5 * 1e9) // to confirm .Run() goroutine
	t := &tuple.Tuple{}
	t.Data = tuple.Map{
		"source": tuple.String(s.initial),
	}
	w.Write(t)
	return nil
}
func (s *DummyDefaultSource) Schema() *Schema {
	var sc Schema = Schema("test")
	return &sc
}

type DummyDefaultSource2 struct {
	initial  string
	initial2 string
}

func (s *DummyDefaultSource2) GenerateStream(w Writer) error {
	t := &tuple.Tuple{}
	t.Data = tuple.Map{
		"source": tuple.String(s.initial),
	}
	w.Write(t)
	t2 := &tuple.Tuple{}
	t2.Data = tuple.Map{
		"source": tuple.String(s.initial2),
	}
	w.Write(t2)
	return nil
}
func (s *DummyDefaultSource2) Schema() *Schema {
	var sc Schema = Schema("test")
	return &sc
}

func dummyToUpperBoxFunc(t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.String()
	t.Data["to-upper"] = tuple.String(strings.ToUpper(string(s)))
	w.Write(t)
	return nil
}

func dummyAddSuffixBoxFunc(t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.String()
	t.Data["add-suffix"] = tuple.String(s + "_1")
	w.Write(t)
	return nil
}

type DummyDefaultSink struct {
	results  []string
	results2 []string
}

func (s *DummyDefaultSink) Write(t *tuple.Tuple) error {
	x, err := t.Data.Get("to-upper")
	if err != nil {
		return nil
	}
	str, _ := x.String()
	s.results = append(s.results, string(str))

	x, err = t.Data.Get("add-suffix")
	if err != nil {
		return nil
	}
	str, _ = x.String()
	s.results2 = append(s.results2, string(str))
	return nil
}
