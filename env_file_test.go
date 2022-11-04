package forge_test

import (
	"os"
	"testing"

	"github.com/kyberbits/forge"
)

func TestSetValueInEnvFile(t *testing.T) {
	targetFile := "test_files/env_file/.env.local"
	os.Remove(targetFile)

	runtime := forge.NewRuntime()
	runtime.Environment = forge.Environment{} // Clear out the environment

	forge.SetValueInEnvFile(targetFile, "FOO", "FOO")
	forge.SetValueInEnvFile(targetFile, "FOO", "BAR")

	// Read in the current status
	if err := runtime.ReadInEnvironmentFile(targetFile); err != nil {
		t.Fatal(err)
	}

	if err := forge.Assert("BAR", runtime.Environment["FOO"]); err != nil {
		t.Fatal(err)
	}

}
