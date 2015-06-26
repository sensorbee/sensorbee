package main

import (
	"bytes"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

var (
	pluginYAMLFilePath   string
	outputMainGoDir      string
	outputMainGoFileName string
)

type PluginPath struct {
	Paths []string
}

func init() {
	flags()
}

func flags() {
	// default that plugin.yaml and this main.go are in same directory,
	// or using customize file path with flag
	//  build_sensorbee -in=/hoge/plugin.yaml
	flag.StringVar(&pluginYAMLFilePath, "in", "plugin.yaml", "The file path of plugin.yaml.")
	// default that main go file are in "sensorbee" directory to create "sensorbee"
	// binary, or using customize directory name with flag
	//  build_sensorbee -outdir=sensorbee2
	flag.StringVar(&outputMainGoDir, "outdir", "sensorbee", "The output file directory of customized main.go.")
	// default that main go file name is "customized_main.go",
	// or using customize output file name with flag
	//  build_sensorbee -outname=foo_main.go
	flag.StringVar(&outputMainGoFileName, "outname", "customized_main.go", "The output file name of customized main.go.")
}

func main() {
	flag.Parse()

	b, err := ioutil.ReadFile(pluginYAMLFilePath)
	if err != nil {
		panic(err)
	}

	m := map[string]interface{}{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		panic(err)
	}

	paths, ok := m["Plugins"]
	if !ok {
		fmt.Fprintf(os.Stdout, "not found plug-in list\n")
	}
	pathList, ok := paths.([]interface{})
	if !ok || len(pathList) < 1 {
		fmt.Fprintf(os.Stdout, "plug-in list is not list style or empty\n")
	}
	pluginPaths := []string{}
	for _, path := range pathList {
		pluginPaths = append(pluginPaths, path.(string))
	}
	create(PluginPath{
		Paths: pluginPaths,
	})

	os.Exit(0)
}

func create(pluginPaths PluginPath) {
	tpl := template.Must(template.New("tpl").Parse(mainGoTemplate))
	var b bytes.Buffer
	if err := tpl.Execute(&b, pluginPaths); err != nil {
		panic(err)
	}

	// file output
	if err := os.MkdirAll(outputMainGoDir, os.ModePerm); err != nil {
		panic(err)
	}
	outFilePath := outputMainGoDir + string(os.PathSeparator) + outputMainGoFileName
	if err := ioutil.WriteFile(outFilePath, b.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}

	// go fmt
	cmd := exec.Command("go", "fmt", outFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

const (
	mainGoTemplate = `package main

import (
    "github.com/codegangsta/cli"
    "os"
    "pfi/sensorbee/sensorbee/client"
    "pfi/sensorbee/sensorbee/server"
    "time"
{{range $_, $path := .Paths}}    _ "{{$path}}"
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
