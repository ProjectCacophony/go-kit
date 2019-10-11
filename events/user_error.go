package events

type UserError struct {
	Message string
	Err     error
}

func NewUserError(message string) *UserError {
	return &UserError{
		Message: message,
	}
}

func AsUserError(err error) *UserError {
	return &UserError{
		Message: err.Error(),
		Err:     err,
	}
}

func (e UserError) Error() string {
	return e.Message
}

func (e UserError) Unwrap() error {
	return e.Err
}
