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

// bufferPool is used by the internal loggers to store string buffers
// that can be used during logging. It helps from an allocation
// perspective because the buffers keep their memory around. Just make
// sure you Reset() the buffer before you use it and then return it
// when you are done with it.
var bufferPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// Level is the event levels that the Telemeters will register and
// process. Setting a Level for a Telemeter tells it to process all
// Levels of equal or higher priority.
type Level int

const (
	// FatalLevel is the highest priority and also suggests that the
	// client application should exit with os.Exit(1) after sending
	// the event.
	FatalLevel Level = iota

	// ErrorLevel suggests that the event is related to an error.
	ErrorLevel

	// WarnLevel suggests that the event is not an error but
	// something that should be dealth with.
	WarnLevel

	// InfoLevel suggests that the event contains information that may
	// be useful.
	InfoLevel

	// DebugLevel suggests that the event contains low level
	// information that probably shouldn't be shown in a production
	// environment but would be more useful for development.
	DebugLevel
)

// ParseLevel turns the given string into a Level. The valuse are
// "fatal", "error", "warn", "info", and "debug". DebugLevel is
// returned if the given string doesn't match.
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

// String implements the fmt.Stringer interface. It prints the level
// in lowercase. An invalid level returns "debug".
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

// MetaData contains additional information about an event.
type MetaData map[string]interface{}

// Event contains information about an event that was
// fired. Telemeters can use this to process information but it's
// mostly useful if you need to transmit the information elsewhere
// using some encoder.
type Event struct {
	Level    Level
	When     time.Time
	What     string
	MetaData MetaData
}

// Telemeter is an instrument for processing events.
type Telemeter interface {

	// SetLevel changes the level at which messages are processed.
	SetLevel(Level)

	// WithFields returns a Telemeter that will include the given
	// MetaData with each event. If the current Telemeter already has
	// MetaData, this should be merged with the existing MetaData.
	WithFields(MetaData) Telemeter

	// Fatalf fires an event with the given message at the FatalLevel.
	Fatalf(format string, args ...interface{})

	// Errorf fires an event with the given message at the ErrorLevel.
	Errorf(format string, args ...interface{})

	// Warnf fires an event with the given message at the WarnLevel.
	Warnf(format string, args ...interface{})

	// Infof fires an event with the given message at the InfoLevel.
	Infof(format string, args ...interface{})

	// Debugf fires an event with the given message at the DebugLevel.
	Debugf(format string, args ...interface{})
}

var t = NewWriter(DebugLevel, os.Stdout)

// SetTelemeter sets the packages default Telemeter to the given
// Telemeter. If not called, events are sent to os.Stdout.
func SetTelemeter(telemeter Telemeter) {
	t = telemeter
}

// SetLevel calls SetLevel on the default Telementer.
func SetLevel(level Level) {
	t.SetLevel(level)
}

// WithFields calls WithFields on the default Telementer.
func WithFields(md MetaData) Telemeter {
	return t.WithFields(md)
}

// Fatalf calls Fatalf on the default Telemeter.
func Fatalf(format string, args ...interface{}) {
	t.Fatalf(format, args...)
}

// Errorf calls Errorf on the default Telemeter.
func Errorf(format string, args ...interface{}) {
	t.Errorf(format, args...)
}

// Warnf calls Warnf on the default Telemeter.
func Warnf(format string, args ...interface{}) {
	t.Warnf(format, args...)
}

// Infof calls Infof on the default Telemeter.
func Infof(format string, args ...interface{}) {
	t.Infof(format, args...)
}

// Debugf calls Debugf on the default Telemeter.
func Debugf(format string, args ...interface{}) {
	t.Debugf(format, args...)
}
