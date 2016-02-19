package shell

import (
	. "github.com/smartystreets/goconvey/convey"
	_ "gopkg.in/sensorbee/sensorbee.v0/client"
	"testing"
)

func TestBQLCommandWithVariousStatements(t *testing.T) {
	Convey("Given a BQL command struct", t, func() {
		cmd := bqlCmd{}
		Convey("When input empty string", func() {
			status, err := cmd.Input("")
			Convey("Then command should not do anything", func() {
				So(err, ShouldBeNil)
				So(status, ShouldEqual, continuousCMD)
			})
		})
		Convey("When input a statement which is end with ';'", func() {
			status, err := cmd.Input("select * from hoge;")
			Convey("Then command should complete prepare and be buffering", func() {
				So(err, ShouldBeNil)
				So(status, ShouldEqual, preparedCMD)
				So(cmd.buffer, ShouldEqual, "select * from hoge;")
				/*
					Convey("And when command eval", func() {
						reqType, uri, body := cmd.Eval()
						Convey("Then command buffer should be flushed", func() {
							So(cmd.buffer, ShouldEqual, "")
							So(reqType, ShouldEqual, client.Post)
							So(uri, ShouldEqual, "/topologies//queries")
							expected := &map[string]interface{}{
								"queries": "select * from hoge;",
							}
							So(body, ShouldResemble, expected)
						})
					})
				*/
			})
		})
		Convey("When input a statement which is not end with ';'", func() {
			status, err := cmd.Input("select")
			Convey("Then command should return continuous state and be buffering", func() {
				So(err, ShouldBeNil)
				So(status, ShouldEqual, continuousCMD)
				So(cmd.buffer, ShouldEqual, "select")
				Convey("And when input next statement which is end with ';'", func() {
					status, err := cmd.Input("* from hoge;")
					Convey("Then command should complete prepare and be buffering", func() {
						So(err, ShouldBeNil)
						So(status, ShouldEqual, preparedCMD)
						So(cmd.buffer, ShouldEqual, "select\n* from hoge;")
						/*
							Convey("And when command eval", func() {
								reqType, url, body := cmd.Eval()
								Convey("Then command buffer should be flushed", func() {
									So(cmd.buffer, ShouldEqual, "")
									So(reqType, ShouldEqual, client.Post)
									So(url, ShouldEqual, "/topologies//queries")
									expected := &map[string]interface{}{
										"queries": "select\n* from hoge;",
									}
									So(body, ShouldResemble, expected)
								})
							})
						*/
					})
				})
				Convey("And when input next statement only ';'", func() {
					status, err := cmd.Input(";")
					Convey("Then command should complete prepare and be buffering", func() {
						So(err, ShouldBeNil)
						So(status, ShouldEqual, preparedCMD)
						So(cmd.buffer, ShouldEqual, "select\n;") // this statement is invalid
						/*
							Convey("And when command eval", func() {
								reqType, url, body := cmd.Eval()
								Convey("Then command buffer should be flushed", func() {
									So(cmd.buffer, ShouldEqual, "")
									So(reqType, ShouldEqual, client.Post)
									So(url, ShouldEqual, "/topologies//queries")
									expected := &map[string]interface{}{
										"queries": "select\n;",
									}
									So(body, ShouldResemble, expected)
								})
							})
						*/
					})
				})
			})
		})
	})
}
