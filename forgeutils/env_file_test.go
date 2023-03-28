package forgeutils_test

import (
	"os"
	"testing"

	"github.com/kyberbits/forge/forge"
	"github.com/kyberbits/forge/forgetest"
	"github.com/kyberbits/forge/forgeutils"
)

func TestSetValueInEnvFile(t *testing.T) {
	targetFile := "test_files/env_file/.env.local"
	os.Remove(targetFile)

	runtime := forge.NewRuntime()
	runtime.Environment = forge.Environment{} // Clear out the environment

	forgeutils.EnvironmentSetValueInFile(targetFile, "FOO", "FOO")
	forgeutils.EnvironmentSetValueInFile(targetFile, "FOO", "BAR")

	// Read in the current status
	if err := runtime.ReadInEnvironmentFile(targetFile); err != nil {
		t.Fatal(err)
	}

	if err := forgetest.Assert("BAR", runtime.Environment["FOO"]); err != nil {
		t.Fatal(err)
	}
}
