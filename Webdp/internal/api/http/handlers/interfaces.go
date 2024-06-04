package handlers

import (
	"encoding/json"
	"net/http"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/response"
)

const (
	CONTENT_TYPE = "Content-Type"
	APP_JSON     = "application/json"
	TEXT_PLAIN   = "text/plain"
	TEXT_CSV     = "text/csv"
	TEXT_MD      = "text/markdown"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func HandlerDecorator(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			RenderError(w, err)
		}
	}
}

func RenderResponse(w http.ResponseWriter, response response.HttpResponse[any], headerKwargs ...string) error {
	if response.HasBody() {
		return RenderJsonResponse(w, response.GetStatusCode(), response, headerKwargs...)
	} else {
		return RenderJsonResponse(w, response.GetStatusCode(), nil, headerKwargs...)
	}

}

func RenderRawResponse(w http.ResponseWriter, httpStatusCode int, response []byte, headerKwargs ...string) error {
	setHeaderKwargs(w, headerKwargs...)
	w.WriteHeader(httpStatusCode)
	_, err := w.Write(response)
	return err
}

func RenderJsonResponse(w http.ResponseWriter, httpStatusCode int, jsonSerialisableContent any, headerKwargs ...string) error {
	setHeaderKwargs(w, headerKwargs...)
	w.Header().Add(CONTENT_TYPE, APP_JSON)
	w.WriteHeader(httpStatusCode)
	if jsonSerialisableContent == nil {
		return nil
	}
	return json.NewEncoder(w).Encode(jsonSerialisableContent)
}

func RenderCSVResponse(w http.ResponseWriter, httpStatusCode int, csv []byte, headerKwargs ...string) error {
	setHeaderKwargs(w, headerKwargs...)
	w.Header().Add(CONTENT_TYPE, TEXT_CSV)
	w.WriteHeader(httpStatusCode)
	_, err := w.Write(csv)
	return err
}

func RenderMDResponse(w http.ResponseWriter, httpStatusCode int, content []byte, headerKwargs ...string) error {
	setHeaderKwargs(w, headerKwargs...)
	w.Header().Add(CONTENT_TYPE, TEXT_MD)
	w.WriteHeader(httpStatusCode)
	_, err := w.Write(content)
	return err
}

func RenderTextResponse(w http.ResponseWriter, httpStatusCode int, text []byte, headerKwargs ...string) error {
	setHeaderKwargs(w, headerKwargs...)
	w.Header().Add(CONTENT_TYPE, TEXT_PLAIN)
	w.WriteHeader(httpStatusCode)
	_, err := w.Write(text)
	return err
}

func setHeaderKwargs(w http.ResponseWriter, headerKwargs ...string) {
	if len(headerKwargs)%2 == 0 {
		for i := 0; i < len(headerKwargs); i += 2 {
			w.Header().Add(headerKwargs[i], headerKwargs[i+1])
		}
	}
}

func RenderError(w http.ResponseWriter, err error) error {
	w.Header().Add("Content-Type", "application/json")
	fail := errors.ExpandError(err)
	w.WriteHeader(fail.GetStatusCode())
	return json.NewEncoder(w).Encode(fail)
}
