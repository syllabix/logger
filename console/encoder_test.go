package console

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syllabix/logger/encode"
	"github.com/syllabix/logger/mode"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var a_config = &zapcore.EncoderConfig{
	MessageKey:   "message",
	LevelKey:     "level",
	EncodeLevel:  zapcore.CapitalLevelEncoder,
	EncodeTime:   zapcore.ISO8601TimeEncoder,
	CallerKey:    "caller",
	EncodeCaller: zapcore.ShortCallerEncoder,
}

var b_config = &zapcore.EncoderConfig{
	MessageKey:   "@msg",
	LevelKey:     "lvl",
	NameKey:      "name_key",
	EncodeLevel:  zapcore.LowercaseLevelEncoder,
	EncodeTime:   zapcore.EpochTimeEncoder,
	CallerKey:    "trace",
	EncodeCaller: zapcore.ShortCallerEncoder,
}

var c_config = &zapcore.EncoderConfig{
	MessageKey:   "message",
	NameKey:      "name_key",
	EncodeLevel:  zapcore.LowercaseLevelEncoder,
	TimeKey:      "@timestamp",
	EncodeTime:   zapcore.EpochMillisTimeEncoder,
	CallerKey:    "trace",
	EncodeCaller: zapcore.ShortCallerEncoder,
}

