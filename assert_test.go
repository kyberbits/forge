package forge_test

import (
	"errors"
	"testing"

	"github.com/kyberbits/forge"
)

func TestAssert(t *testing.T) {
	// LOL, this test is funny!

	{ // Equal
		expected := errors.New("not equal... {\"Expected\":false,\"Actual\":true}")
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
