// go:generate stringer -type=Code
package errors

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Code type
type Code int32

// Codes identifiers
const (
	bad_request    Code = 400
	unauthorized   Code = 401
	delinquent     Code = 402
	forbidden      Code = 403
	not_found      Code = 404
	not_acceptable Code = 406
	invalid_params Code = 422
	rate_limit     Code = 429

	internal_server Code = 500
)

// Exportable aliases from real codes
const (
	StatusBadRequest          = bad_request
	StatusUnauthorized        = unauthorized
	StatusPaymentRequired     = delinquent
	StatusForbidden           = forbidden
	StatusNotFound            = not_found
	StatusNotAcceptable       = not_acceptable
	StatusUnprocessableEntity = invalid_params
	StatusTooManyRequests     = rate_limit

	StatusInternalServerError = internal_server
)

// Error type
type Error struct {
	StatusCode Code
	Meta       Meta
	Message    string

	InternalError error // internal information used for debugging
}

// Meta stores metadata that can be visible for end users and developers
type Meta map[string]interface{}

// New returns a new Error
func New(code Code, msg string, setters ...errorParamsSetter) *Error {
	var meta Meta
	for _, fn := range setters {
		fn(&meta)
	}
	return &Error{
		StatusCode: code,
		Meta:       meta,
		Message:    msg,
	}
}

// NewFromError returns a New Error with description of the error given
func NewFromError(code Code, err error, msg string, setters ...errorParamsSetter) *Error {
	var meta Meta
	for _, fn := range setters {
		fn(&meta)
	}
	return &Error{
		StatusCode: code,
		Meta:       meta,
		Message:    msg,

		InternalError: err,
	}
}

// FromGRPC returns a new Error from an error received by grpc. If the
// error was encoded with ToGPC method then the full Error passed is
// returned.
func FromGRPC(err error) *Error {
	var raw struct {
		Meta          Meta   `json:"meta, omitempty"`
		Message       string `json:"msg, omitempty"`
		InternalError error  `json:"internal_error,omitempty"`
	}

	code := grpc.Code(err)
	desc := grpc.ErrorDesc(err)

	if unmarshalError := json.Unmarshal([]byte(desc), &raw); unmarshalError != nil {
		return InternalServerFromError(err, "unexpected error")
	}

	return &Error{
		StatusCode: Code(code),
		Meta:       raw.Meta,
		Message:    raw.Message,

		InternalError: raw.InternalError,
	}
}

// ToGRPC ecode error into a grpc error
func (e *Error) ToGRPC() error {
	buff, _ := json.Marshal(struct {
		Meta    Meta   `json:"meta,omitempty"`
		Message string `json:"msg,omitempty"`

		InternalError error `json:"internal_error,omitempty"`
	}{
		Meta:    e.Meta,
		Message: e.Message,

		InternalError: e.InternalError,
	})

	return grpc.Errorf(codes.Code(e.StatusCode), string(buff))
}

// Code returns error StatusCode casted to int
func (e *Error) Code() int {
	return int(e.StatusCode)
}

// Error method return string representation of error.
func (e *Error) Error() string {
	str := fmt.Sprintf("status_code=%d error_id=%q", e.StatusCode, fmt.Sprint(e.StatusCode))

	if len(e.Message) > 0 {
		str += fmt.Sprintf(" msg=%q", e.Message)
	}

	if e.InternalError != nil {
		str += fmt.Sprintf(" desc=%q", e.InternalError.Error())
	}

	for key, value := range e.Meta {
		str += fmt.Sprintf(" %s=%q", key, value)
	}

	return str
}

// ErrorID returns string representation of the error StatusCode.
func (e *Error) ErrorID() string {
	return fmt.Sprint(e.StatusCode)
}

type errorParamsSetter func(*Meta)

// SetMeta sets the given key values into the Meta of the error.
func SetMeta(m Meta) errorParamsSetter {
	return func(params *Meta) {
		if (*params) == nil {
			(*params) = m
			return
		}

		for key, value := range m {
			(*params)[key] = value
		}
	}
}

// BadRequest returns an Error with bad_request code
func BadRequest(message string, setters ...errorParamsSetter) *Error {
	return New(StatusBadRequest, message, setters...)
}

// BadRequestFromError returns an Error with bad_request code with err as a
// internalError.
func BadRequestFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusBadRequest, err, msg, setters...)
}

// Unauthorized returns an Error with unauthorized code
func Unauthorized(message string, setters ...errorParamsSetter) *Error {
	return New(StatusUnauthorized, message, setters...)
}

// UnauthorizedFromError returns an Error with unauthorized code with err as a
// internalError.
func UnauthorizedFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusUnauthorized, err, msg, setters...)
}

// Delinquent returns an Error with delinquent code
func Delinquent(message string, setters ...errorParamsSetter) *Error {
	return New(StatusPaymentRequired, message, setters...)
}

// DelinquentFromError returns an Error with delinquent code with err as a
// internalError.
func DelinquentFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusPaymentRequired, err, msg, setters...)
}

// Forbidden returns an Error with forbidden code
func Forbidden(message string, setters ...errorParamsSetter) *Error {
	return New(StatusForbidden, message, setters...)
}

// ForbiddenFromError returns an Error with forbidden code with err as a
// internalError.
func ForbiddenFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusForbidden, err, msg, setters...)
}

// NotFound returns an Error with not_found code
func NotFound(message string, setters ...errorParamsSetter) *Error {
	return New(StatusNotFound, message, setters...)
}

// NotFoundFromError returns an Error with not_found code with err as a
// internalError.
func NotFoundFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusNotFound, err, msg, setters...)
}

// NotAcceptable returns an Error with not_acceptable code
func NotAcceptable(message string, setters ...errorParamsSetter) *Error {
	return New(StatusNotAcceptable, message, setters...)
}

// NotAcceptableFromError returns an Error with not_acceptable code with err as a
// internalError.
func NotAcceptableFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusNotAcceptable, err, msg, setters...)
}

// InvalidParams returns an Error with invalid_params code
func InvalidParams(message string, setters ...errorParamsSetter) *Error {
	return New(StatusUnprocessableEntity, message, setters...)
}

// InvalidParamsFromError returns an Error with invalid_params code with err as a
// internalError.
func InvalidParamsFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusUnprocessableEntity, err, msg, setters...)
}

// RateLimit returns an Error with rate_limit code
func RateLimit(message string, setters ...errorParamsSetter) *Error {
	return New(StatusTooManyRequests, message, setters...)
}

// RateLimitFromError returns an Error with rate_limit code with err as a
// internalError.
func RateLimitFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusTooManyRequests, err, msg, setters...)
}

// InternalServer returns an Error with internal_server code
func InternalServer(message string, setters ...errorParamsSetter) *Error {
	return New(StatusInternalServerError, message, setters...)
}

// InternalServerFromError returns an Error with internal_server code with err as a
// internalError.
func InternalServerFromError(err error, msg string, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusInternalServerError, err, msg, setters...)
}

// MarshalJSON serialize error to json
func (e *Error) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		Meta       Meta   `json:"meta,omitempty"`
		Message    string `json:"msg,omitempty"`
		ErrorID    string `json:"error_id"`
		StatusCode Code   `json:"status_code"`
	}{e.Meta, e.Message, fmt.Sprint(e.StatusCode), e.StatusCode})
}