func TestEncoder_Clone(t *testing.T) {
	type fields struct {
		config *zapcore.EncoderConfig
		buf    *buffer.Buffer
		level  zapcore.Level
		mode   mode.Kind
	}
	tests := []struct {
		name   string
		fields fields
		want   zapcore.Encoder
	}{
		{
			name: "clone empty buffer",
			fields: fields{
				config: a_config,
				buf:    bufferpool.Get(),
				level:  zapcore.InfoLevel,
				mode:   mode.Development,
			},
			want: &Encoder{
				config: a_config,
				buf:    bufferpool.Get(),
				level:  zapcore.InfoLevel,
				mode:   mode.Development,
			},
		},
		{
			name: "clone existing buffer contents",
			fields: fields{
				config: a_config,
				buf: func() *buffer.Buffer {
					b := bufferpool.Get()
					b.AppendString("message=ik denk we hebben een kleine problem")
					return b
				}(),
				level: zapcore.InfoLevel,
				mode:  mode.Development,
			},
			want: &Encoder{
				config: a_config,
				buf: func() *buffer.Buffer {
					b := bufferpool.Get()
					b.AppendString("message=ik denk we hebben een kleine problem")
					return b
				}(),
				level: zapcore.InfoLevel,
				mode:  mode.Development,
			},
		},
		{
			name: "zero value clone",
			fields: fields{
				config: b_config,
				buf:    bufferpool.Get(),
			},
			want: &Encoder{
				config: b_config,
				buf:    bufferpool.Get(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Encoder{
				config: tt.fields.config,
				buf:    tt.fields.buf,
				level:  tt.fields.level,
				mode:   tt.fields.mode,
			}
			if got := e.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncoder_devmode(t *testing.T) {
	type fields struct {
		config *zapcore.EncoderConfig
		buf    *buffer.Buffer
		level  zapcore.Level
		mode   mode.Kind
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "in dev mode",
			fields: fields{
				mode: mode.Development,
			},
			want: true,
		},
		{
			name: "in pro mode",
			fields: fields{
				mode: mode.Production,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Encoder{
				config: tt.fields.config,
				buf:    tt.fields.buf,
				level:  tt.fields.level,
				mode:   tt.fields.mode,
			}
			if got := e.devmode(); got != tt.want {
				t.Errorf("Encoder.devmode() = %v, want %v", got, tt.want)
			}
		})
	}
}

var dev_cfg = Config{
	Config: a_config,
	Mode:   mode.Development,
}

var pro_cfg = Config{
	Config: a_config,
	Mode:   mode.Production,
}

func TestEncoder_write(t *testing.T) {
	tests := []struct {
		name   string
		enc    *Encoder
		buffer []byte
		want   string
	}{
		{
			name: "dev mode debug level",
			enc: func() *Encoder {
				e := NewEncoder(dev_cfg)
				e.level = zapcore.DebugLevel
				return e
			}(),
			buffer: func() []byte {
				b := fmt.Sprintf("%s=dev %s=core-app %s=los.12314",
					encode.Cyan.Add("env"),
					encode.Cyan.Add("app"),
					encode.Cyan.Add("host"))

				return []byte(b)
			}(),
			want: "\x1b[35menv\x1b[0m=dev \x1b[35mapp\x1b[0m=core-app \x1b[35mhost\x1b[0m=los.12314",
		},
		{
			name: "dev mode info level",
			enc: func() *Encoder {
				e := NewEncoder(dev_cfg)
				e.level = zapcore.InfoLevel
				return e
			}(),
			buffer: func() []byte {
				b := fmt.Sprintf("%s=dev %s=core-app %s=los.12314",
					encode.Red.Add("env"),
					encode.Red.Add("app"),
					encode.Red.Add("host"))

				return []byte(b)
			}(),
			want: "\x1b[36menv\x1b[0m=dev \x1b[36mapp\x1b[0m=core-app \x1b[36mhost\x1b[0m=los.12314",
		},
		{
			name: "dev mode warn level",
			enc: func() *Encoder {
				e := NewEncoder(dev_cfg)
				e.level = zapcore.WarnLevel
				return e
			}(),
			buffer: func() []byte {
				b := fmt.Sprintf("%s=dev %s=core-app %s=los.12314",
					encode.Cyan.Add("env"),
					encode.Cyan.Add("app"),
					encode.Cyan.Add("host"))

				return []byte(b)
			}(),
			want: "\x1b[33menv\x1b[0m=dev \x1b[33mapp\x1b[0m=core-app \x1b[33mhost\x1b[0m=los.12314",
		},
		{
			name: "dev mode error level",
			enc: func() *Encoder {
				e := NewEncoder(dev_cfg)
				e.level = zapcore.ErrorLevel
				return e
			}(),
			buffer: func() []byte {
				b := fmt.Sprintf("%s=dev %s=core-app %s=los.12314",
					encode.Cyan.Add("env"),
					encode.Cyan.Add("app"),
					encode.Cyan.Add("host"))

				return []byte(b)
			}(),
			want: "\x1b[31menv\x1b[0m=dev \x1b[31mapp\x1b[0m=core-app \x1b[31mhost\x1b[0m=los.12314",
		},
		{
			name: "pro mode info level",
			enc: func() *Encoder {
				e := NewEncoder(pro_cfg)
				e.level = zapcore.InfoLevel
				return e
			}(),
			buffer: func() []byte {
				b := "env=dev app=core-app host=los.12314"

				return []byte(b)
			}(),
			want: "env=dev app=core-app host=los.12314",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.enc.write(tt.buffer)
			actual := tt.enc.buf.String()
			assert.Equal(t, tt.want, actual)
		})
	}
}

var info_entry = zapcore.Entry{
	Level:      zapcore.InfoLevel,
	Time:       time.Date(2020, time.March, 22, 13, 42, 12, 8, time.UTC),
	LoggerName: "test_log",
	Message:    "hello world, this is a log",
	Caller:     zapcore.NewEntryCaller(uintptr(123), "foo.go", 18, true),
	Stack:      "callstack...",
}

var debug_entry = zapcore.Entry{
	Level:      zapcore.DebugLevel,
	Time:       time.Date(2020, time.March, 22, 13, 42, 12, 8, time.UTC),
	LoggerName: "hello_loggr",
	Message:    "excellent day for a bike ride",
	Caller:     zapcore.NewEntryCaller(uintptr(123), "buzz.go", 18, true),
	Stack:      "callstack...",
}

var warn_entry = zapcore.Entry{
	Level:      zapcore.WarnLevel,
	Time:       time.Date(2020, time.March, 22, 13, 42, 12, 8, time.UTC),
	LoggerName: "hello_loggr",
	Message:    "excellent day for a bike ride",
	Caller:     zapcore.NewEntryCaller(uintptr(123), "buzz.go", 18, true),
	Stack:      "callstack...",
}

var error_entry = zapcore.Entry{
	Level:      zapcore.ErrorLevel,
	Time:       time.Date(2020, time.March, 22, 13, 42, 12, 8, time.UTC),
	LoggerName: "hello_loggr",
	Message:    "excellent day for a bike ride",
	Caller:     zapcore.NewEntryCaller(uintptr(123), "buzz.go", 18, true),
	Stack:      "callstack...",
}

func TestEncoder_EncodeEntry(t *testing.T) {
	type fields struct {
		config *zapcore.EncoderConfig
		buf    *buffer.Buffer
		level  zapcore.Level
		mode   mode.Kind
	}
	type args struct {
		ent    zapcore.Entry
		fields []zapcore.Field
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *buffer.Buffer
		wantErr bool
	}{
		{
			name: "info dev mode",
			fields: fields{
				config: a_config,
				buf:    bufferpool.Get(),
				mode:   mode.Development,
			},
			args: args{
				ent: info_entry,
				fields: []zapcore.Field{
					zapcore.Field{
						Key:       "application",
						Type:      zapcore.StringType,
						Integer:   0,
						String:    "bol.kit",
						Interface: nil,
					},
					zapcore.Field{
						Key:       "errcount",
						Type:      zapcore.Int32Type,
						Integer:   230,
						String:    "",
						Interface: nil,
					},
				},
			},
			want: func() *buffer.Buffer {
				b := bufferpool.Get()
				str := "INFO 2020-03-22T13:42:12.000Z [36mcaller[0m=foo.go:18 [36mmessage[0m=hello world, this is a log [36mapplication[0m=bol.kit [36merrcount[0m=230\n"
				b.Write([]byte(str))
				return b
			}(),
			wantErr: false,
		},
		{
			name: "info pro mode",
			fields: fields{
				config: a_config,
				buf:    bufferpool.Get(),
				mode:   mode.Production,
			},
			args: args{
				ent: info_entry,
				fields: []zapcore.Field{
					zapcore.Field{
						Key:       "application",
						Type:      zapcore.StringType,
						Integer:   0,
						String:    "bol.kit",
						Interface: nil,
					},
					zapcore.Field{
						Key:       "errcount",
						Type:      zapcore.Int32Type,
						Integer:   230,
						String:    "",
						Interface: nil,
					},
				},
			},
			want: func() *buffer.Buffer {
				b := bufferpool.Get()
				str := "INFO 2020-03-22T13:42:12.000Z caller=foo.go:18 message=hello world, this is a log application=bol.kit errcount=230\n"
				b.Write([]byte(str))
				return b
			}(),
			wantErr: false,
		},
		{
			name: "debug pro mode",
			fields: fields{
				config: b_config,
				buf:    bufferpool.Get(),
				mode:   mode.Production,
			},
			args: args{
				ent: debug_entry,
				fields: []zapcore.Field{
					zapcore.Field{
						Key:       "@source_host",
						Type:      zapcore.StringType,
						Integer:   0,
						String:    "xavier-910",
						Interface: nil,
					},
					zapcore.Field{
						Key:       "errcount",
						Type:      zapcore.Int32Type,
						Integer:   122,
						String:    "",
						Interface: nil,
					},
				},
			},
			want: func() *buffer.Buffer {
				b := bufferpool.Get()
				str := "debug 1584884532 name_key=hello_loggr trace=buzz.go:18 @msg=excellent day for a bike ride @source_host=xavier-910 errcount=122\n"
				b.Write([]byte(str))
				return b
			}(),
			wantErr: false,
		},
		{
			name: "warn dev mode",
			fields: fields{
				config: a_config,
				buf:    bufferpool.Get(),
				mode:   mode.Development,
			},
			args: args{
				ent: warn_entry,
				fields: []zapcore.Field{
					zapcore.Field{
						Key:       "@source_host",
						Type:      zapcore.StringType,
						Integer:   0,
						String:    "xavier-910",
						Interface: nil,
					},
					zapcore.Field{
						Key:       "errcount",
						Type:      zapcore.Int32Type,
						Integer:   122,
						String:    "",
						Interface: nil,
					},
				},
			},
			want: func() *buffer.Buffer {
				b := bufferpool.Get()
				str := "WARN 2020-03-22T13:42:12.000Z \x1b[33mcaller\x1b[0m=buzz.go:18 \x1b[33mmessage\x1b[0m=excellent day for a bike ride \x1b[33m@source_host\x1b[0m=xavier-910 \x1b[33merrcount\x1b[0m=122\n"
				b.Write([]byte(str))
				return b
			}(),
			wantErr: false,
		},
		{
			name: "warn pro mode",
			fields: fields{
				config: a_config,
				buf:    bufferpool.Get(),
				mode:   mode.Production,
			},
			args: args{
				ent: warn_entry,
				fields: []zapcore.Field{
					zapcore.Field{
						Key:       "@source_host",
						Type:      zapcore.StringType,
						Integer:   0,
						String:    "xavier-910",
						Interface: nil,
					},
					zapcore.Field{
						Key:       "errcount",
						Type:      zapcore.Int32Type,
						Integer:   122,
						String:    "",
						Interface: nil,
					},
				},
			},
			want: func() *buffer.Buffer {
				b := bufferpool.Get()
				str := "WARN 2020-03-22T13:42:12.000Z caller=buzz.go:18 message=excellent day for a bike ride @source_host=xavier-910 errcount=122\n"
				b.Write([]byte(str))
				return b
			}(),
			wantErr: false,
		},
		{
			name: "error dev mode",
			fields: fields{
				config: c_config,
				buf:    bufferpool.Get(),
				mode:   mode.Development,
			},
			args: args{
				ent: error_entry,
				fields: []zapcore.Field{
					zapcore.Field{
						Key:       "@source_host",
						Type:      zapcore.StringType,
						Integer:   0,
						String:    "xavier-910",
						Interface: nil,
					},
					zapcore.Field{
						Key:       "errcount",
						Type:      zapcore.Int32Type,
						Integer:   122,
						String:    "",
						Interface: nil,
					},
				},
			},
			want: func() *buffer.Buffer {
				b := bufferpool.Get()
				str := "error @timestamp=1584884532000 \x1b[31mname_key\x1b[0m=hello_loggr \x1b[31mtrace\x1b[0m=buzz.go:18 \x1b[31mmessage\x1b[0m=excellent day for a bike ride \x1b[31m@source_host\x1b[0m=xavier-910 \x1b[31merrcount\x1b[0m=122\n"
				b.Write([]byte(str))
				return b
			}(),
			wantErr: false,
		},
		{
			name: "error pro mode",
			fields: fields{
				config: c_config,
				buf:    bufferpool.Get(),
				mode:   mode.Production,
			},
			args: args{
				ent: error_entry,
				fields: []zapcore.Field{
					zapcore.Field{
						Key:       "@source_host",
						Type:      zapcore.StringType,
						Integer:   0,
						String:    "xavier-910",
						Interface: nil,
					},
					zapcore.Field{
						Key:       "errcount",
						Type:      zapcore.Int32Type,
						Integer:   122,
						String:    "",
						Interface: nil,
					},
				},
			},
			want: func() *buffer.Buffer {
				b := bufferpool.Get()
				str := "error @timestamp=1584884532000 name_key=hello_loggr trace=buzz.go:18 message=excellent day for a bike ride @source_host=xavier-910 errcount=122\n"
				b.Write([]byte(str))
				return b
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Encoder{
				config: tt.fields.config,
				buf:    tt.fields.buf,
				level:  tt.fields.level,
				mode:   tt.fields.mode,
			}

			got, err := e.EncodeEntry(tt.args.ent, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder.EncodeEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}
