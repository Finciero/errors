package errors

import (
	"encoding/json"
	"fmt"
)

const (
	StatusContinue           = 100
	StatusSwitchingProtocols = 101

	StatusOK                   = 200
	StatusCreated              = 201
	StatusAccepted             = 202
	StatusNonAuthoritativeInfo = 203
	StatusNoContent            = 204
	StatusResetContent         = 205
	StatusPartialContent       = 206

	StatusMultipleChoices   = 300
	StatusMovedPermanently  = 301
	StatusFound             = 302
	StatusSeeOther          = 303
	StatusNotModified       = 304
	StatusUseProxy          = 305
	StatusTemporaryRedirect = 307

	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
	StatusUnprocessableEntity          = 422
	StatusPreconditionRequired         = 428
	StatusTooManyRequests              = 429
	StatusRequestHeaderFieldsTooLarge  = 431
	StatusUnavailableForLegalReasons   = 451

	StatusInternalServerError           = 500
	StatusNotImplemented                = 501
	StatusBadGateway                    = 502
	StatusServiceUnavailable            = 503
	StatusGatewayTimeout                = 504
	StatusHTTPVersionNotSupported       = 505
	StatusNetworkAuthenticationRequired = 511
)

const (
	MsgBadRequest          = "bad request"
	MsgInvalidParams       = "invalid parameters"
	MsgMissingParams       = "missing parameters"
	MsgIntervalServerError = "internal server error"
	MsgNotFound            = "not found"
)

var (
	ErrNotFound            = NotFound()
	ErrMissingParams       = MissingParams()
	ErrInvalidParams       = InvalidParams()
	ErrBadRequest          = BadRequest()
	ErrIntervalServerError = InternalServerError()
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
	return New(StatusBadRequest, MsgBadRequest, setters...)
}

func InvalidParams(setters ...errorParamsSetter) *Error {
	return New(StatusUnprocessableEntity, MsgInvalidParams, setters...)
}

func MissingParams(setters ...errorParamsSetter) *Error {
	return New(StatusUnprocessableEntity, MsgMissingParams, setters...)
}

func InternalServerError(setters ...errorParamsSetter) *Error {
	return New(StatusInternalServerError, MsgIntervalServerError, setters...)
}

func NotFound(setters ...errorParamsSetter) *Error {
	return New(StatusNotFound, MsgNotFound, setters...)
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
