package logger

import (
	"os"
	"runtime"
	"strings"

	"github.com/syllabix/logger/internal/registry"
)

// GetHostname returns hostname
func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return h
}

func pkgname() registry.Package {
	pc, _, _, _ := runtime.Caller(2)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	pkg := ""

	if parts[pl-2][0] == '(' {
		pkg = strings.Join(parts[0:pl-2], ".")
	} else {
		pkg = strings.Join(parts[0:pl-1], ".")
	}
	return registry.Package(pkg)
}
