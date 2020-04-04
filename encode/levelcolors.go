package encode

import (
	"strconv"

	"go.uber.org/zap/zapcore"
)

var (
	levelColors = map[zapcore.Level]Color{
		zapcore.DebugLevel:  Magenta,
		zapcore.InfoLevel:   Cyan,
		zapcore.WarnLevel:   Yellow,
		zapcore.ErrorLevel:  Red,
		zapcore.PanicLevel:  Red,
		zapcore.FatalLevel:  Red,
		zapcore.DPanicLevel: Red, // uber's development panic level
	}

	unknownLevelColor = Red

	levelcache = make(map[zapcore.Level]string, len(levelColors))
)

func init() {
	for level, color := range levelColors {
		levelcache[level] = color.Add(level.CapitalString())
	}
}

// ColorKey will color the provided key at the associated log level color
func ColorKey(key string, level zapcore.Level) string {
	c, ok := levelColors[level]
	if !ok {
		c = unknownLevelColor
	}
	return c.Add(key)
}

// CapitalColorLevel will apply coloring to the log level indicator
func CapitalColorLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	s, ok := levelcache[l]
	if !ok {
		s = unknownLevelColor.Add(l.CapitalString())
	}
	enc.AppendString(s)
}

// ColorBytesForLevel returns the byte sequence used to color keys for console output
func ColorBytesForLevel(lvl zapcore.Level) (byte, byte) {
	color, exists := levelColors[lvl]
	if !exists {
		color = Red
	}
	strcolor := strconv.Itoa(int(color))
	return strcolor[0], strcolor[1]
}
