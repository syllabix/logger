package logger

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syllabix/logger/mode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestAppName(t *testing.T) {

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple name",
			args: args{
				name: "my-app",
			},
			want: "my-app",
		},
		{
			name: "empty name",
			args: args{
				name: "",
			},
			want: "",
		},
		{
			name: "simple name",
			args: args{
				name: "my-app",
			},
			want: "my-app",
		},
		{
			name: "long name",
			args: args{
				name: "my-incredible-application-that-will-enable-massive-profits-and-fame-for-all",
			},
			want: "my-incredible-application-that-will-enable-massive-profits-and-fame-for-all",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := new(Config)
			opt := AppName(tt.args.name)
			opt(cfg)
			assert.Equal(t, tt.want, cfg.appname)
		})
	}
}

func TestConsoleWriter(t *testing.T) {

	w := new(bytes.Buffer)
	type args struct {
		w io.Writer
	}
	tests := []struct {
		name string
		args args
		want io.Writer
	}{
		{
			name: "assigns writer",
			args: args{w},
			want: w,
		},
		{
			name: "assigns nil",
			args: args{nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := new(Config)
			opt := ConsoleWriter(tt.args.w)
			opt(cfg)
			assert.Equal(t, tt.want, cfg.csink)
		})
	}
}

func TestJSONWriter(t *testing.T) {
	w := new(bytes.Buffer)
	type args struct {
		w io.Writer
	}
	tests := []struct {
		name string
		args args
		want io.Writer
	}{
		{
			name: "assigns writer",
			args: args{w},
			want: w,
		},
		{
			name: "assigns nil",
			args: args{nil},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := new(Config)
			opt := JSONWriter(tt.args.w)
			opt(cfg)
			assert.Equal(t, tt.want, cfg.jsink)
		})
	}
}

func TestMode(t *testing.T) {
	type args struct {
		mode mode.Kind
	}
	tests := []struct {
		name string
		args args
		want mode.Kind
	}{
		{
			name: "dev",
			args: args{
				mode: mode.Development,
			},
			want: mode.Development,
		},
		{
			name: "pro",
			args: args{
				mode: mode.Production,
			},
			want: mode.Production,
		},
		{
			name: "none",
			args: args{
				mode: mode.None,
			},
			want: mode.None,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := new(Config)
			opt := Mode(tt.args.mode)
			opt(cfg)
			assert.Equal(t, tt.want, cfg.mode)
		})
	}
}

func TestLevel(t *testing.T) {
	type args struct {
		lvl zapcore.Level
	}
	tests := []struct {
		name string
		args args
		want zapcore.Level
	}{
		{
			name: "info",
			args: args{
				lvl: zap.InfoLevel,
			},
			want: zap.InfoLevel,
		},
		{
			name: "warn",
			args: args{
				lvl: zap.WarnLevel,
			},
			want: zap.WarnLevel,
		},
		{
			name: "error",
			args: args{
				lvl: zap.ErrorLevel,
			},
			want: zap.ErrorLevel,
		},
		{
			name: "debug",
			args: args{
				lvl: zap.DebugLevel,
			},
			want: zap.DebugLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := new(Config)
			opt := Level(tt.args.lvl)
			opt(cfg)
			assert.Equal(t, tt.want, cfg.level)
		})
	}
}

func TestConfigure(t *testing.T) {
	before()
	defer after()

	writer := new(bytes.Buffer)

	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			name: "defaults",
			args: args{
				options: make([]Option, 0, 1),
			},
			want: Config{
				mode:    mode.Development,
				csink:   os.Stdout,
				jsink:   nil,
				appname: "",
				level:   zap.InfoLevel,
			},
		},
		{
			name: "app_name",
			args: args{
				options: []Option{
					AppName("awesome-app"),
				},
			},
			want: Config{
				mode:    mode.Development,
				csink:   os.Stdout,
				jsink:   nil,
				appname: "awesome-app",
				level:   zap.InfoLevel,
			},
		},
		{
			name: "app_name/mode/level",
			args: args{
				options: []Option{
					AppName("awesome-app"),
					Level(zap.ErrorLevel),
					Mode(mode.None),
				},
			},
			want: Config{
				mode:    mode.None,
				csink:   os.Stdout,
				jsink:   nil,
				appname: "awesome-app",
				level:   zap.ErrorLevel,
			},
		},
		{
			name: "app_name/mode/level/csink",
			args: args{
				options: []Option{
					AppName("awesome-app"),
					Level(zap.ErrorLevel),
					Mode(mode.None),
					ConsoleWriter(os.Stderr),
				},
			},
			want: Config{
				mode:    mode.None,
				csink:   os.Stderr,
				jsink:   nil,
				appname: "awesome-app",
				level:   zap.ErrorLevel,
			},
		},
		{
			name: "app_name/mode/level/csink/jsink",
			args: args{
				options: []Option{
					AppName("awesome-app"),
					Level(zap.ErrorLevel),
					Mode(mode.None),
					ConsoleWriter(os.Stderr),
					JSONWriter(writer),
				},
			},
			want: Config{
				mode:    mode.None,
				csink:   os.Stderr,
				jsink:   writer,
				appname: "awesome-app",
				level:   zap.ErrorLevel,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Configure(tt.args.options...)
			assert.Equal(t, tt.want, *global)
		})
	}
}
