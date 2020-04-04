package console

import (
	"github.com/syllabix/logger/encode"
	"github.com/syllabix/logger/mode"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Config is an Encoder purposed configuration struct
// containing a reference to the zacore EncoderConfig as well
// as the encoder mode it should run in
type Config struct {
	Config *zapcore.EncoderConfig
	Mode   mode.Kind
}

// Encoder is a bol.com tailored zap encoder for
// writing human readable logs to the console
type Encoder struct {
	config *zapcore.EncoderConfig
	buf    *buffer.Buffer
	level  zapcore.Level
	mode   mode.Kind
}

// Clone implements the Clone method of the zapcore Encoder interface
func (e *Encoder) Clone() zapcore.Encoder {
	clone := e.clone(e.level)
	clone.buf.Write(e.buf.Bytes())
	return clone
}

func (e *Encoder) clone(level zapcore.Level) *Encoder {
	clone := get()
	clone.config = e.config
	clone.level = level
	clone.mode = e.mode
	clone.buf = bufferpool.Get()
	return clone
}

func (e *Encoder) devmode() bool {
	return e.mode == mode.Development
}

func (e *Encoder) write(buffer []byte) {
	if e.devmode() {
		recolor(buffer, e.level)
	}
	e.buf.Write(buffer)
}

func (e *Encoder) addKey(key string) {
	e.buf.AppendByte(' ')
	if e.devmode() {
		e.buf.AppendString(encode.ColorKey(key, e.level))
	} else {
		e.buf.AppendString(key)
	}
	e.buf.AppendByte('=')
}

// EncodeEntry implements the EncodeEntry method of the zapcore Encoder interface
func (e *Encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := e.clone(ent.Level)
	config := final.config

	config.EncodeLevel(ent.Level, final)
	final.buf.AppendByte(' ')

	if !isEmpty(config.TimeKey) {
		final.AppendString(config.TimeKey)
		final.buf.AppendByte('=')
	}

	config.EncodeTime(ent.Time, final)

	if !isEmpty(ent.LoggerName) && !isEmpty(config.NameKey) {
		final.addKey(config.NameKey)
		cur := final.buf.Len()
		nameEncoder := config.EncodeName
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}
		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			final.AppendString(ent.LoggerName)
		}
	}

	if ent.Caller.Defined && !isEmpty(config.CallerKey) {
		final.addKey(config.CallerKey)
		cur := final.buf.Len()
		config.EncodeCaller(ent.Caller, final)
		if cur == final.buf.Len() {
			final.AppendString(ent.Caller.String())
		}
	}

	if !isEmpty(config.MessageKey) {
		final.addKey(config.MessageKey)
		final.AppendString(ent.Message)
	}

	for i := range fields {
		fields[i].AddTo(final)
	}

	if final.buf.Len() > 0 {
		final.write(e.buf.Bytes())
	}

	if !isEmpty(ent.Stack) && !isEmpty(config.StacktraceKey) {
		final.AddString(config.StacktraceKey, ent.Stack)
	}

	if !isEmpty(config.LineEnding) {
		final.buf.AppendString(config.LineEnding)
	} else {
		final.buf.AppendString("\n")
	}

	ret := final.buf
	put(final)

	return ret, nil
}

// Register is used to register the console Encoder to the zap framework
func Register(mode mode.Kind) func(zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return func(zcfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		c := Config{
			Mode:   mode,
			Config: &zcfg,
		}
		return NewEncoder(c), nil
	}
}

// NewEncoder initializes a a bol.com tailored Encoder
func NewEncoder(cfg Config) *Encoder {
	return &Encoder{
		buf:    bufferpool.Get(),
		mode:   cfg.Mode,
		config: cfg.Config,
	}
}
