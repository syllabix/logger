package encode

import (
	"go.uber.org/zap/zapcore"
)

// DevConsoleConfig is a bol.com tailored development optimized console encoder config
var DevConsoleConfig = &zapcore.EncoderConfig{
	MessageKey:     "message",
	LevelKey:       "level",
	EncodeLevel:    CapitalColorLevel,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	CallerKey:      "caller",
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	StacktraceKey:  "stacktrace",
}

// ProConsoleConfig is a bol.com tailored production optimized console encoder config
var ProConsoleConfig = &zapcore.EncoderConfig{
	MessageKey:     "message",
	LevelKey:       "level",
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	CallerKey:      "caller",
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	StacktraceKey:  "stacktrace",
}

// JSONConfig is a bol.com tailored production optimized json encoder config
var JSONConfig = zapcore.EncoderConfig{
	MessageKey:     "@message",
	LevelKey:       "level",
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	TimeKey:        "@timestamp",
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	CallerKey:      "caller",
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	StacktraceKey:  "stacktrace",
}
