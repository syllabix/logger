package registry

import (
	"errors"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ErrPkgNotRegistered is returned when a package is not in the registry
var ErrPkgNotRegistered = errors.New("the provided package is not in the registry")

// Package is string representation of a Go package
type Package string

var (
	levels       = make(map[Package]zap.AtomicLevel)
	defaultLevel = zapcore.InfoLevel
	mutex        sync.Mutex
)

// SetDefaultLevel sets the default log level applied to the registry
func SetDefaultLevel(level zapcore.Level) {
	defaultLevel = level
}

// Set a log level for logger instances in the provided package
// Set returns a non nil error
func Set(pkg Package, level zapcore.Level) error {
	lvl, ok := levels[pkg]
	if !ok {
		return ErrPkgNotRegistered
	}
	lvl.SetLevel(level)
	return nil
}

// Get returns the atomic log level for the provided package name
func Get(pkg Package) zap.AtomicLevel {
	mutex.Lock()
	defer mutex.Unlock()
	lvl, ok := levels[pkg]
	if !ok {
		lvl = zap.NewAtomicLevelAt(defaultLevel)
		levels[pkg] = lvl
	}
	return lvl
}

// GetPackages returns all packages in the registry
func GetPackages() []Package {
	pkgs := make([]Package, 0, len(levels))
	for pkg := range levels {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
