// Package telemetry provides a method for logging.
package telemetry

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	now  = time.Now
	exit = os.Exit
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

type Level int

const (
	FatalLevel Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

// DebugLevel is returned if the given string doesn't match.
func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	default:
		return DebugLevel
	}
}

func (l Level) String() string {
	switch l {
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	case WarnLevel:
		return "warn"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	default:
		return "debug"
	}
}

type MetaData map[string]interface{}

type Event struct {
	Level    Level
	When     time.Time
	What     string
	MetaData MetaData
}

type Telemeter interface {
	// Should merge with existing.
	WithFields(MetaData) Telemeter
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}
