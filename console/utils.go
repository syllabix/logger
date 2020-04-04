package console

import (
	"github.com/syllabix/logger/encode"
	"go.uber.org/zap/zapcore"
)

func isEmpty(str string) bool {
	return len(str) < 1
}

func recolor(buffer []byte, lvl zapcore.Level) {
	i := 0
	for i < len(buffer) {
		b1 := buffer[i]
		if b1 == '\x1b' {
			i++
			b2 := buffer[i]
			if b2 == '[' {
				if buffer[i+1] != '0' {
					buffer[i+1], buffer[i+2] = encode.ColorBytesForLevel(lvl)
					i = i + 2
				}
			}
		}
		i++
	}
}
