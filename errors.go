// go:generate stringer -type=Code
package errors

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Code int32

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

type Error struct {
	StatusCode  Code
	Meta        Meta
	Description string

	internal string // internal information used for debugging
}

type Meta map[string]interface{}

func New(code Code, message string, setters ...errorParamsSetter) *Error {
	var m Meta
	for _, s := range setters {
		s(&m)
	}
	return &Error{StatusCode: code, Meta: m, Description: message}
}

func NewFromError(code Code, err error, setters ...errorParamsSetter) *Error {
	var m Meta
	for _, s := range setters {
		s(&m)
	}
	return &Error{StatusCode: code, Meta: m, Description: err.Error()}
}

func FromGRPC(err error) *Error {
	var raw = struct {
		Meta        Meta   `json:"meta,omitempty"`
		Description string `json:"msg,omitempty"`
	}{}

	code := grpc.Code(err)
	desc := grpc.ErrorDesc(err)
	if err = json.Unmarshal([]byte(desc), &raw); err != nil {
		return &Error{
			StatusCode:  Code(code),
			Description: desc,
		}
	}

	return &Error{
		StatusCode:  Code(code),
		Meta:        raw.Meta,
		Description: raw.Description,
	}
}

func (e *Error) Code() int {
	return int(e.StatusCode)
}

func (e *Error) Error() string {
	str := fmt.Sprintf("status_code=%d error_id=%q", e.StatusCode, fmt.Sprint(e.StatusCode))

	if len(e.Description) > 0 {
		str += fmt.Sprintf(" msg=%q", e.Description)
	}

	for key, value := range e.Meta {
		str += fmt.Sprintf(" %s=%s", key, value)
	}

	return str
}

func (e *Error) ErrorID() string {
	return fmt.Sprint(e.StatusCode)
}

func (e *Error) ToGRPC() error {
	buff, _ := json.Marshal(struct {
		Meta        Meta   `json:"meta,omitempty"`
		Description string `json:"msg,omitempty"`
	}{e.Meta, e.Description})
	return grpc.Errorf(codes.Code(e.StatusCode), string(buff))
}

type errorParamsSetter func(*Meta)

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

func BadRequest(message string, setters ...errorParamsSetter) *Error {
	return New(StatusBadRequest, message, setters...)
}

func BadRequestFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusBadRequest, err, setters...)
}

func Unauthorized(message string, setters ...errorParamsSetter) *Error {
	return New(StatusUnauthorized, message, setters...)
}

func UnauthorizedFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusUnauthorized, err, setters...)
}

func Delinquent(message string, setters ...errorParamsSetter) *Error {
	return New(StatusPaymentRequired, message, setters...)
}

func DelinquentFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusPaymentRequired, err, setters...)
}

func Forbidden(message string, setters ...errorParamsSetter) *Error {
	return New(StatusForbidden, message, setters...)
}

func ForbiddenFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusForbidden, err, setters...)
}

func Suspended(message string, setters ...errorParamsSetter) *Error {
	return New(StatusForbidden, message, setters...)
}

func SuspendedFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusForbidden, err, setters...)
}

func NotFound(message string, setters ...errorParamsSetter) *Error {
	return New(StatusNotFound, message, setters...)
}

func NotFoundFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusNotFound, err, setters...)
}

func NotAcceptable(message string, setters ...errorParamsSetter) *Error {
	return New(StatusNotAcceptable, message, setters...)
}

func NotAcceptableFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusNotAcceptable, err, setters...)
}

func InvalidParams(message string, setters ...errorParamsSetter) *Error {
	return New(StatusUnprocessableEntity, message, setters...)
}

func InvalidParamsFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusUnprocessableEntity, err, setters...)
}

func RateLimit(message string, setters ...errorParamsSetter) *Error {
	return New(StatusTooManyRequests, message, setters...)
}

func RateLimitFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusTooManyRequests, err, setters...)
}

func InternalServer(message string, setters ...errorParamsSetter) *Error {
	return New(StatusInternalServerError, message, setters...)
}

func InternalServerFromError(err error, setters ...errorParamsSetter) *Error {
	return NewFromError(StatusInternalServerError, err, setters...)
}

func (e *Error) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		Meta       Meta   `json:"meta,omitempty"`
		Msg        string `json:"msg,omitempty"`
		ErrorID    string `json:"error_id"`
		StatusCode Code   `json:"status_code"`
	}{e.Meta, e.Description, fmt.Sprint(e.StatusCode), e.StatusCode})
}
