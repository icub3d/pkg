package telemetry

// NewMulti creates a telemeter that sends events to all of the given
// telementers.
//
// NOTE: Fatalf will probably not send to all of the telemeters
// because they are likely to call os.Exit.
func NewMulti(tt ...Telemeter) Telemeter {
	return &multi{
		tt: tt,
	}
}

type multi struct {
	tt []Telemeter
	md MetaData
}

func (m *multi) SetLevel(level Level) {
	for _, t := range m.tt {
		t.SetLevel(level)
	}
}

func (m *multi) WithFields(md MetaData) Telemeter {
	return &multi{
		tt: m.tt,
		md: md,
	}
}

func (m *multi) Fatalf(format string, args ...interface{}) {
	for _, t := range m.tt {
		t.WithFields(m.md).Fatalf(format, args...)
	}
}

func (m *multi) Errorf(format string, args ...interface{}) {
	for _, t := range m.tt {
		t.WithFields(m.md).Errorf(format, args...)
	}
}

func (m *multi) Warnf(format string, args ...interface{}) {
	for _, t := range m.tt {
		t.WithFields(m.md).Warnf(format, args...)
	}
}

func (m *multi) Infof(format string, args ...interface{}) {
	for _, t := range m.tt {
		t.WithFields(m.md).Infof(format, args...)
	}
}

func (m *multi) Debugf(format string, args ...interface{}) {
	for _, t := range m.tt {
		t.WithFields(m.md).Debugf(format, args...)
	}
}
