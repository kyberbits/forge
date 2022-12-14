package forge_test

import (
	"os"
	"testing"

	"github.com/kyberbits/forge/forge"
	"github.com/kyberbits/forge/forgetest"
)

func TestRuntime(t *testing.T) {
	runtime := forge.NewRuntime()
	runtime.FS = os.DirFS("./test_files/environment")

	runtime.Environment = forge.Environment{
		"SYSTEM": "system",
	}
	if err := runtime.ReadInDefaultEnvironmentFiles(); err != nil {
		panic(err)
	}

	actual := runtime.Environment
	expected := forge.Environment{
		"SYSTEM":  "system",
		"DEFAULT": "default",
		"LOCAL":   "local",
	}

	if err := forgetest.Assert(expected, actual); err != nil {
		t.Error(err)
	}
}
