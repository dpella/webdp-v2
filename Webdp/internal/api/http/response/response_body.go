package response

import (
	"encoding/json"
	"net/http"
)

type HttpResponse[T any] interface {
	response()
	MarshalJSON() ([]byte, error)
	GetStatusCode() int
	HasBody() bool
}

type Void interface {
	void()
}

type empty struct {
}

func (e empty) void() {}

// ===== NO CONTENT ============================

// syntactic sugar for a success response with code 204

type noContent struct{}

func NoContent() HttpResponse[Void] { return noContent{} }

func (noContent) GetStatusCode() int { return 204 }

func (noContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(empty{})
}

func (noContent) HasBody() bool { return false }

func (noContent) response() {}

// ====== empty resonse =======================
type emptyResponse struct {
	code int
}

func EmptyResponse(statusCode int) HttpResponse[Void] {
	return &emptyResponse{code: statusCode}
}

func (e emptyResponse) GetStatusCode() int { return e.code }

func (emptyResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(nil)
}
func (e *emptyResponse) response() {}

func (emptyResponse) HasBody() bool { return false }

// ============= SUCCESS ======================

type success[T any] struct {
	body T
	code int
}

// GetStatusCode implements HttpResponse.
func (s success[T]) GetStatusCode() int {
	return s.code
}

// MarshalJSON implements HttpResponse.
func (s success[T]) MarshalJSON() ([]byte, error) {
	if s.code == http.StatusNoContent {
		return json.Marshal(empty{})
	}
	return json.Marshal(s.body)
}

// response implements HttpResponse.
func (success[T]) response() {}

func (s *success[T]) WithBody(body T) *success[T] {
	s.body = body
	return s
}

func NewSuccess[T any](statusCode int, responseBody T) *success[T] {
	return &success[T]{code: statusCode, body: responseBody}
}

func (s success[T]) HasBody() bool {
	return true
}

// just type checking
var _ HttpResponse[string] = success[string]{}
var _ HttpResponse[string] = fail[string]{}
var _ HttpResponse[Void] = noContent{}

// ==== FAIL =============

type fail[T any] struct {
	title  string
	detail string
	status int
	ftype  string
}

// GetStatusCode implements HttpResponse.
func (f fail[T]) GetStatusCode() int {
	return f.status
}

// MarshalJSON implements HttpResponse.
func (f fail[T]) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["title"] = f.title
	m["detail"] = f.detail
	m["status"] = f.status
	m["type"] = f.ftype
	return json.Marshal(m)
}

func (f fail[T]) response() {}

func NewFail[T any](statusCode int) *fail[T] {
	return &fail[T]{status: statusCode}
}

func (f *fail[T]) WithTitle(title string) *fail[T] {
	f.title = title
	return f
}

func (f *fail[T]) WithDetail(detail string) *fail[T] {
	f.detail = detail
	return f
}

func (f *fail[T]) WithType(failType string) *fail[T] {
	f.ftype = failType
	return f
}

func (f fail[T]) HasBody() bool {
	return true
}
