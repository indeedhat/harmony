package common

type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}

// WrapError into another error
func WrapError(original, aditional error) error {
	if aditional == nil {
		return original
	}

	if original == nil {
		return aditional
	}

	return &wrapError{
		msg: original.Error(),
		err: aditional,
	}
}
