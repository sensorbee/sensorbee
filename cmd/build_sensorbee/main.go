package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func main() {
	app := cli.NewApp()
	app.Name = "build_sensorbee"
	app.Usage = "Build an custom sensorbee command"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "build.yaml",
			Usage: "path to a config file",
		},
		cli.StringFlag{
			Name:  "output-dir",
			Value: "sensorbee",
			Usage: "the output directory in which the generated source code is written",
		},
		cli.StringFlag{
			Name:  "output-filename",
			Value: "customized_main.go",
			Usage: "the name of the filename containing func main()",
		},

		// TODO: an bool option to run go build (this should be true by default, so maybe --no-build should be provided)
		// TODO: add option to go get plugins (maybe fetch-plugins or something like that?)
	}
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) {
	func() {
		if e := recover(); e != nil {
			fmt.Fprintln(os.Stderr, e)
			os.Exit(1)
		}
	}()
	if fn := c.String("output-filename"); fn != filepath.Base(fn) {
		panic(fmt.Errorf("the output file name must only contain a filename: %v", fn))
	}

	config := loadConfig(c.String("config"))

	// TODO: go get all plugins if option is specified

	create(c, config)
}

type Config struct {
	PluginPaths []string `yaml:"plugins"`
}

func loadConfig(path string) *Config {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("cannot load the config file '%v': %v\n", path, err))
	}

	config := &Config{}
	if err := yaml.Unmarshal(b, config); err != nil {
		panic(fmt.Errorf("cannot parse the config file '%v': %v\n", path, err))
	}

	// TODO: validation
	return config
}

func create(c *cli.Context, config *Config) {
	tpl := template.Must(template.New("tpl").Parse(mainGoTemplate))
	var b bytes.Buffer
	if err := tpl.Execute(&b, config); err != nil {
		panic(fmt.Errorf("cannot generate a template source code: %v\n", err))
	}

	// file output
	outputDir := c.String("output-dir")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(fmt.Errorf("cannot create a directory '%v': %v", outputDir, err))
	}
	outFilePath := filepath.Join(outputDir, c.String("output-filename"))
	if err := ioutil.WriteFile(outFilePath, b.Bytes(), 0644); err != nil {
		panic(fmt.Errorf("cannot generate an output file '%v': %v", outFilePath, err))
	}

	// go fmt
	cmd := exec.Command("go", "fmt", outFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(fmt.Errorf("cannot apply go fmt to the generated file: %v", err))
	}

	// TODO: build the command if the option is given
}

const (
	mainGoTemplate = `package main

import (
    "github.com/codegangsta/cli"
    "os"
    "pfi/sensorbee/sensorbee/client"
    "pfi/sensorbee/sensorbee/server"
    "time"
{{range $_, $path := .PluginPaths}}    _ "{{$path}}"
{{end}})

type commandGenerator func() cli.Command

func init() {
    // TODO
    time.Local = time.UTC
}

func main() {
    app := setUpApp([]commandGenerator{
        server.SetUpRunCommand,
        client.SetUpCMDLineToolCommand,
    })

    if err := app.Run(os.Args); err != nil {
        os.Exit(1)
    }
}

func setUpApp(cmds []commandGenerator) *cli.App {
    app := cli.NewApp()
    app.Name = "sensorbee"
    app.Usage = "SenserBee"
    app.Version = "0.0.1" // TODO get dynamic, will be get from external file
    app.Flags = []cli.Flag{
        cli.StringFlag{ // TODO get configuration from external file
            Name:   "config, c",
            Value:  "/etc/sersorbee/sensorbee.config",
            Usage:  "path to the config file",
            EnvVar: "SENSORBEE_CONFIG",
        },
    }
    app.Before = appBeforeHook

    for _, c := range cmds {
        app.Commands = append(app.Commands, c())
    }
    return app
}

func appBeforeHook(c *cli.Context) error {
    if err := loadConfig(c); err != nil {
        return err
    }
    return nil
}

func loadConfig(c *cli.Context) error {
    // TODO load configuration file (YAML)
    return nil
}
`
)
