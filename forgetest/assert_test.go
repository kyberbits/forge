package forgetest_test

import (
	"errors"
	"testing"

	"github.com/kyberbits/forge/forgetest"
)

var errFoobar = errors.New("[false != true]")

func TestAssert(t *testing.T) {
	// LOL, this test is funny!
	{ // Equal
		expected := errFoobar
		actual := forgetest.Assert(false, true)
		if err := forgetest.Assert(expected, actual); err != nil {
			t.Error(err)
		}
	}

	{ // Not Equal
		var expected error
		actual := forgetest.Assert(true, true)
		if err := forgetest.Assert(expected, actual); err != nil {
			t.Error(err)
		}
	}
}
