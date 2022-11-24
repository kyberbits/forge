package forge_test

import (
	"errors"
	"testing"

	"github.com/kyberbits/forge"
)

func TestAssert(t *testing.T) {
	// LOL, this test is funny!

	{ // Equal
		expected := errors.New("{\n\t\"Expected\": false,\n\t\"Actual\": true\n}")
		actual := forge.Assert(false, true)
		if err := forge.Assert(expected, actual); err != nil {
			t.Error(err)
		}
	}

	{ // Not Equal
		var expected error = nil
		actual := forge.Assert(true, true)
		if err := forge.Assert(expected, actual); err != nil {
			t.Error(err)
		}
	}
}
