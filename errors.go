package errors

import (
	"encoding/json"
	"fmt"
)

const (
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusPaymentRequired     = 402
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusNotAcceptable       = 406
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429

	StatusInternalServerError = 500
)

const (
	IDBadRequest     = "bad_request"
	IDUnauthorized   = "unauthorized"
	IDDelinquent     = "delinquent"
	IDForbidden      = "forbidden"
	IDSuspended      = "suspended"
	IDNotFound       = "not_found"
	IDNotAcceptable  = "not_acceptable"
	IDInvalidParams  = "invalid_params"
	IDRateLimit      = "rate_limit"
	IDInternalServer = "internal_server"
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
	StatusCode int
	ErrorID    string
	params     errorParams
}

type Meta map[string]interface{}

type errorParams struct {
	Meta        Meta
	Description string
}

func (params errorParams) String() string {
	var desc string
	if params.Description != "" {
		desc += fmt.Sprintf("description=%q", params.Description)
	}
	for key, value := range params.Meta {
		desc += fmt.Sprintf(" %s=%s", key, value)
	}
	return desc
}

func New(sc int, id string, setters ...errorParamsSetter) *Error {
	var p errorParams
	for _, s := range setters {
		s(&p)
	}
	return &Error{StatusCode: sc, ErrorID: id, params: p}
}

func (e *Error) Error() string {
	params := e.params.String()

	if len(params) > 0 {
		return fmt.Sprintf("status_code=%d error_id=%q %s", e.StatusCode, e.ErrorID, params)
	}

	return fmt.Sprintf("status_code=%d error_id=%q", e.StatusCode, e.ErrorID)
}

type errorParamsSetter func(*errorParams)

func SetMeta(m Meta) errorParamsSetter {
	return func(e *errorParams) {
		e.Meta = m
	}
}

func SetDescription(d string) errorParamsSetter {
	return func(e *errorParams) {
		e.Description = d
	}
}

func BadRequest(setters ...errorParamsSetter) *Error {
	return New(StatusBadRequest, IDBadRequest, setters...)
}

func Unauthorized(setters ...errorParamsSetter) *Error {
	return New(StatusUnauthorized, IDUnauthorized, setters...)
}

func Delinquent(setters ...errorParamsSetter) *Error {
	return New(StatusPaymentRequired, IDDelinquent, setters...)
}

func Forbidden(setters ...errorParamsSetter) *Error {
	return New(StatusForbidden, IDForbidden, setters...)
}

func Suspended(setters ...errorParamsSetter) *Error {
	return New(StatusForbidden, IDSuspended, setters...)
}

func NotFound(setters ...errorParamsSetter) *Error {
	return New(StatusNotFound, IDNotFound, setters...)
}

func NotAcceptable(setters ...errorParamsSetter) *Error {
	return New(StatusNotAcceptable, IDNotAcceptable, setters...)
}

func InvalidParams(setters ...errorParamsSetter) *Error {
	return New(StatusUnprocessableEntity, IDInvalidParams, setters...)
}

func RateLimit(setters ...errorParamsSetter) *Error {
	return New(StatusTooManyRequests, IDRateLimit, setters...)
}

func InternalServer(setters ...errorParamsSetter) *Error {
	return New(StatusInternalServerError, IDInternalServer, setters...)
}

func (e *Error) MarshalJSON() (b []byte, err error) {
	switch e.StatusCode {
	case StatusInternalServerError:
		return json.Marshal(struct {
			ErrorID    string `json:"error_id"`
			StatusCode int    `json:"status_code"`
		}{e.ErrorID, e.StatusCode})
	default:
		return json.Marshal(struct {
			Meta        Meta   `json:"meta,omitempty"`
			Description string `json:"description,omitempty"`
			ErrorID     string `json:"error_id"`
			StatusCode  int    `json:"status_code"`
		}{e.params.Meta, e.params.Description, e.ErrorID, e.StatusCode})
	}
}
