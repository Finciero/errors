package errors

// BuildError ...
func BuildError(err error) *Error {
	if err, ok := (err).(*Error); ok {
		return err
	}
	return InternalServerFromError(err)
}
