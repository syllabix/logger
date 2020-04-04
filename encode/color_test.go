package encode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor_Add(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		c    Color
		args args
		want string
	}{
		{
			name: "Black",
			c:    Black,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[30m@source_host\x1b[0m",
		},
		{
			name: "Red",
			c:    Red,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[31m@source_host\x1b[0m",
		},
		{
			name: "Green",
			c:    Green,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[32m@source_host\x1b[0m",
		},
		{
			name: "Yellow",
			c:    Yellow,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[33m@source_host\x1b[0m",
		},
		{
			name: "Blue",
			c:    Blue,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[34m@source_host\x1b[0m",
		},
		{
			name: "Magenta",
			c:    Magenta,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[35m@source_host\x1b[0m",
		},
		{
			name: "Cyan",
			c:    Cyan,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[36m@source_host\x1b[0m",
		},
		{
			name: "White",
			c:    White,
			args: args{
				s: "@source_host",
			},
			want: "\x1b[37m@source_host\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.Add(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
