package topology

import (
	"bytes"
	"github.com/codegangsta/cli"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"pfi/sensorbee/sensorbee/server/testutil"
	"testing"
)

// Do NOT run these tests in parallel!

type app struct {
	app *cli.App
	buf *bytes.Buffer
	url string
}

func newApp(url string) *app {
	a := cli.NewApp()
	a.Name = "sensorbee"
	a.Usage = "SenserBee"
	a.Version = "test"
	a.Commands = []cli.Command{
		SetUp(),
	}
	buf := bytes.NewBuffer(nil)
	a.Writer = buf
	return &app{
		app: a,
		buf: buf,
		url: url,
	}
}

func (a *app) run(sub string, args ...string) (string, error) {
	return a.rawRun(sub, append([]string{"--uri", a.url}, args...)...)
}

func (a *app) rawRun(sub string, args ...string) (string, error) {
	if err := a.app.Run(append([]string{"sensorbee", "topology", sub}, args...)); err != nil {
		return "", err
	}

	res, err := ioutil.ReadAll(a.buf)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func TestTopologyCommand(t *testing.T) {
	testMode = true
	testutil.TestAPIWithRealHTTPServer = true
	s := testutil.NewServer()
	defer s.Close()

	Convey("Given a sensorbee topology command", t, func() {
		Convey("When creating a topology", func() {
			out, err := newApp(s.URL()).run("create", "test_topology")
			So(err, ShouldBeNil)
			So(out, ShouldBeBlank)
			So(testExitCode, ShouldEqual, 0)
			Reset(func() {
				newApp(s.URL()).run("drop", "test_topology")
				newApp(s.URL()).run("drop", "test_topology2")
			})

			Convey("Then the topology should be listed", func() {
				out, err := newApp(s.URL()).run("list")
				So(err, ShouldBeNil)
				So(testExitCode, ShouldEqual, 0)

				So(out, ShouldContainSubstring, "test_topology")
			})

			Convey("Then dropping the topology should succeed", func() {
				out, err := newApp(s.URL()).run("drop", "test_topology")
				So(err, ShouldBeNil)
				So(out, ShouldBeBlank)
				So(testExitCode, ShouldEqual, 0)

				Convey("And the topology shouldn't be listed", func() {
					out, err := newApp(s.URL()).run("list")
					So(err, ShouldBeNil)
					So(out, ShouldBeBlank)
					So(testExitCode, ShouldEqual, 0)
				})
			})

			Convey("And creating another topology", func() {
				out, err := newApp(s.URL()).run("create", "another_topology")
				So(err, ShouldBeNil)
				So(testExitCode, ShouldEqual, 0)
				So(out, ShouldBeBlank)

				Convey("Then they should be listed", func() {
					out, err := newApp(s.URL()).run("list")
					So(err, ShouldBeNil)
					So(testExitCode, ShouldEqual, 0)

					So(out, ShouldContainSubstring, "test_topology")
					So(out, ShouldContainSubstring, "another_topology")
				})
			})

			Convey("Then creating another topology having the same name should fail", func() {
				out, err := newApp(s.URL()).run("create", "test_topology")
				So(err, ShouldBeNil)
				So(out, ShouldBeBlank)
				So(testExitCode, ShouldNotEqual, 0)
			})
		})

		Convey("When dropping a nonexistent topology", func() {
			out, err := newApp(s.URL()).run("drop", "test_topology")

			Convey("Then it shouldn't fail", func() {
				So(err, ShouldBeNil)
				So(out, ShouldBeBlank)
				So(testExitCode, ShouldEqual, 0)
			})
		})
	})
}

func TestTopologyCreateCommandValidation(t *testing.T) {
	testMode = true
	testutil.TestAPIWithRealHTTPServer = true
	tmp := testutil.NewServer()
	dummyURL := tmp.URL()
	s := testutil.NewServer()
	defer s.Close()
	url := s.URL()
	tmp.Close() // hope tmp's URL won't be reused too soon.

	Convey("Given a sensorbee topology create command", t, func() {
		// invalid url and api-version won't be tested in other commands
		// because they share the same implementation.

		cases := []struct {
			title string
			args  []string
		}{
			{"When a topology name is missing", []string{"--uri", url}},
			{"When there're too many arguments", []string{"--uri", url, "a", "b", "c"}},
			{"When a url is wrong", []string{"--uri", dummyURL, "test_topology"}},
			{"When a url is invalid", []string{"--uri", ":/hoge/", "test_topology"}},
			{"When a version is invalid", []string{"--uri", url, "--api-version", "v2", "test_topology"}},
			{"When a topology name is invalid", []string{"--uri", url, "test-topology"}},
		}

		for _, c := range cases {
			c := c
			Convey(c.title, func() {
				out, err := newApp(url).rawRun("create", c.args...)
				So(err, ShouldBeNil)
				So(out, ShouldBeBlank)

				Convey("Then the exit code shouldn't be 0", func() {
					So(testExitCode, ShouldNotEqual, 0)
				})
			})
		}
	})
}

func TestTopologyListCommandValidation(t *testing.T) {
	testMode = true
	testutil.TestAPIWithRealHTTPServer = true
	tmp := testutil.NewServer()
	dummyURL := tmp.URL()
	s := testutil.NewServer()
	defer s.Close()
	url := s.URL()
	tmp.Close() // hope tmp's URL won't be reused too soon.

	Convey("Given a sensorbee topology list command", t, func() {
		cases := []struct {
			title string
			args  []string
		}{
			{"When there're too many arguments", []string{"--uri", url, "test_topology"}},
			{"When a url is wrong", []string{"--uri", dummyURL}},
		}

		for _, c := range cases {
			c := c
			Convey(c.title, func() {
				out, err := newApp(url).rawRun("list", c.args...)
				So(err, ShouldBeNil)
				So(out, ShouldBeBlank)

				Convey("Then the exit code shouldn't be 0", func() {
					So(testExitCode, ShouldNotEqual, 0)
				})
			})
		}
	})
}

func TestTopologyDropCommandValidation(t *testing.T) {
	testMode = true
	testutil.TestAPIWithRealHTTPServer = true
	tmp := testutil.NewServer()
	dummyURL := tmp.URL()
	s := testutil.NewServer()
	defer s.Close()
	url := s.URL()
	tmp.Close() // hope tmp's URL won't be reused too soon.

	Convey("Given a sensorbee topology drop command", t, func() {
		cases := []struct {
			title string
			args  []string
		}{
			{"When a topology name is missing", []string{"--uri", url}},
			{"When there're too many arguments", []string{"--uri", url, "a", "b", "c"}},
			{"When a url is wrong", []string{"--uri", dummyURL, "test_topology"}},
			{"When a topology name is invalid", []string{"--uri", url, "test/topology"}},
		}

		for _, c := range cases {
			c := c
			Convey(c.title, func() {
				out, err := newApp(url).rawRun("drop", c.args...)
				So(err, ShouldBeNil)
				So(out, ShouldBeBlank)

				Convey("Then the exit code shouldn't be 0", func() {
					So(testExitCode, ShouldNotEqual, 0)
				})
			})
		}
	})
}
