package encode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/syllabix/logger/encode/internal/mocks"
	"go.uber.org/zap/zapcore"
)

func TestColorKey(t *testing.T) {
	type args struct {
		key   string
		level zapcore.Level
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "info",
			args: args{
				key:   "INFO",
				level: zapcore.InfoLevel,
			},
			want: "\x1b[36mINFO\x1b[0m",
		},
		{
			name: "debug",
			args: args{
				key:   "DEBUG",
				level: zapcore.DebugLevel,
			},
			want: "\x1b[35mDEBUG\x1b[0m",
		},
		{
			name: "warn",
			args: args{
				key:   "WARN",
				level: zapcore.WarnLevel,
			},
			want: "\x1b[33mWARN\x1b[0m",
		},
		{
			name: "warn",
			args: args{
				key:   "WARN",
				level: zapcore.WarnLevel,
			},
			want: "\x1b[33mWARN\x1b[0m",
		},
		{
			name: "error",
			args: args{
				key:   "ERROR",
				level: zapcore.ErrorLevel,
			},
			want: "\x1b[31mERROR\x1b[0m",
		},
		{
			name: "unknown",
			args: args{
				key:   "CRAZY",
				level: zapcore.Level(22),
			},
			want: "\x1b[31mCRAZY\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ColorKey(tt.args.key, tt.args.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

var config = &zapcore.EncoderConfig{
	MessageKey:   "message",
	LevelKey:     "level",
	EncodeLevel:  zapcore.CapitalLevelEncoder,
	EncodeTime:   zapcore.ISO8601TimeEncoder,
	CallerKey:    "caller",
	EncodeCaller: zapcore.ShortCallerEncoder,
}

func makeEncoder() *mocks.PrimitiveArrayEncoder {
	mockEncoder := new(mocks.PrimitiveArrayEncoder)
	mockEncoder.On("AppendString", mock.Anything)
	return mockEncoder
}

func TestCapitalColorLevel(t *testing.T) {

	type args struct {
		l   zapcore.Level
		enc zapcore.PrimitiveArrayEncoder
	}
	tests := []struct {
		name      string
		args      args
		assertion func(t *testing.T, mock *mocks.PrimitiveArrayEncoder)
	}{
		{
			name: "debug",
			args: args{
				l: zapcore.DebugLevel,
			},
			assertion: func(t *testing.T, mock *mocks.PrimitiveArrayEncoder) {
				mock.AssertCalled(t, "AppendString", "\x1b[35mDEBUG\x1b[0m")
			},
		},
		{
			name: "info",
			args: args{
				l: zapcore.InfoLevel,
			},
			assertion: func(t *testing.T, mock *mocks.PrimitiveArrayEncoder) {
				mock.AssertCalled(t, "AppendString", "\x1b[36mINFO\x1b[0m")
			},
		},
		{
			name: "warn",
			args: args{
				l: zapcore.WarnLevel,
			},
			assertion: func(t *testing.T, mock *mocks.PrimitiveArrayEncoder) {
				mock.AssertCalled(t, "AppendString", "\x1b[33mWARN\x1b[0m")
			},
		},
		{
			name: "error",
			args: args{
				l: zapcore.ErrorLevel,
			},
			assertion: func(t *testing.T, mock *mocks.PrimitiveArrayEncoder) {
				mock.AssertCalled(t, "AppendString", "\x1b[31mERROR\x1b[0m")
			},
		},
		{
			name: "unknown level",
			args: args{
				l: zapcore.Level(19),
			},
			assertion: func(t *testing.T, mock *mocks.PrimitiveArrayEncoder) {
				mock.AssertCalled(t, "AppendString", "\x1b[31mLEVEL(19)\x1b[0m")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := makeEncoder()
			tt.args.enc = encoder
			CapitalColorLevel(tt.args.l, tt.args.enc)
			tt.assertion(t, encoder)
		})
	}
}

func TestColorBytesForLevel(t *testing.T) {
	type args struct {
		lvl zapcore.Level
	}
	tests := []struct {
		name  string
		args  args
		want  byte
		want1 byte
	}{
		{
			name: "debug",
			args: args{
				lvl: zapcore.DebugLevel,
			},
			want:  '3',
			want1: '5',
		},
		{
			name: "info",
			args: args{
				lvl: zapcore.InfoLevel,
			},
			want:  '3',
			want1: '6',
		},
		{
			name: "warn",
			args: args{
				lvl: zapcore.WarnLevel,
			},
			want:  '3',
			want1: '3',
		},
		{
			name: "error",
			args: args{
				lvl: zapcore.ErrorLevel,
			},
			want:  '3',
			want1: '1',
		},
		{
			name: "unknown",
			args: args{
				lvl: zapcore.Level(19),
			},
			want:  '3',
			want1: '1',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ColorBytesForLevel(tt.args.lvl)
			assert.Equal(t, tt.want, got, "first byte was not expected")
			assert.Equal(t, tt.want1, got1, "second byte was not expected")
		})
	}
}
