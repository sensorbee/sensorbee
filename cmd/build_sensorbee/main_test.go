package main

import (
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/version"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "build_sensorbee_load_config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	Convey("Given build_sensorbee tool", t, func() {
		Convey("When load config with plugin and commands", func() {
			cfgstr := `plugins:
  - path/to/plugin
commands:
  run:
  runfile:
  repo1:
    path: path/to/repo
  repo2:
    path: path/to/repo2.v1
`
			confName := filepath.Join(dir, "build_test.yaml")
			So(ioutil.WriteFile(confName, []byte(cfgstr), 0644), ShouldBeNil)
			Convey("Then the tool should load the config file", func() {
				conf, err := loadConfig(confName)
				So(err, ShouldBeNil)
				expectedConf := Config{
					PluginPaths: []string{"path/to/plugin"},
					SubCommands: map[string]commandDetail{
						"run":     commandDetail{},
						"runfile": commandDetail{},
						"repo1":   commandDetail{Path: "path/to/repo"},
						"repo2":   commandDetail{Path: "path/to/repo2.v1"},
					},
					Version: version.Version,
				}
				So(*conf, ShouldResemble, expectedConf)
			})
		})

		Convey("When load config with empty yaml", func() {
			confName := filepath.Join(dir, "build_test_empty.yaml")
			So(ioutil.WriteFile(confName, []byte(""), 0644), ShouldBeNil)
			Convey("Then the tool should load the config file", func() {
				conf, err := loadConfig(confName)
				So(err, ShouldBeNil)
				expectedConf := Config{
					PluginPaths: []string(nil),
					SubCommands: map[string]commandDetail{
						"run":      commandDetail{},
						"shell":    commandDetail{},
						"topology": commandDetail{},
						"exp":      commandDetail{},
						"runfile":  commandDetail{},
					},
					Version: version.Version,
				}
				So(*conf, ShouldResemble, expectedConf)
			})
		})
	})
}

func TestCreateMainFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "build_sensorbee_create_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	outFilename := filepath.Join(dir, "test_main.go")

	set := flag.NewFlagSet("dummy", flag.ExitOnError)
	set.String("source-filename", outFilename, "")
	c := cli.NewContext(cli.NewApp(), set, nil)

	Convey("Given build_sensorbee tool", t, func() {
		Convey("When create a main file with a plugin and a buildin command", func() {
			config := &Config{
				PluginPaths: []string{"path/to/plugin"},
				SubCommands: map[string]commandDetail{
					"run": commandDetail{},
				},
				Version: version.Version,
			}
			So(create(c, config), ShouldBeNil)
			Convey("Then the main file should be created", func() {
				b, err := ioutil.ReadFile(outFilename)
				So(err, ShouldBeNil)
				expectedMainFile := `package main

import (
	_ "gopkg.in/sensorbee/sensorbee.v0/bql/udf/builtin"
	"gopkg.in/sensorbee/sensorbee.v0/cmd/lib/run"
	"gopkg.in/sensorbee/sensorbee.v0/version"
	"gopkg.in/urfave/cli.v1"
	"os"
	_ "path/to/plugin"
	"time"
)

func init() {
	// TODO
	time.Local = time.UTC
}

func main() {
	app := cli.NewApp()
	app.Name = "sensorbee"
	app.Usage = "SensorBee built with build_sensorbee ` + version.Version + `"
	app.Version = version.Version
	app.Commands = []cli.Command{
		run.SetUp(),
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
`
				So(string(b), ShouldEqual, expectedMainFile)
			})
		})

		Convey("When create a main file with a plugin and a custom command", func() {
			config := &Config{
				PluginPaths: []string{"path/to/plugin"},
				SubCommands: map[string]commandDetail{
					"repo1": commandDetail{Path: "path/to/repo"},
				},
				Version: version.Version,
			}
			So(create(c, config), ShouldBeNil)
			Convey("Then the main file should be created", func() {
				b, err := ioutil.ReadFile(outFilename)
				So(err, ShouldBeNil)
				expectedMainFile := `package main

import (
	_ "gopkg.in/sensorbee/sensorbee.v0/bql/udf/builtin"
	"gopkg.in/sensorbee/sensorbee.v0/version"
	"gopkg.in/urfave/cli.v1"
	"os"
	_ "path/to/plugin"
	repo1 "path/to/repo"
	"time"
)

func init() {
	// TODO
	time.Local = time.UTC
}

func main() {
	app := cli.NewApp()
	app.Name = "sensorbee"
	app.Usage = "SensorBee built with build_sensorbee ` + version.Version + `"
	app.Version = version.Version
	app.Commands = []cli.Command{
		repo1.SetUp(),
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
`
				So(string(b), ShouldEqual, expectedMainFile)
			})
		})
	})
}
