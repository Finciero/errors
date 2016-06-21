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

var (
	ErrBadRequest     = BadRequest()
	ErrUnauthorized   = Unauthorized()
	ErrDelinquent     = Delinquent()
	ErrForbidden      = Forbidden()
	ErrSuspended      = Suspended()
	ErrNotFound       = NotFound()
	ErrNotAcceptable  = NotAcceptable()
	ErrInvalidParams  = InvalidParams()
	ErrRateLimit      = RateLimit()
	ErrInternalServer = InternalServer()
)

type Error struct {
	StatusCode  Code
	Meta        Meta
	Description string
}

type Meta map[string]interface{}

type errorParams struct {
	meta        Meta
	description string
}

func New(code Code, setters ...errorParamsSetter) *Error {
	var p errorParams
	for _, s := range setters {
		s(&p)
	}
	return &Error{StatusCode: code, Meta: p.meta, Description: p.description}
}

func NewFromError(code Code, err error, setters ...errorParamsSetter) *Error {
	var p errorParams
	for _, s := range setters {
		s(&p)
	}
	return &Error{StatusCode: code, Meta: p.meta, Description: err.Error()}
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
	} else {
		return &Error{
			StatusCode:  Code(code),
			Meta:        raw.Meta,
			Description: raw.Description,
		}
	}
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

type errorParamsSetter func(*errorParams)

func SetMeta(m Meta) errorParamsSetter {
	return func(e *errorParams) {
		if e.meta == nil {
			e.meta = m
			return
		}

		for key, value := range m {
			e.meta[key] = value
		}
	}
}

func SetDescription(d string) errorParamsSetter {
	return func(e *errorParams) {
		e.description = d
	}
}

func BadRequest(setters ...errorParamsSetter) *Error {
	return New(StatusBadRequest, setters...)
}

func Unauthorized(setters ...errorParamsSetter) *Error {
	return New(StatusUnauthorized, setters...)
}

func Delinquent(setters ...errorParamsSetter) *Error {
	return New(StatusPaymentRequired, setters...)
}

func Forbidden(setters ...errorParamsSetter) *Error {
	return New(StatusForbidden, setters...)
}

func Suspended(setters ...errorParamsSetter) *Error {
	return New(StatusForbidden, setters...)
}

func NotFound(setters ...errorParamsSetter) *Error {
	return New(StatusNotFound, setters...)
}

func NotAcceptable(setters ...errorParamsSetter) *Error {
	return New(StatusNotAcceptable, setters...)
}

func InvalidParams(setters ...errorParamsSetter) *Error {
	return New(StatusUnprocessableEntity, setters...)
}

func RateLimit(setters ...errorParamsSetter) *Error {
	return New(StatusTooManyRequests, setters...)
}

func InternalServer(setters ...errorParamsSetter) *Error {
	return New(StatusInternalServerError, setters...)
}

func (e *Error) MarshalJSON() (b []byte, err error) {
	return json.Marshal(struct {
		Meta       Meta   `json:"meta,omitempty"`
		Msg        string `json:"msg,omitempty"`
		ErrorID    string `json:"error_id"`
		StatusCode Code   `json:"status_code"`
	}{e.Meta, e.Description, fmt.Sprint(e.StatusCode), e.StatusCode})
}
