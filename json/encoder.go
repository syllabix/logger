package json

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Encoder encodes log messages in logstash json format
type Encoder struct {
	enc zapcore.Encoder
}

// Clone implements the zapcore Encoder interface
func (e *Encoder) Clone() zapcore.Encoder {
	encoder := e.enc.Clone()
	return &Encoder{encoder}
}

// EncodeEntry encodes the log entry in logstash json format
func (e *Encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {

	cfield, defined := caller(ent)
	if defined {
		fields = append(fields, cfield)
	}

	fields = append(fields, level(ent))
	return e.enc.EncodeEntry(ent, fields)
}

// NewEncoder returns a json Encoder
func NewEncoder(cfg zapcore.EncoderConfig) *Encoder {
	return &Encoder{
		enc: zapcore.NewJSONEncoder(cfg),
	}
}

func level(ent zapcore.Entry) zapcore.Field {
	return zapcore.Field{
		Key:    "level",
		Type:   zapcore.StringType,
		String: ent.Level.String(),
	}
}

func caller(ent zapcore.Entry) (zapcore.Field, bool) {
	if !ent.Caller.Defined {
		return zapcore.Field{}, false
	}

	return zapcore.Field{
		Key:    "caller",
		Type:   zapcore.StringType,
		String: ent.Caller.String(),
	}, true
}
