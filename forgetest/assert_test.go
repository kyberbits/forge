package forgetest_test

import (
	"errors"
	"testing"

	"github.com/kyberbits/forge/forgetest"
)

func TestAssert(t *testing.T) {
	// LOL, this test is funny!

	{ // Equal
		expected := errors.New("{\n\t\"Expected\": false,\n\t\"Actual\": true\n}")
		actual := forgetest.Assert(false, true)
		if err := forgetest.Assert(expected, actual); err != nil {
			t.Error(err)
		}
	}

	{ // Not Equal
		var expected error = nil
		actual := forgetest.Assert(true, true)
		if err := forgetest.Assert(expected, actual); err != nil {
			t.Error(err)
		}
	}
}
