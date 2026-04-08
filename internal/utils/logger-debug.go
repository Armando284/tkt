//go:build debug

package utils

import (
	"fmt"
	"time"
)

// Esta implementación SÓLO se compila cuando haces: go build -tags=debug
type debugLogger struct{}

// var Dev Logger = &debugLogger{}

func (l *debugLogger) Debugf(format string, args ...any) {
	l.log("DEBUG", format, args...)
}

func (l *debugLogger) Infof(format string, args ...any) {
	l.log("INFO", format, args...)
}

func (l *debugLogger) Warnf(format string, args ...any) {
	l.log("WARN", format, args...)
}

func (l *debugLogger) Errorf(format string, args ...any) {
	l.log("ERROR", format, args...)
}

func (l *debugLogger) log(level, format string, args ...any) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[tkt %s] %s: %s\n", level, timestamp, fmt.Sprintf(format, args...))
}