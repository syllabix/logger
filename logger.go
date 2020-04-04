package logger

import (
	"github.com/syllabix/logger/console"
	"github.com/syllabix/logger/encode"
	"github.com/syllabix/logger/internal/registry"
	"github.com/syllabix/logger/json"
	"github.com/syllabix/logger/mode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func consoleConfig() console.Config {
	config := console.Config{
		Mode: global.mode,
	}

	if global.mode == mode.Production {
		config.Config = encode.ProConsoleConfig
	} else {
		config.Config = encode.DevConsoleConfig
	}
	return config
}

// New returns an instance of a logger configured via the logger package
// global options
func New() *zap.Logger {

	level := registry.Get(pkgname())

	// configure console encoder
	cEncoder := console.NewEncoder(consoleConfig())
	cout := zapcore.AddSync(global.csink)
	core := zapcore.NewCore(cEncoder, cout, level)

	// if a json sink has been set, configure it
	// with the zapcore JSON Encoder, and Tee the console
	// encoder with it
	if global.jsink != nil {
		jEncoder := json.NewEncoder(encode.JSONConfig)
		rsink := zapcore.AddSync(global.jsink)
		core = zapcore.NewTee(
			core,
			zapcore.NewCore(jEncoder, rsink, level),
		)
	}

	// TODO: determine the most efficient way to allocate
	// core = zapcore.NewSampler(
	// 	core,
	// 	500*time.Millisecond,
	// 	10, // first
	// 	50, // thereafter
	// )

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.PanicLevel)).
		With(zap.String("@source_host", hostname()))

	logger = logger.With(zap.Namespace("@fields"))

	if len(global.appname) > 0 {
		logger = logger.With(zap.String("application", global.appname))
	}

	return logger
}

// SetLevelForPackage will set the log level for all instances
// of a logger in the provided package, returning an error if the
// package name provided does not exist in the registry
func SetLevelForPackage(pkg string, level zapcore.Level) error {
	return registry.Set(registry.Package(pkg), level)
}

// GetPackages retuns all package names that logger instances
// have been created in
func GetPackages() []string {
	pkgs := registry.GetPackages()
	strpkgs := make([]string, len(pkgs))
	for i := range pkgs {
		strpkgs[i] = string(pkgs[i])
	}
	return strpkgs
}
