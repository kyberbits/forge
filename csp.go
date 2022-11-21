package forge

import (
	"fmt"
	"strings"
)

type CSP struct {
	Default []string
	Script  []string
	Connect []string
	Frame   []string
}

func (csp CSP) String() string {
	result := ""
	if len(csp.Default) > 0 {
		result += fmt.Sprintf("default-src %s;", strings.Join(csp.Default, " "))
	}

	if len(csp.Script) > 0 {
		result += fmt.Sprintf("script-src %s;", strings.Join(csp.Script, " "))
	}

	if len(csp.Connect) > 0 {
		result += fmt.Sprintf("connect-src %s;", strings.Join(csp.Connect, " "))
	}

	if len(csp.Frame) > 0 {
		result += fmt.Sprintf("frame-src %s;", strings.Join(csp.Frame, " "))
	}

	return result
}
