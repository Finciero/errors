package errors

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestNew(t *testing.T) {
	tests := []struct {
		code    Code
		id      string
		setters []errorParamsSetter
		desc    string
		meta    Meta
	}{
		{
			code:    0,
			id:      "Code(0)",
			setters: nil,
			desc:    "",
			meta:    nil,
		},
		{
			code:    1,
			id:      "Code(1)",
			setters: nil,
			desc:    "hi",
			meta:    nil,
		},
		{
			code:    4,
			id:      "Code(4)",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"})},
			desc:    "let's go",
			meta:    Meta{"hi": "ho"},
		},
		{
			code:    5,
			id:      "Code(5)",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetMeta(Meta{"ho": "hi"})},
			desc:    "let's go",
			meta:    Meta{"ho": "hi", "hi": "ho"},
		},
	}

	for _, tt := range tests {
		got := New(tt.code, tt.desc, tt.setters...)
		if got.StatusCode != tt.code {
			t.Errorf("New(%d, %v) = %v, unexpected status\n exp: %d\n got: %d\n", tt.code, tt.setters, got, tt.code, got.StatusCode)
		}
		if got.ErrorID() != tt.id {
			t.Errorf("New(%d, %v) = %v, unexpected error_id\n exp: %q\n got:  %q\n", tt.code, tt.setters, got, tt.id, got.ErrorID())
		}
		if got.Description != tt.desc {
			t.Errorf("New(%d, %v) = %v, unexpected description\n exp: %q\n got: %q\n", tt.code, tt.setters, got, tt.desc, got.Description)
		}
		if !reflect.DeepEqual(tt.meta, got.Meta) {
			t.Errorf("New(%d, %v) = %v, unexpected meta\n exp: %v\n got: %v\n", tt.code, tt.setters, got, tt.meta, got.Meta)
		}
	}
}

func TestNewFromError(t *testing.T) {
	var (
		errMsg  = "test: new error"
		errTest = errors.New(errMsg)
	)

	tests := []struct {
		code    Code
		id      string
		setters []errorParamsSetter
		meta    Meta
	}{
		{
			code:    1,
			id:      "Code(1)",
			setters: nil,
			meta:    nil,
		},
		{
			code:    4,
			id:      "Code(4)",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"})},
			meta:    Meta{"hi": "ho"},
		},
		{
			code:    5,
			id:      "Code(5)",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetMeta(Meta{"ho": "hi"})},
			meta:    Meta{"ho": "hi", "hi": "ho"},
		},
	}

	for _, tt := range tests {
		got := NewFromError(tt.code, errTest, tt.setters...)
		if got.StatusCode != tt.code {
			t.Errorf("New(%d, %v) = %v, unexpected status\n exp: %d\n got: %d\n", tt.code, tt.setters, got, tt.code, got.StatusCode)
		}
		if got.ErrorID() != tt.id {
			t.Errorf("New(%d, %v) = %v, unexpected error_id\n exp: %q\n got:  %q\n", tt.code, tt.setters, got, tt.id, got.ErrorID())
		}
		if got.Description != errMsg {
			t.Errorf("New(%d, %v) = %v, unexpected description\n exp: %q\n got: %q\n", tt.code, tt.setters, got, errMsg, got.Description)
		}
		if !reflect.DeepEqual(tt.meta, got.Meta) {
			t.Errorf("New(%d, %v) = %v, unexpected meta\n exp: %v\n got: %v\n", tt.code, tt.setters, got, tt.meta, got.Meta)
		}
	}
}

func TestFromGRPC(t *testing.T) {
	tests := []struct {
		code int
		desc string
		exp  *Error
	}{
		{
			code: int(StatusBadRequest),
			desc: `{"meta":{"hi":"ho"},"msg":"let's go"}`,
			exp:  New(StatusBadRequest, "let's go", SetMeta(Meta{"hi": "ho"})),
		},
		{
			code: int(StatusBadRequest),
			desc: `{"meta":{"hi":"ho"},"msg":"let's go"}`,
			exp:  BadRequest("let's go", SetMeta(Meta{"hi": "ho"})),
		},
		{
			code: int(StatusUnauthorized),
			desc: `{"msg":"let's go"}`,
			exp:  New(StatusUnauthorized, "let's go"),
		},
		{
			code: int(StatusUnauthorized),
			desc: `{"msg":"let's go"}`,
			exp:  Unauthorized("let's go"),
		},
		{
			code: int(StatusUnauthorized),
			desc: "let's go",
			exp:  Unauthorized("let's go"),
		},
	}

	for _, tt := range tests {
		in := grpc.Errorf(codes.Code(tt.code), tt.desc)
		err := FromGRPC(in)

		if !reflect.DeepEqual(err, tt.exp) {
			t.Errorf("FromGRPC(%v) = %v\n exp: %v\n got: %v\n", in, err, tt.exp, err)
		}
	}
}

func TestToGRPCFromGRPC(t *testing.T) {

	tests := []struct {
		err *Error
	}{
		{New(StatusBadRequest, "let's go", SetMeta(Meta{"hi": "ho"}))},
		{BadRequest("let's go", SetMeta(Meta{"hi": "ho"}))},
		{New(StatusUnauthorized, "let's go")},
		{Unauthorized("let's go")},
	}

	for _, tt := range tests {
		in := tt.err.ToGRPC()
		err := FromGRPC(in)

		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("FromGRPC(%v) = %v\n exp: %v\n got: %v\n", in, err, tt.err, err)
		}
	}
}

