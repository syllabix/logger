package console

import (
	"sync"

	"go.uber.org/zap"

	"github.com/syllabix/logger/mode"

	"go.uber.org/zap/buffer"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &Encoder{}
	},
}

func get() *Encoder {
	return pool.Get().(*Encoder)
}

func put(enc *Encoder) {
	enc.config = nil
	enc.buf = nil
	enc.mode = mode.None
	enc.level = zap.InfoLevel
	pool.Put(enc)
}

var bufferpool = buffer.NewPool()
