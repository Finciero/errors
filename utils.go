package errors

// BuildError ...
func BuildError(err error) *Error {
	if err == nil {
		return nil
	}

	if err, ok := (err).(*Error); ok {
		return err
	}

	return InternalServerFromError(err, "unexpected error")
}