func TestToGRPC(t *testing.T) {
	tests := []struct {
		err *Error
		exp string
	}{
		{Unauthorized("", SetMeta(Meta{"hi": "ho"})), `{"meta":{"hi":"ho"}}`},
		{InternalServer(""), `{}`},
		{BadRequest(""), `{}`},
		{Forbidden(""), `{}`},
		{InvalidParams(""), `{}`},
		{NotAcceptable(""), `{}`},
		{NotFound(""), `{}`},
		{Delinquent(""), `{}`},
		{RateLimit(""), `{}`},
		{Unauthorized(""), `{}`},
		{Unauthorized("some error", SetMeta(Meta{"hi": "ho"})), `{"meta":{"hi":"ho"},"msg":"some error"}`},
		{RateLimit("some error", SetMeta(Meta{"hi": "ho"}), SetMeta(Meta{"hi": "hi"})), `{"meta":{"hi":"hi"},"msg":"some error"}`},
	}

	for _, tt := range tests {
		got := tt.err.ToGRPC() // grpc error
		if (int32)(grpc.Code(got)) != (int32)(tt.err.StatusCode) || grpc.ErrorDesc(got) != tt.exp {
			t.Errorf("(%v).ToGRPC()\n got: {code: %d, desc: %q}\n exp: {code: %d, desc: %q}\n",
				tt.err, grpc.Code(got), string(grpc.ErrorDesc(got)), tt.err.StatusCode, tt.exp)
		}
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		code    Code
		msg     string
		setters []errorParamsSetter
		exp     string
	}{
		{
			code:    0,
			msg:     "",
			setters: nil,
			exp:     `status_code=0 error_id="Code(0)"`,
		},
		{
			code:    1,
			msg:     "hi",
			setters: nil,
			exp:     `status_code=1 error_id="Code(1)" msg="hi"`,
		},
		{
			code:    2,
			msg:     "ho",
			setters: nil,
			exp:     `status_code=2 error_id="Code(2)" msg="ho"`,
		},
		{
			code:    3,
			msg:     "",
			setters: nil,
			exp:     `status_code=3 error_id="Code(3)"`,
		},
		{
			code:    4,
			msg:     "let's go",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"})},
			exp:     `status_code=4 error_id="Code(4)" msg="let's go" hi=ho`,
		},
		{
			code:    5,
			msg:     "let's go",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetMeta(Meta{"hi": "hi"})},
			exp:     `status_code=5 error_id="Code(5)" msg="let's go" hi=hi`,
		},
		{
			code:    6,
			msg:     "let's go",
			setters: []errorParamsSetter{SetMeta(Meta{"ho": "hi"})},
			exp:     `status_code=6 error_id="Code(6)" msg="let's go" ho=hi`,
		},
	}

	for _, tt := range tests {
		err := New(tt.code, tt.msg, tt.setters...)
		got := err.Error()
		if got != tt.exp {
			t.Errorf("(%v).Error() = %q\n exp: %q\n got: %q\n\n", err, got, tt.exp, got)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		code    Code
		msg     string
		setters []errorParamsSetter
		exp     []byte
	}{
		{
			code:    0,
			setters: nil,
			exp:     []byte(`{"error_id":"Code(0)","status_code":0}`),
		},
		{
			code:    1,
			msg:     "hi",
			setters: nil,
			exp:     []byte(`{"msg":"hi","error_id":"Code(1)","status_code":1}`),
		},
		{
			code:    4,
			msg:     "let's go",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"})},
			exp:     []byte(`{"meta":{"hi":"ho"},"msg":"let's go","error_id":"Code(4)","status_code":4}`),
		},
		{
			code:    5,
			msg:     "let's go",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho"}), SetMeta(Meta{"ho": "hi"})},
			exp:     []byte(`{"meta":{"hi":"ho","ho":"hi"},"msg":"let's go","error_id":"Code(5)","status_code":5}`),
		},
		{
			code:    6,
			msg:     "let's go",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho", "ho": "hi"})},
			exp:     []byte(`{"meta":{"hi":"ho","ho":"hi"},"msg":"let's go","error_id":"Code(6)","status_code":6}`),
		},
		{
			code:    StatusInternalServerError,
			msg:     "let's go",
			setters: []errorParamsSetter{SetMeta(Meta{"hi": "ho", "ho": "hi"})},
			exp:     []byte(`{"meta":{"hi":"ho","ho":"hi"},"msg":"let's go","error_id":"internal_server","status_code":500}`),
		},
		{
			code:    StatusBadRequest,
			setters: nil,
			exp:     []byte(`{"error_id":"bad_request","status_code":400}`),
		},
		{
			code:    StatusForbidden,
			setters: nil,
			exp:     []byte(`{"error_id":"forbidden","status_code":403}`),
		},
		{
			code:    StatusUnprocessableEntity,
			setters: nil,
			exp:     []byte(`{"error_id":"invalid_params","status_code":422}`),
		},
		{
			code:    StatusNotAcceptable,
			setters: nil,
			exp:     []byte(`{"error_id":"not_acceptable","status_code":406}`),
		},
		{
			code:    StatusNotFound,
			setters: nil,
			exp:     []byte(`{"error_id":"not_found","status_code":404}`),
		},
		{
			code:    StatusPaymentRequired,
			setters: nil,
			exp:     []byte(`{"error_id":"delinquent","status_code":402}`),
		},
		{
			code:    StatusTooManyRequests,
			setters: nil,
			exp:     []byte(`{"error_id":"rate_limit","status_code":429}`),
		},
		{
			code:    StatusUnauthorized,
			setters: nil,
			exp:     []byte(`{"error_id":"unauthorized","status_code":401}`),
		},
	}

	for _, tt := range tests {
		err := New(tt.code, tt.msg, tt.setters...)
		got, _ := json.Marshal(err)
		if !reflect.DeepEqual(got, tt.exp) {
			t.Errorf("json.Marshal(%v) = %q\n exp: %q\n got: %q\n", err, got, tt.exp, got)
		}
	}
}
