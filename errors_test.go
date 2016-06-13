package errors

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		code    int
		id      string
		setters []errorParamsSetter
		desc    string
		meta    Meta
	}{
		{0, "foo", nil, "", nil},
		{1, "bar", []errorParamsSetter{SetDescription("hi")}, "hi", nil},
		{2, "baz", []errorParamsSetter{SetDescription("hi"), SetDescription("ho")}, "ho", nil},
		{3, "bax", nil, "", nil},
		{4, "xab", []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetDescription("let's go")}, "let's go", Meta{"hi": "ho"}},
		{5, "zab", []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetDescription("let's go"), SetMeta(Meta{"ho": "hi"})}, "let's go", Meta{"ho": "hi"}},
	}

	for _, tt := range tests {
		got := New(tt.code, tt.id, tt.setters...)
		if got.StatusCode != tt.code {
			t.Errorf("New(%d, %q, %v) = %v, exptected status %d, got %d\n", tt.code, tt.id, tt.setters, got, tt.code, got.StatusCode)
		}
		if got.ErrorID != tt.id {
			t.Errorf("New(%d, %q, %v) = %v, exptected error id %q, got %q\n", tt.code, tt.id, tt.setters, got, tt.id, got.ErrorID)
		}
		if got.params.Description != tt.desc {
			t.Errorf("New(%d, %q, %v) = %v, exptected description %q, got %q\n", tt.code, tt.id, tt.setters, got, tt.desc, got.params.Description)
		}
		if !reflect.DeepEqual(tt.meta, got.params.Meta) {
			t.Errorf("New(%d, %q, %v) = %v, exptected meta %v, got %v\n", tt.code, tt.id, tt.setters, got, tt.meta, got.params.Meta)
		}
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		code    int
		id      string
		setters []errorParamsSetter
		exp     string
	}{
		{0, "foo", nil, `status_code=0 error_id="foo"`},
		{1, "bar", []errorParamsSetter{SetDescription("hi")}, `status_code=1 error_id="bar" description="hi"`},
		{2, "baz", []errorParamsSetter{SetDescription("hi"), SetDescription("ho")}, `status_code=2 error_id="baz" description="ho"`},
		{3, "bax", nil, `status_code=3 error_id="bax"`},
		{4, "xab", []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetDescription("let's go")}, `status_code=4 error_id="xab" description="let's go" hi=ho`},
		{5, "zab", []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetDescription("let's go"), SetMeta(Meta{"ho": "hi"})}, `status_code=5 error_id="zab" description="let's go" ho=hi`},
		{6, "rab", []errorParamsSetter{SetMeta(Meta{"ho": "hi"}), SetDescription("let's go")}, `status_code=6 error_id="rab" description="let's go" ho=hi`},
	}

	for _, tt := range tests {
		err := New(tt.code, tt.id, tt.setters...)
		got := err.Error()
		if got != tt.exp {
			t.Errorf("(%v).Error() = %q\n exp: %q\n got: %q\n", err, got, tt.exp, got)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		code    int
		id      string
		setters []errorParamsSetter
		exp     []byte
	}{
		{0, "foo", nil, []byte(`{"error_id":"foo","status_code":0}`)},
		{1, "bar", []errorParamsSetter{SetDescription("hi")}, []byte(`{"description":"hi","error_id":"bar","status_code":1}`)},
		{2, "baz", []errorParamsSetter{SetDescription("hi"), SetDescription("ho")}, []byte(`{"description":"ho","error_id":"baz","status_code":2}`)},
		{3, "bax", nil, []byte(`{"error_id":"bax","status_code":3}`)},
		{4, "xab", []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetDescription("let's go")}, []byte(`{"meta":{"hi":"ho"},"description":"let's go","error_id":"xab","status_code":4}`)},
		{5, "zab", []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetDescription("let's go"), SetMeta(Meta{"ho": "hi"})}, []byte(`{"meta":{"ho":"hi"},"description":"let's go","error_id":"zab","status_code":5}`)},
		{6, "rab", []errorParamsSetter{SetMeta(Meta{"hi": "ho", "ho": "hi"}), SetDescription("let's go")}, []byte(`{"meta":{"hi":"ho","ho":"hi"},"description":"let's go","error_id":"rab","status_code":6}`)},
		{StatusInternalServerError, "rab", []errorParamsSetter{SetMeta(Meta{"hi": "ho", "ho": "hi"}), SetDescription("let's go")}, []byte(`{"error_id":"rab","status_code":500}`)},
	}

	for _, tt := range tests {
		err := New(tt.code, tt.id, tt.setters...)
		got, _ := json.Marshal(err)
		if !reflect.DeepEqual(got, tt.exp) {
			t.Errorf("json.Marshal(%v) = %q\n exp: %q\n got: %q\n", err, got, tt.exp, got)
		}
	}
}
