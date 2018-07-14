package telemetry

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

func NewWriter(l Level, w io.Writer) Telemeter {
	return &writer{
		lock: &sync.Mutex{},
		w:    w,
		l:    l,
	}
}

type writer struct {
	lock *sync.Mutex
	w    io.Writer
	l    Level
	md   MetaData
}

func (w *writer) log(l Level, m string) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	buf.WriteString(fmt.Sprintf(`level="%s" when="%s" what="%s"`,
		l.String(), now(), m))
	for k, v := range w.md {
		buf.WriteString(fmt.Sprintf(` %s="%v"`, k, v))
	}
	buf.WriteString("\n")

	w.lock.Lock()
	defer w.lock.Unlock()

	_, err := w.w.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("failed to write to log: %v", err)
	}
}

func (w *writer) WithFields(md MetaData) Telemeter {
	n := MetaData{}
	for k, v := range w.md {
		n[k] = v
	}
	for k, v := range md {
		n[k] = v
	}
	return &writer{
		lock: w.lock,
		w:    w.w,
		l:    w.l,
		md:   n,
	}
}

func (w *writer) Fatalf(format string, args ...interface{}) {
	if w.l >= FatalLevel {
		w.log(FatalLevel, fmt.Sprintf(format, args...))
	}
	exit(1)
}

func (w *writer) Errorf(format string, args ...interface{}) {
	if w.l >= ErrorLevel {
		w.log(ErrorLevel, fmt.Sprintf(format, args...))
	}
}

func (w *writer) Warnf(format string, args ...interface{}) {
	if w.l >= WarnLevel {
		w.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (w *writer) Infof(format string, args ...interface{}) {
	if w.l >= InfoLevel {
		w.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (w *writer) Debugf(format string, args ...interface{}) {
	if w.l >= DebugLevel {
		w.log(DebugLevel, fmt.Sprintf(format, args...))
	}
}
