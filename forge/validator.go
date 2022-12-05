package forge

type ValidatorError struct {
	Message string
}

func (err ValidatorError) Error() string {
	return err.Message
}

type Validator interface {
	Validate() error
}
