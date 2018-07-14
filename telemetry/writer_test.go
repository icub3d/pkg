package telemetry

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestWriteWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(DebugLevel, buf).WithFields(MetaData{"first": 1}).
		WithFields(MetaData{"second": 2})
	w.Debugf("%s", "debug")
	if !strings.Contains(buf.String(), "first") &&
		!strings.Contains(buf.String(), "second") {
		t.Errorf("couldn't find 'first' and 'second' in: %v", buf.String())
	}
}

func TestWriter(t *testing.T) {
	now = func() time.Time { return time.Unix(0, 0) }
	defer func() { now = time.Now }()
	exit = func(code int) {}
	defer func() { exit = os.Exit }()

	buf := &bytes.Buffer{}
	w := NewWriter(DebugLevel, buf)

	for x := DebugLevel; x >= 0; x-- {
		w.SetLevel(x)
		t := w.WithFields(MetaData{
			"cat": true,
		})
		t.Debugf("%s", "debug")
		t.Infof("%s", "info")
		t.Warnf("%s", "warn")
		t.Errorf("%s", "error")
		t.Fatalf("%s", "fatal")
	}

	exp := `level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug" cat="true"
level="info" when="1969-12-31 17:00:00 -0700 MST" what="info" cat="true"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn" cat="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" cat="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" cat="true"
level="info" when="1969-12-31 17:00:00 -0700 MST" what="info" cat="true"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn" cat="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" cat="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" cat="true"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn" cat="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" cat="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" cat="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" cat="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" cat="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" cat="true"
`

	if buf.String() != exp {
		t.Errorf("log mismatch. got:\n%vexpected:\n%v", buf.String(), exp)
	}
}
