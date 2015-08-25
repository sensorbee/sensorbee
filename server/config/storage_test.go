package config

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStorage(t *testing.T) {
	Convey("Given a JSON config for storage section", t, func() {
		Convey("When the config is valid", func() {
			s, err := NewStorage(toMap(`{"uds":{"type":"in_memory"}}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters", func() {
				So(s.UDS.Type, ShouldEqual, "in_memory")
				So(s.UDS.Params, ShouldBeEmpty)
			})
		})

		Convey("When the config only has required parameters", func() {
			// no required parameter at the moment
			s, err := NewStorage(toMap(`{}`))

			Convey("Then it should have given parameters and default values", func() {
				So(err, ShouldBeNil)
				So(s.UDS.Type, ShouldEqual, "in_memory")
			})
		})

		Convey("When the config has an undefined field", func() {
			_, err := NewStorage(toMap(`{"undefined":"invalid"}`))

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When the uds.type parameter is missing", func() {
			_, err := NewStorage(toMap(`{"uds":{}}`))

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestUDSStorageInMemory(t *testing.T) {
	Convey("Given a JSON config for storage.uds section with type 'in_memory'", t, func() {
		Convey("When the config is valid", func() {
			s, err := NewStorage(toMap(`{"uds":{"type":"in_memory"}}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters", func() {
				So(s.UDS.Type, ShouldEqual, "in_memory")
				So(s.UDS.Params, ShouldBeEmpty)
			})
		})

		Convey("When the config has an undefined field", func() {
			_, err := NewStorage(toMap(`{"uds":{"type":"in_memory","unknown":"invalid"}}`))

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When params is null", func() {
			s, err := NewStorage(toMap(`{"uds":{"type":"in_memory","params":null}}`))

			Convey("Then it should be valid", func() {
				So(err, ShouldBeNil)
				So(s.UDS.Params, ShouldNotBeNil)
			})
		})

		Convey("When params is an empty object rather than null", func() {
			s, err := NewStorage(toMap(`{"uds":{"type":"in_memory","params":{}}}`))

			Convey("Then it should be valid", func() {
				So(err, ShouldBeNil)
				So(s.UDS.Params, ShouldNotBeNil)
			})
		})

		Convey("When params has an additional field", func() {
			_, err := NewStorage(toMap(`{"uds":{"type":"in_memory","params":{"a":"b"}}}`))

			Convey("Then it shouldn't be valid", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestUDSStorageFS(t *testing.T) {
	Convey("Given a JSON config for storage.uds section with type 'fs'", t, func() {
		Convey("When the config is valid", func() {
			s, err := NewStorage(toMap(`{"uds":{"type":"fs"}}`))
			So(err, ShouldBeNil)

			Convey("Then it should have given parameters", func() {
				So(s.UDS.Type, ShouldEqual, "fs")
				So(s.UDS.Params, ShouldBeEmpty)
			})
		})

		Convey("When the config has an undefined field", func() {
			_, err := NewStorage(toMap(`{"uds":{"type":"fs","unknown":"invalid"}}`))

			Convey("Then it should be invalid", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When params is null", func() {
			s, err := NewStorage(toMap(`{"uds":{"type":"fs","params":null}}`))

			Convey("Then it should be valid", func() {
				So(err, ShouldBeNil)
				So(s.UDS.Params, ShouldNotBeNil)
			})
		})

		Convey("When params is an empty object rather than null", func() {
			s, err := NewStorage(toMap(`{"uds":{"type":"fs","params":{"dir":"/path/to/dir"}}}`))

			Convey("Then it should be valid", func() {
				So(err, ShouldBeNil)
				So(s.UDS.Params, ShouldNotBeNil)
			})
		})

		Convey("When params has an additional field", func() {
			_, err := NewStorage(toMap(`{"uds":{"type":"fs","params":{"a":"b"}}}`))

			Convey("Then it shouldn't be valid", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When validating dir", func() {
			for _, dir := range []string{"storage", "/path/to/storage"} {
				Convey(fmt.Sprint("Then it should accept ", dir), func() {
					s, err := NewStorage(toMap(fmt.Sprintf(`{"uds":{"type":"fs","params":{"dir":"%v"}}}`, dir)))
					So(err, ShouldBeNil)
					So(s.UDS.Params["dir"], ShouldEqual, dir)
				})
			}

			// TODO: test invalid format
		})

		Convey("When validating temp_dir", func() {
			for _, dir := range []string{"storage", "/path/to/storage"} {
				Convey(fmt.Sprint("Then it should accept ", dir), func() {
					s, err := NewStorage(toMap(fmt.Sprintf(`{"uds":{"type":"fs","params":{"dir":"/path/to/dir","temp_dir":"%v"}}}`, dir)))
					So(err, ShouldBeNil)
					So(s.UDS.Params["temp_dir"], ShouldEqual, dir)
				})
			}

			// TODO: test invalid format
		})
	})
}
