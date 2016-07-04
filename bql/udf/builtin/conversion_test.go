package builtin

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
)

func TestBlobToString(t *testing.T) {
	ctx := core.NewContext(nil)
	Convey("Given blob_to_string UDF", t, func() {
		Convey("When passing valid values", func() {
			Convey("Then it should accept an empty blob", func() {
				s, err := blobToRawString(data.Blob(""))
				So(err, ShouldBeNil)
				So(string(s), ShouldBeBlank)
			})

			Convey("Then it should accept JSON", func() {
				a := `{
					"a": 1,
					"b": "2",
					"c": [3.4, "5", {"6": 7}]
				}`
				s, err := blobToRawString(data.Blob(a))
				So(err, ShouldBeNil)
				So(string(s), ShouldEqual, a)
			})

			Convey("Then it should accept Japanese", func() {
				a := "日本語でおｋ"
				s, err := blobToRawString(data.Blob(a))
				So(err, ShouldBeNil)
				So(string(s), ShouldEqual, a)
			})
		})

		Convey("When passing invalid values", func() {
			Convey("Then it should reject binary values", func() {
				_, err := blobToRawString(data.Blob{1, 0xff, 3})
				So(err, ShouldNotBeNil)
			})

			Convey("Then it should reject a broken string", func() {
				b := data.Blob("日本語でおｋ")
				b[5] = 0
				_, err := blobToRawString(b)
				So(err, ShouldNotBeNil)
			})

			Convey("Then it should reject invalid types", func() {
				f := udf.MustConvertGeneric(blobToRawString)
				_, err := f.Call(ctx, data.Map{})
				So(err, ShouldNotBeNil)
			})
		})
	})
}
