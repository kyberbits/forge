package stackparse_test

import (
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/kyberbits/forge/forge/internal/stackparse"
)

func TestParse(t *testing.T) {
	want := 4
	if got := stackparse.Parse(debug.Stack()); !reflect.DeepEqual(len(got.Frames), want) {
		// e := json.NewEncoder(os.Stdout)
		// e.SetIndent("", "\t")
		// e.Encode(got)
		t.Errorf("Parse() = %v, want %v", len(got.Frames), want)
	}
}
