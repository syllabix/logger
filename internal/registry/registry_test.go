package registry

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestSetDefaultLevel(t *testing.T) {
	type args struct {
		level zapcore.Level
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "info",
			args: args{
				level: zap.InfoLevel,
			},
		},
		{
			name: "debug",
			args: args{
				level: zap.DebugLevel,
			},
		},
		{
			name: "warn",
			args: args{
				level: zap.WarnLevel,
			},
		},
		{
			name: "error",
			args: args{
				level: zap.ErrorLevel,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefaultLevel(tt.args.level)
			assert.Equal(t, defaultLevel, tt.args.level)
		})
	}
}

func TestGet(t *testing.T) {

	type args struct {
		pkg Package
	}
	tests := []struct {
		name  string
		args  args
		setup func()
		want  zap.AtomicLevel
	}{
		{
			name: "new package @ default",
			args: args{
				pkg: "my/package/name",
			},
			setup: func() {
				defaultLevel = zap.InfoLevel
			},
			want: zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name: "existing package @ initial level",
			args: args{
				pkg: "my/package/name",
			},
			setup: func() {
				defaultLevel = zap.WarnLevel
			},
			want: zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name: "existing package @ updated",
			args: args{
				pkg: "my/package/name",
			},
			setup: func() {
				Set("my/package/name", zap.DebugLevel)
			},
			want: zap.NewAtomicLevelAt(zap.DebugLevel),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := Get(tt.args.pkg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	type args struct {
		pkg   Package
		level zapcore.Level
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    zapcore.Level
		setup   func()
	}{
		{
			name: "non registered",
			args: args{
				pkg:   "not/registered/pkg",
				level: zap.InfoLevel,
			},
			setup:   func() { /*no-op*/ },
			wantErr: true,
		},
		{
			name: "update registered package",
			args: args{
				pkg:   "my/package/name",
				level: zap.WarnLevel,
			},
			setup: func() {
				Get("my/package/name")
			},
			want:    zap.WarnLevel,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Set(tt.args.pkg, tt.args.level); (err != nil) != tt.wantErr {
				tt.setup()
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)

				if !tt.wantErr {
					assert.Equal(t, tt.want, Get(tt.args.pkg))
				}
			}
		})
	}
}

func TestGetPackages(t *testing.T) {

	tests := []struct {
		name  string
		setup func()
		want  []Package
	}{
		{
			name: "empty",
			setup: func() {
				levels = make(map[Package]zap.AtomicLevel)
			},
			want: []Package{},
		},
		{
			name: "all registered",
			setup: func() {
				levels = make(map[Package]zap.AtomicLevel)
				levels[Package("alpha")] = zap.NewAtomicLevelAt(zapcore.InfoLevel)
				levels[Package("beta")] = zap.NewAtomicLevelAt(zapcore.WarnLevel)
				levels[Package("kappa")] = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
			},
			want: []Package{"alpha", "beta", "kappa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := GetPackages()
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
