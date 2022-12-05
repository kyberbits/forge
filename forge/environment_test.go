package forge_test

import (
	"testing"

	"github.com/kyberbits/forge/forge"
	"github.com/kyberbits/forge/forgetest"
)

func TestEnvironmentSuccess(t *testing.T) {
	type Config struct {
		Debug bool `env:"DEBUG"`
		Port  int  `env:"PORT"`
	}

	config := Config{
		Port: 22, // Default value
	}

	environment := forge.Environment{
		"PORT":  "33",   // Env Override
		"DEBUG": "true", // Env Override
	}

	if err := environment.Decode(&config); err != nil {
		t.Error(err)
	}

	actual := config
	expected := Config{
		Port:  33,
		Debug: true,
	}

	if err := forgetest.Assert(expected, actual); err != nil {
		t.Error(err)
	}
}
