package forgetest

import (
	"encoding/json"
	"fmt"
	"reflect"
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

	if !reflect.DeepEqual(expected, actual) {
		actualBytes, _ := json.MarshalIndent(AssertFailure{
			Expected: expected,
			Actual:   actual,
		}, "", "\t")

		return fmt.Errorf("%s", string(actualBytes))
	}

	return nil
}
