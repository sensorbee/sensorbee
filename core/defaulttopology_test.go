package core

import (
	. "github.com/smartystreets/goconvey/convey"
	"pfi/sensorbee/sensorbee/core/tuple"
	"strings"
	"testing"
)

func TestDefaultTopology(t *testing.T) {
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

	/* TODO Add tests that show the correct data transport in DefaultTopology.
	 * We should remove the Default* structs from defaulttopology.go and
	 * add some replacements in this file that allow us to check data
	 * flow. For example, a TestSource could just emit two messages with
	 * different IDs and a TestSink could store all received messages
	 * in a local array. Then after Run() is finished, check whether all
	 * items arrived at the sink(s) in the correct order and where processed
	 * correctly by all intermediate boxes.
	 */
	Convey("Given a default topology", t, func() {

		tb := NewDefaultStaticTopologyBuilder()
		s1 := &DummyDefaultSource{"value"}
		tb.AddSource("Source1", s1)
		s2 := &DummyDefaultSource{"test"}
		tb.AddSource("Source2", s2)
		s3 := &DummyDefaultSource{"hoge"}
		tb.AddSource("Source3", s3)
		s4 := &DummyDefaultSource{"fuga"}
		tb.AddSource("Source4", s4)
		s5 := &DummyDefaultSource{"foo"}
		tb.AddSource("Source5", s5)
		b1 := BoxFunc(dummyToUpperBoxFunc)
		tb.AddBox("aBox", &b1).Input("Source1", nil).
			Input("Source2", nil).
			Input("Source3", nil).
			Input("Source4", nil).
			Input("Source5", nil)
		b2 := BoxFunc(dummyFilterBoxFunc)
		tb.AddBox("bBox", &b2).Input("aBox", nil)
		si := &DummyDefaultSink{}
		tb.AddSink("si", si).Input("aBox")
		t := tb.Build()
		Convey("Run topology", func() {
			t.Run()
			So(si.result, ShouldEqual, "HOGE")
		})
	})

}

type DummyDefaultSource struct{ initial string }

func (this *DummyDefaultSource) GenerateStream(w Writer) error {
	t := &tuple.Tuple{}
	t.Data = tuple.Map{
		"source": tuple.String(this.initial),
	}
	w.Write(t)
	return nil
}
func (this *DummyDefaultSource) Schema() *Schema {
	var s Schema = Schema("test")
	return &s
}

func dummyToUpperBoxFunc(t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.String()
	t.Data = tuple.Map{
		"source": tuple.String(strings.ToUpper(string(s))),
	}
	w.Write(t)
	return nil
}

func dummyFilterBoxFunc(t *tuple.Tuple, w Writer) error {
	x, _ := t.Data.Get("source")
	s, _ := x.String()
	if s == "HOGE" {
		t.Data["filtered"] = s
	}
	return nil
}

type DummyDefaultSink struct{ result string }

func (this *DummyDefaultSink) Write(t *tuple.Tuple) error {
	x, err := t.Data.Get("filtered")
	if err != nil {
		return nil
	}
	s, _ := x.String()
	this.result = string(s)
	return nil
}
