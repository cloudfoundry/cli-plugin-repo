package test_helpers

import "strings"

type TestLogger struct {
	contents []string
	cursor   int
}

func NewTestLogger() *TestLogger {
	return &TestLogger{
		cursor: 0,
	}
}

func (l *TestLogger) Write(p []byte) (n int, err error) {
	l.contents = append(l.contents, string(p))
	return len(p), nil
}

func (l *TestLogger) ContainsSubstring(strs []string) bool {
	for _, str := range strs {
		for l.cursor < len(l.contents) {
			if strings.Contains(l.contents[l.cursor], str) {
				break
			}
			l.cursor = l.cursor + 1
		}
	}

	if l.cursor == len(l.contents) {
		return false
	}
	return true
}
