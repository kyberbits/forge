package forgetest

import (
	"fmt"

	"github.com/go-test/deep"
)

type AssertFailure struct {
	Expected interface{}
	Actual   interface{}
}

func Assert(expected interface{}, actual interface{}) error {
	if actualErr, ok := actual.(error); ok {
		actual = actualErr.Error()
	}

	if expectedErr, ok := expected.(error); ok {
		expected = expectedErr.Error()
	}

	if diff := deep.Equal(expected, actual); diff != nil {
		return fmt.Errorf("%v", diff)
	}

	return nil
}
