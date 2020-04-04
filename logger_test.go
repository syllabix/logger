package logger

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syllabix/logger/console"
	"github.com/syllabix/logger/encode"
	"github.com/syllabix/logger/internal/registry"
	"github.com/syllabix/logger/mode"
	"go.uber.org/zap"
)

var tmpconfig Config

func before() {
	tmpconfig = *global
}

func after() {
	*global = tmpconfig
}

const testStateError = "as the tests in this package share global state, it appears another test has corrupted the expected state for this one - please use the before() and after() funcs to reset state"

func Test_consoleConfig(t *testing.T) {
	before()
	defer after()

	tests := []struct {
		name  string
		setup func()
		want  console.Config
	}{
		{
			name: "pro",
			setup: func() {
				global.mode = mode.Production
			},
			want: console.Config{
				Mode:   mode.Production,
				Config: encode.ProConsoleConfig,
			},
		},
		{
			name: "dev",
			setup: func() {
				global.mode = mode.Development
			},
			want: console.Config{
				Mode:   mode.Development,
				Config: encode.DevConsoleConfig,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := consoleConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("consoleConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPackages(t *testing.T) {
	before()
	defer after()

	tests := []struct {
		name  string
		setup func()
		want  []string
	}{
		{
			name:  "empty",
			setup: func() { /* */ },
			want:  []string{},
		},
		{
			name: "returns registered pkgs",
			setup: func() {
				registry.Get("my/awesome/pkg")
				registry.Get("main")
				registry.Get("cool/code/man")
				registry.Get("sometimes/there/is/a/thing/that/just/feels/slightly/ever/so/slightly/too-much")
			},
			want: []string{
				"my/awesome/pkg",
				"main",
				"cool/code/man",
				"sometimes/there/is/a/thing/that/just/feels/slightly/ever/so/slightly/too-much",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := GetPackages()
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

// discarder is used to capture output written to an io.Writer
type discarder struct {
	log string
}

func (d *discarder) Write(p []byte) (n int, err error) {
	d.log = string(p)
	return len(p), nil
}

func (d *discarder) reset() {
	d.log = ""
}

func correctFormat(timestamp string) bool {
	_, err := time.Parse("2006-01-02T15:04:05.000Z0700", timestamp)
	return err == nil
}

func TestNew(t *testing.T) {
	before()
	defer after()

	consolew := new(discarder)
	jsonw := new(discarder)

	tests := []struct {
		name       string
		setup      func()
		checkinfo  func(t *testing.T)
		checkwarn  func(t *testing.T)
		checkerror func(t *testing.T)
	}{
		{
			name: "dev config with sensible defaults",
			setup: func() {
				Configure(
					AppName("test-app"),
					ConsoleWriter(consolew),
				)
			},
			checkinfo: func(t *testing.T) {
				output := strings.Split(consolew.log, " ")

				if len(output) < 13 {
					t.Error(testStateError)
					t.FailNow()
				}

				assert.Equal(t, "\x1b[36mINFO\x1b[0m", output[0])
				assert.True(t, correctFormat(output[1]))
				assert.Equal(t, "\x1b[36mcaller\x1b[0m=logger/logger_test.go:272", output[2])
				assert.Equal(t, "\x1b[36mmessage\x1b[0m=hello", output[3])
				assert.Equal(t, "\x1b[36mstatus\x1b[0m=blue", output[10])
				assert.Equal(t, "\x1b[36mcount\x1b[0m=12", output[11])
				assert.Equal(t, "\x1b[36mapplication\x1b[0m=test-app\n", output[13])
			},
			checkwarn: func(t *testing.T) {
				output := strings.Split(consolew.log, " ")
				if len(output) < 13 {
					t.Error(testStateError)
					t.FailNow()
				}

				assert.Equal(t, "\x1b[33mWARN\x1b[0m", output[0])
				assert.True(t, correctFormat(output[1]))
				assert.Equal(t, "\x1b[33mcaller\x1b[0m=logger/logger_test.go:278", output[2])
				assert.Equal(t, "\x1b[33mmessage\x1b[0m=hello", output[3])
				assert.Equal(t, "\x1b[33mstatus\x1b[0m=yellow", output[10])
				assert.Equal(t, "\x1b[33mcount\x1b[0m=54", output[11])
				assert.Equal(t, "\x1b[33mapplication\x1b[0m=test-app\n", output[13])
			},
			checkerror: func(t *testing.T) {
				output := strings.Split(consolew.log, " ")
				if len(output) < 13 {
					t.Error(testStateError)
					t.FailNow()
				}

				assert.Equal(t, "\x1b[31mERROR\x1b[0m", output[0])
				assert.True(t, correctFormat(output[1]))
				assert.Equal(t, "\x1b[31mcaller\x1b[0m=logger/logger_test.go:284", output[2])
				assert.Equal(t, "\x1b[31mmessage\x1b[0m=hello", output[3])
				assert.Equal(t, "\x1b[31mstatus\x1b[0m=red", output[10])
				assert.Equal(t, "\x1b[31mcount\x1b[0m=9102", output[11])
				assert.Equal(t, "\x1b[31mapplication\x1b[0m=test-app\n", output[13])
			},
		},
		{
			name: "pro config with json and console writers",
			setup: func() {
				Configure(
					AppName("test-app"),
					ConsoleWriter(consolew),
					JSONWriter(jsonw),
					Mode(mode.Production),
				)
			},
			checkinfo: func(t *testing.T) {
				output := strings.Split(consolew.log, " ")
				joutput := jsonw.log

				if len(output) < 13 {
					t.Error(testStateError)
					t.FailNow()
				}

				// TODO: assert json output format
				assert.NotEmpty(t, joutput)

				assert.Equal(t, "INFO", output[0])
				assert.True(t, correctFormat(output[1]))
				assert.Equal(t, "caller=logger/logger_test.go:272", output[2])
				assert.Equal(t, "message=hello", output[3])
				assert.Equal(t, "status=blue", output[10])
				assert.Equal(t, "count=12", output[11])
				assert.Equal(t, "application=test-app\n", output[13])
			},
			checkwarn: func(t *testing.T) {
				output := strings.Split(consolew.log, " ")
				joutput := jsonw.log

				if len(output) < 13 {
					t.Error(testStateError)
					t.FailNow()
				}

				// TODO: assert json output format
				assert.NotEmpty(t, joutput)

				assert.Equal(t, "WARN", output[0])
				assert.True(t, correctFormat(output[1]))
				assert.Equal(t, "caller=logger/logger_test.go:278", output[2])
				assert.Equal(t, "message=hello", output[3])
				assert.Equal(t, "status=yellow", output[10])
				assert.Equal(t, "count=54", output[11])
				assert.Equal(t, "application=test-app\n", output[13])
			},
			checkerror: func(t *testing.T) {
				output := strings.Split(consolew.log, " ")
				joutput := jsonw.log

				if len(output) < 13 {
					t.Error(testStateError)
					t.FailNow()
				}

				// TODO: assert json output format
				assert.NotEmpty(t, joutput)

				assert.Equal(t, "ERROR", output[0])
				assert.True(t, correctFormat(output[1]))
				assert.Equal(t, "caller=logger/logger_test.go:284", output[2])
				assert.Equal(t, "message=hello", output[3])
				assert.Equal(t, "status=red", output[10])
				assert.Equal(t, "count=9102", output[11])
				assert.Equal(t, "application=test-app\n", output[13])
			},
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			logger := New()
			logger.Info("hello world - this is some info",
				zap.String("status", "blue"),
				zap.Int("count", 12),
			)
			tt.checkinfo(t)

			logger.Warn("hello world - this is a warning",
				zap.String("status", "yellow"),
				zap.Int("count", 54),
			)
			tt.checkwarn(t)

			logger.Error("hello world - this is an error",
				zap.String("status", "red"),
				zap.Int("count", 9102),
			)
			tt.checkerror(t)
		})
	}
}
