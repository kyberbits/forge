package forge

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ValidatorError struct {
	Message string
}

func (err ValidatorError) Error() string {
	return err.Message
}

type Validator interface {
	Validate() error
}

func DecodeRequestBody(r *http.Request, target interface{}) error {
	err := json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		if err.Error() == "EOF" {
			return errors.New("Request Body Can Not Be Blank")
		}
		return err
	}

	if validator, ok := target.(Validator); ok {
		if err := validator.Validate(); err != nil {
			if _, ok := err.(ValidatorError); ok {
				return err
			}

			return ValidatorError{Message: err.Error()}
		}
	}

	return nil
}
