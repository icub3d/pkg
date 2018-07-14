package telemetry

import "testing"

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
