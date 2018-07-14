package telemetry

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestMulti(t *testing.T) {
	now = func() time.Time { return time.Unix(0, 0) }
	defer func() { now = time.Now }()
	exit = func(code int) {}
	defer func() { exit = os.Exit }()

	buf := &bytes.Buffer{}
	w := NewMulti(
		NewWriter(DebugLevel, buf),
		NewWriter(InfoLevel, buf),
		NewWriter(WarnLevel, buf),
		NewWriter(ErrorLevel, buf),
		NewWriter(FatalLevel, buf),
	).WithFields(MetaData{"top": true})
	w.Debugf("%s", "debug")
	w.Infof("%s", "info")
	w.Warnf("%s", "warn")
	w.Errorf("%s", "error")
	w.Fatalf("%s", "fatal")

	w.SetLevel(DebugLevel)
	w.Debugf("%s", "debug")

	exp := `level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug" top="true"
level="info" when="1969-12-31 17:00:00 -0700 MST" what="info" top="true"
level="info" when="1969-12-31 17:00:00 -0700 MST" what="info" top="true"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn" top="true"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn" top="true"
level="warn" when="1969-12-31 17:00:00 -0700 MST" what="warn" top="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" top="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" top="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" top="true"
level="error" when="1969-12-31 17:00:00 -0700 MST" what="error" top="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" top="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" top="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" top="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" top="true"
level="fatal" when="1969-12-31 17:00:00 -0700 MST" what="fatal" top="true"
level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug" top="true"
level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug" top="true"
level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug" top="true"
level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug" top="true"
level="debug" when="1969-12-31 17:00:00 -0700 MST" what="debug" top="true"
`

	if buf.String() != exp {
		t.Errorf("log mismatch. got:\n%vexpected:\n%v", buf.String(), exp)
	}
}
