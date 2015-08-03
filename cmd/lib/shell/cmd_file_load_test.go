package shell

import (
	. "github.com/smartystreets/goconvey/convey"
	_ "pfi/sensorbee/sensorbee/client"
	"testing"
)

func TestValidFilePathLoading(t *testing.T) {
	Convey("Given a valid BQL file", t, func() {
		filePath := "dummy_file_test.bql"
		targetCmd := fileLoadCmd{}
		Convey("When input the file", func() {
			statusType, err := targetCmd.Input("file " + filePath)
			So(err, ShouldBeNil)
			Convey("Then buffering the queries", func() {
				expected := `CREATE SOURCE dummy TYPE homhom;
SELECT foo FROM dummy;
`
				So(statusType, ShouldEqual, preparedCMD)
				So(targetCmd.queries, ShouldEqual, expected)
				/*
					Convey("And then request the body include queries", func() {
						m := map[string]interface{}{
							"queries": expected,
						}
						method, uri, body := targetCmd.Eval()
						So(method, ShouldEqual, client.Post)
						So(uri, ShouldEqual, "/topologies//queries")
						So(body, ShouldResemble, m)
					})
				*/
			})
		})
	})
}

func TestInvalidFilePathLoading(t *testing.T) {
	Convey("Given a invalid file path", t, func() {
		filePath := "invalid_file.bql"
		targetCmd := fileLoadCmd{}
		Convey("When input the file", func() {
			statusType, err := targetCmd.Input("file " + filePath)
			So(err, ShouldNotBeNil)
			Convey("Then return invalid status", func() {
				So(statusType, ShouldEqual, invalidCMD)
			})
		})
	})
}

func TestEmptyFilePathLoading(t *testing.T) {
	Convey("Given a file load command", t, func() {
		targetCmd := fileLoadCmd{}
		Convey("When input empty path", func() {
			statusType, err := targetCmd.Input("file")
			So(err, ShouldNotBeNil)
			Convey("Then return invalid status", func() {
				So(statusType, ShouldEqual, invalidCMD)
			})
		})
	})
}
