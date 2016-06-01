package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Error struct {
	statusCode int
	errorID    string
	params     errorParams
}

func (e *Error) build(setters ...errorParamsSetter) {
	for _, setter := range setters {
		setter(&e.params)
	}
}

type errorParams struct {
	Meta        map[string]interface{} `json:"omitempty"`
	Description string                 `json:"omitempty"`
}

type errorParamsSetter func(*errorParams)

func SetMeta(m map[string]interface{}) errorParamsSetter {
	return func(e *errorParams) {
		e.Meta = m
	}
}

func SetDescription(d string) errorParamsSetter {
	return func(e *errorParams) {
		e.Description = d
	}
}

func NewBadRequest(setters ...errorParamsSetter) *Error {
	err := &Error{statusCode: 400, errorID: "bad_request"}

	err.build(setters...)
	return err
}

func NewInvalidParams(setters ...errorParamsSetter) *Error {
	err := &Error{statusCode: 422, errorID: "invalid_params"}

	err.build(setters...)
	return err
}

func NewMissingParams(setters ...errorParamsSetter) *Error {
	err := &Error{statusCode: 422, errorID: "missing_params"}

	err.build(setters...)
	return err
}

func NewInternalServer(setters ...errorParamsSetter) *Error {
	err := &Error{statusCode: 422, errorID: "missing_params"}

	err.build(setters...)
	return err
}

func notifyError(e Error) {
	fmt.Println("Sending to sentry...", e)
	time.Sleep(3 * time.Second)
}

func (e *Error) WriteJSON(w http.ResponseWriter) error {
	defer notifyError(*e)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(e.statusCode)

	return json.NewEncoder(w).Encode(e)
}

func (e *Error) MarshalJSON() (b []byte, err error) {
	s := struct {
		Meta        interface{} `json:"meta,omitempty"`
		Description string      `json:"description,omitempty"`
		ErrorID     string      `json:"error_id"`
	}{
		Description: e.params.Description,
		Meta:        e.params.Meta,
		ErrorID:     e.errorID,
	}

	return json.Marshal(s)
}
