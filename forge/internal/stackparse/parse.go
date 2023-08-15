package stackparse

import (
	"strings"
)

type StackTrace struct {
	Name   string
	Frames []Frame
}

type Frame struct {
	Call     string
	Location string
}

func Parse(stBytes []byte) StackTrace {
	st := StackTrace{}

	lines := strings.Split(string(stBytes), "\n")
	for lineNumber, line := range lines {
		if lineNumber+1 == len(lines) {
			continue
		}

		if lineNumber == 0 {
			st.Name = line

			continue
		}

		if lineNumber%2 == 0 {
			continue
		}

		st.Frames = append(st.Frames, Frame{
			Call:     line,
			Location: strings.TrimSpace(lines[lineNumber+1]),
		})
	}

	return st
}
