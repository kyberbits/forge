package forge_test

import (
	"testing"

	"github.com/kyberbits/forge/forge"
	"github.com/kyberbits/forge/forgetest"
)

func TestCSPFull(t *testing.T) {
	csp := forge.CSP{
		Default: []string{
			"'self'",
			"example.com",
		},
		Script: []string{
			"'self'",
			"example.com",
		},
		Connect: []string{
			"'self'",
			"example.com",
		},
		Frame: []string{
			"'self'",
			"example.com",
		},
	}
	expected := "default-src 'self' example.com;script-src 'self' example.com;connect-src 'self' example.com;frame-src 'self' example.com;"
	if err := forgetest.Assert(expected, csp.String()); err != nil {
		t.Error(err)
	}
}

func TestCSPBlank(t *testing.T) {
	csp := forge.CSP{}
	expected := ""
	if err := forgetest.Assert(expected, csp.String()); err != nil {
		t.Error(err)
	}
}
