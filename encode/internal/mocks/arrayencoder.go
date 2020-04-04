package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/buffer"
)

// PrimitiveArrayEncoder is a mocked implementation of the respective
// zap interface
type PrimitiveArrayEncoder struct {
	mock.Mock
	buf *buffer.Buffer
}

func (p *PrimitiveArrayEncoder) AppendBool(bool) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendByteString(message []byte) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendComplex128(complex128) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendComplex64(complex64) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendFloat64(float64) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendFloat32(float32) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendInt(int) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendInt64(int64) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendInt32(int32) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendInt16(int16) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendInt8(int8) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendString(str string) {
	p.Called(str)
}

func (p *PrimitiveArrayEncoder) AppendUint(uint) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendUint64(uint64) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendUint32(uint32) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendUint16(uint16) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendUint8(uint8) {
	panic("not implemented")
}

func (p *PrimitiveArrayEncoder) AppendUintptr(uintptr) {
	panic("not implemented")
}
