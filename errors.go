package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	statusCode int
	errorID    string
	params     errorParams
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"apierror: status_code: %d error_id: %s params: %+v\n",
		e.statusCode,
		e.errorID,
		e.params,
	)
}

var (
	NotFound = &Error{statusCode: 404, errorID: "not found"}
)

func (e *Error) build(setters ...errorParamsSetter) {
	for _, setter := range setters {
		setter(&e.params)
	}
}

type Meta map[string]interface{}

type errorParams struct {
	Meta        Meta
	Description string
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
	err := &Error{statusCode: 500, errorID: "missing_params"}

	err.build(setters...)
	return err
}

func compareMeta(m1, m2 Meta) bool {
	if m1 == nil && m2 == nil {
		return true
	}

	if (m1 == nil && m2 != nil) || (m1 != nil && m2 == nil) {
		return false
	}

	if len(m1) != len(m2) {
		return false
	}

	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; !ok || v1 != v2 {
			return false
		}
	}

	for k1, v1 := range m2 {
		if v2, ok := m1[k1]; !ok || v1 != v2 {
			return false
		}
	}

	return true
}

func Compare(e1, e2 *Error) bool {
	return (e1.statusCode == e2.statusCode &&
		e1.errorID == e2.errorID &&
		e1.params.Description == e2.params.Description &&
		compareMeta(e1.params.Meta, e2.params.Meta))
}

func notifyError(e Error) {
	fmt.Println("Send to sentry...", e)
}

func (e *Error) WriteJSON(w http.ResponseWriter) error {
	go notifyError(*e)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(e.statusCode)

	return json.NewEncoder(w).Encode(e)
}

func (e *Error) MarshalJSON() (b []byte, err error) {
	var description string
	var meta Meta

	if e.statusCode != 500 {
		description = e.params.Description
		meta = e.params.Meta
	}

	s := struct {
		Meta        Meta   `json:"meta,omitempty"`
		Description string `json:"description,omitempty"`
		ErrorID     string `json:"error_id"`
	}{
		Description: description,
		Meta:        meta,
		ErrorID:     e.errorID,
	}

	return json.Marshal(s)
}
