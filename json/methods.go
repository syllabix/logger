package json

import (
	"time"

	"go.uber.org/zap/zapcore"
)

func (e *Encoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	return e.enc.AddArray(key, marshaler)
}

func (e *Encoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	return e.enc.AddObject(key, marshaler)
}

func (e *Encoder) AddBinary(key string, value []byte) {
	e.enc.AddBinary(key, value)
}

func (e *Encoder) AddByteString(key string, value []byte) {
	e.enc.AddByteString(key, value)
}

func (e *Encoder) AddBool(key string, value bool) {
	e.enc.AddBool(key, value)
}

func (e *Encoder) AddComplex128(key string, value complex128) {
	e.enc.AddComplex128(key, value)
}

func (e *Encoder) AddComplex64(key string, value complex64) {
	e.enc.AddComplex64(key, value)
}

func (e *Encoder) AddDuration(key string, value time.Duration) {
	e.enc.AddDuration(key, value)
}

func (e *Encoder) AddFloat64(key string, value float64) {
	e.enc.AddFloat64(key, value)
}

func (e *Encoder) AddFloat32(key string, value float32) {
	e.enc.AddFloat32(key, value)
}

func (e *Encoder) AddInt(key string, value int) {
	e.enc.AddInt(key, value)
}

func (e *Encoder) AddInt64(key string, value int64) {
	e.enc.AddInt64(key, value)
}

func (e *Encoder) AddInt32(key string, value int32) {
	e.enc.AddInt32(key, value)
}

func (e *Encoder) AddInt16(key string, value int16) {
	e.enc.AddInt16(key, value)
}

func (e *Encoder) AddInt8(key string, value int8) {
	e.enc.AddInt8(key, value)
}

func (e *Encoder) AddString(key string, value string) {
	e.enc.AddString(key, value)
}

func (e *Encoder) AddTime(key string, value time.Time) {
	e.enc.AddTime(key, value)
}

func (e *Encoder) AddUint(key string, value uint) {
	e.enc.AddUint(key, value)
}

func (e *Encoder) AddUint64(key string, value uint64) {
	e.enc.AddUint64(key, value)
}

func (e *Encoder) AddUint32(key string, value uint32) {
	e.enc.AddUint32(key, value)
}

func (e *Encoder) AddUint16(key string, value uint16) {
	e.enc.AddUint16(key, value)
}

func (e *Encoder) AddUint8(key string, value uint8) {
	e.enc.AddUint8(key, value)
}

func (e *Encoder) AddUintptr(key string, value uintptr) {
	e.enc.AddUintptr(key, value)
}

func (e *Encoder) AddReflected(key string, value interface{}) error {
	return e.enc.AddReflected(key, value)
}

func (e *Encoder) OpenNamespace(key string) {
	e.enc.OpenNamespace(key)
}
