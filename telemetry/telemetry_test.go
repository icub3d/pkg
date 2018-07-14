package telemetry

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		level string
		exp   Level
	}{
		{"fatal", FatalLevel},
		{"error", ErrorLevel},
		{"warn", WarnLevel},
		{"info", InfoLevel},
		{"debug", DebugLevel},
		{"Fatal", FatalLevel},
		{"Error", ErrorLevel},
		{"Warn", WarnLevel},
		{"Info", InfoLevel},
		{"Debug", DebugLevel},
		{"foo", DebugLevel},
	}
	for _, test := range tests {
		t.Run(test.level, func(t *testing.T) {
			l := ParseLevel(test.level)
			if l != test.exp {
				t.Errorf("ParseLevel(%v) = %v; expected %v", test.level, l, test.exp)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		exp   string
		level Level
	}{
		{"fatal", FatalLevel},
		{"error", ErrorLevel},
		{"warn", WarnLevel},
		{"info", InfoLevel},
		{"debug", DebugLevel},
		{"debug", 10},
	}
	for _, test := range tests {
		t.Run(test.exp, func(t *testing.T) {
			l := test.level.String()
			if l != test.exp {
				t.Errorf("%v.String() = %v; expected %v", test.level, l, test.exp)
			}
		})
	}

}

func TestDefaultTelemeterWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	SetTelemeter(NewWriter(DebugLevel, buf))
	w := WithFields(MetaData{"first": 1}).
		WithFields(MetaData{"second": 2})
	w.Debugf("%s", "debug")
	if !strings.Contains(buf.String(), "first") &&
		!strings.Contains(buf.String(), "second") {
		t.Errorf("couldn't find 'first' and 'second' in: %v", buf.String())
	}
}

func TestDefaultTelemeter(t *testing.T) {
	now = func() time.Time { return time.Unix(0, 0) }
	defer func() { now = time.Now }()
	exit = func(code int) {}
	defer func() { exit = os.Exit }()

	buf := &bytes.Buffer{}
	SetTelemeter(NewWriter(DebugLevel, buf))

	for x := DebugLevel; x >= 0; x-- {
		SetLevel(x)
		Debugf("%s", "debug")
		Infof("%s", "info")
		Warnf("%s", "warn")
		Errorf("%s", "error")
		Fatalf("%s", "fatal")
	}

	exp := `level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug"
level="info" when="1969-12-31 17:00:00 -0700 MST" what="info"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal"
level="info" when="1969-12-31 17:00:00 -0700 MST" what="info"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal"
`

	if buf.String() != exp {
		t.Errorf("log mismatch. got:\n%vexpected:\n%v", buf.String(), exp)
	}
}
