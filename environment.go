package forge

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrInvalidValue returned when the value passed to Unmarshal is nil or not a pointer to a struct.
	ErrInvalidValue = errors.New("value must be a non-nil pointer to a struct")

	// ErrUnsupportedFieldType returned when a field with tag "env" is unsupported.
	ErrUnsupportedFieldType = errors.New("field is an unsupported type")

	// ErrUnexportedField returned when a field with tag "env" is not exported.
	ErrUnexportedField = errors.New("field must be exported")
)

func NewEnvironment() Environment {
	environment := Environment{}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		environment[pair[0]] = pair[1]
	}

	return environment
}

type Environment map[string]string

func (env Environment) Decode(target interface{}) error {
	valueOf := reflect.ValueOf(target)
	element := valueOf.Elem()
	for i := 0; i < element.NumField(); i++ {
		fieldInstance := element.Field(i)
		fieldDefinition := element.Type().Field(i)

		// Get the tag value
		tag := fieldDefinition.Tag.Get("env")
		if tag == "" {
			continue
		}

		// Get the matching environment variable
		valueFromEnv, ok := env[tag]
		if !ok {
			continue
		}

		switch fieldDefinition.Type.Kind() {
		case reflect.String:
			fieldInstance.SetString(valueFromEnv)
		case reflect.Bool:
			convertedValue, err := strconv.ParseBool(valueFromEnv)
			if err != nil {
				return err
			}
			fieldInstance.SetBool(convertedValue)

		case reflect.Int:
			convertedValue, err := strconv.Atoi(valueFromEnv)
			if err != nil {
				return err
			}
			fieldInstance.SetInt(int64(convertedValue))
		}
	}

	return nil
}

func (environment Environment) ImportEnvFileContents(fileContents string) error {
	lines := strings.Split(fileContents, "\n")
	for _, line := range lines {
		// Skip blank lines
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		value := parts[1]

		// Skip the env variable is already set
		if _, alreadySet := environment[key]; alreadySet {
			continue
		}

		// Set the value
		environment[key] = value
	}

	return nil
}
