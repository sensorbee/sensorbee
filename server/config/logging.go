package config

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io"
	"os"
)

// Logging has configuration parameters for logging.
type Logging struct {
	// Target is a logging target. It can be one of followings:
	//
	//	- stdout
	//	- stderr
	//	- file path
	Target string `json:"target" yaml:"target"`

	// MinLogLevel specifies the minimum level of log entries which should
	// actually be written to the target. Possible levels are "debug", "info",
	// "warn"/"warning", "error", or "fatal".
	MinLogLevel string `json:"min_log_level" yaml:"min_log_level"`

	// LogDroppedTuples controls logging of dropped tuples. If this parameter
	// is true, dropped tuples are logged as JSON objects in logs. It might
	// affect the overall performance of the server.
	LogDroppedTuples bool `json:"log_dropped_tuples" yaml:"log_dropped_tuples"`

	// SummarizeDroppedTuples controls summarization of dropped tuples. If this
	// parameter is true, only a portion of a dropped tuple is logged. The data
	// logged looks like a JSON object but might not be able to be parsed by
	// JSON parsers. This parameter only works when LogDroppedTuples is true.
	SummarizeDroppedTuples bool `json:"summarize_dropped_tuples" yaml:"summarize_dropped_tuples"`

	// TODO: add log rotation
	// TODO: add log formatting
}

var (
	loggingSchemaString = `{
	"type": "object",
	"properties": {
		"target": {
			"type": "string"
		},
		"min_log_level": {
			"enum": ["debug", "info", "warn", "warning", "error", "fatal"]
		},
		"log_dropped_tuples": {
			"type": "boolean"
		},
		"summarize_dropped_tuples": {
			"type": "boolean"
		}
	},
	"additionalProperties": false
}`
	loggingSchema *gojsonschema.Schema
)

func init() {
	s, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(loggingSchemaString))
	if err != nil {
		panic(err)
	}
	loggingSchema = s
}

// NewLogging create a Logging config parameters from a given map.
func NewLogging(m data.Map) (*Logging, error) {
	if err := validate(loggingSchema, m); err != nil {
		return nil, err
	}
	return newLogging(m), nil
}

func newLogging(m data.Map) *Logging {
	return &Logging{
		Target:                 mustAsString(getWithDefault(m, "target", data.String("stderr"))),
		MinLogLevel:            mustAsString(getWithDefault(m, "min_log_level", data.String("info"))),
		LogDroppedTuples:       mustToBool(getWithDefault(m, "log_dropped_tuples", data.False)),
		SummarizeDroppedTuples: mustToBool(getWithDefault(m, "summarize_dropped_tuples", data.False)),
	}
}

type nopCloser struct {
	w io.Writer
}

func (n *nopCloser) Write(p []byte) (int, error) {
	return n.w.Write(p)
}

func (n *nopCloser) Close() error {
	return nil
}

// CreateWriter creates io.Writer for loggers. When Target is a file, the writer
// supports log rotation using lumberjack.
func (l *Logging) CreateWriter() (io.WriteCloser, error) {
	// TODO: config package should probably concentrate on parsing and validating
	// config files and this should be moved to the server.
	switch l.Target {
	case "stdout":
		return &nopCloser{os.Stdout}, nil
	case "stderr":
		return &nopCloser{os.Stderr}, nil

	default:
		// Currently, only file path is supported
		f, err := os.OpenFile(l.Target, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return nil, fmt.Errorf("cannot open the file %v: %v", l.Target, err)
		}
		f.Close()

		return &lumberjack.Logger{
			Filename: l.Target,
			// TODO: set rotation options
		}, nil
	}
}
