package logger

import (
	"io"
	"os"

	"github.com/syllabix/logger/internal/registry"

	"github.com/syllabix/logger/mode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config has settings that are globally applied to all
// logging instances. It can be configured via a call to Configure
// with a list of options
type Config struct {
	mode mode.Kind
	// console/local sink
	csink io.Writer
	// json sink
	jsink   io.Writer
	appname string
	level   zapcore.Level
}

// sane defaults
var global = &Config{
	mode:    mode.Development,
	csink:   os.Stdout,
	jsink:   nil,
	appname: "",
	level:   zap.InfoLevel,
}

// An Option can be used to apply a value to a setting
// on the global config
type Option func(config *Config)

// AppName sets the "application" field to the provided value
// on all logging contexts
func AppName(name string) Option {
	return func(config *Config) {
		config.appname = name
	}
}

// ConsoleWriter sets the writer that will receive console formatted output
// from a logger
func ConsoleWriter(w io.Writer) Option {
	return func(config *Config) {
		config.csink = w
	}
}

// JSONWriter sets the writer that will receive json formatted output
// from a logger
func JSONWriter(w io.Writer) Option {
	return func(config *Config) {
		config.jsink = w
	}
}

// Mode sets the kind of mode loggers and their respective encoders
// should run in
func Mode(m mode.Kind) Option {
	return func(config *Config) {
		config.mode = m
	}
}

// Level sets the default log level of all logger instances
func Level(lvl zapcore.Level) Option {
	return func(config *Config) {
		config.level = lvl
	}
}

// Configure will apply all the supplied options to a global configuration
// that will be applied to all logger instances.
func Configure(options ...Option) {
	for _, opt := range options {
		opt(global)
	}

	registry.SetDefaultLevel(global.level)
}
