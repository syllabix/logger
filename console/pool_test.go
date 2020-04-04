package console

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syllabix/logger/mode"
)

func Test_put(t *testing.T) {
	type args struct {
		enc *Encoder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "proper reset development",
			args: args{
				enc: NewEncoder(dev_cfg),
			},
		},
		{
			name: "proper reset from pro",
			args: args{
				enc: NewEncoder(pro_cfg),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			put(tt.args.enc)
			assert.Nil(t, tt.args.enc.buf, "buffer should be nil'd out on put to pool")
			assert.Nil(t, tt.args.enc.config, "config should be nil'd out on put to pool")
			assert.Equal(t, tt.args.enc.mode, mode.None, "mode should be set to none on put to pool")
		})
	}
}

func Test_get(t *testing.T) {
	tests := []struct {
		name string
		want *Encoder
	}{
		{
			name: "instance from pool is reset",
			want: &Encoder{
				mode: mode.None,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() = %v, want %v", got, tt.want)
			}
		})
	}
}
