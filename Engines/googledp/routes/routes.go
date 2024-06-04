package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"googledp/client"
	"googledp/dpfuncs"
	"googledp/requests"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type handlers struct {
	Cache map[int64][][]string
}

func RegisterRoutes(r *mux.Router) {

	handlers := &handlers{
		Cache: make(map[int64][][]string),
	}

	r.HandleFunc("/evaluate", wrapperHttp(handlers.evaluate)).Methods("POST")
	r.HandleFunc("/validate", wrapperHttp(handlers.validate)).Methods("POST")
	r.HandleFunc("/accuracy", wrapperHttp(handlers.accuracy)).Methods("POST")
	r.HandleFunc("/documentation", wrapperHttp(handlers.documentation)).Methods("GET")
	r.HandleFunc("/functions", wrapperHttp(handlers.functions)).Methods("GET")
	r.HandleFunc("/cache/{datasetId}", wrapperHttp(handlers.cache)).Methods("DELETE")
}

func (h handlers) evaluate(w http.ResponseWriter, r *http.Request) error {
	var eval requests.Evaluate
	if err := parseJsonRequestBody(r, &eval); err != nil {
		fmt.Printf("%s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}

	records, ok := h.Cache[eval.Dataset]

	if !ok {
		data, err := client.GetCSVData(eval.CallbackUrl)
		if err != nil {
			fmt.Println(err.Error())
		}
		h.Cache[eval.Dataset] = data
		records = data
	}

	result, err := dpfuncs.NewEvalQuery(eval, records[1:])

	if err != nil {
		fmt.Printf("%s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, err.Error())
	}

	return WriteJSON(w, http.StatusOK, &result)
}

type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	Status string `json:"status"`
}

func (h handlers) validate(w http.ResponseWriter, r *http.Request) error {
	var eval requests.Evaluate
	if err := parseJsonRequestBody(r, &eval); err != nil {
		fmt.Printf("%s", err.Error())
		return WriteJSON(w, http.StatusOK, ValidateResponse{Valid: false, Status: err.Error()})
	}

	records, ok := h.Cache[eval.Dataset]

	if !ok {
		data, err := client.GetCSVData(eval.CallbackUrl)
		if err != nil {
			fmt.Println(err.Error())
			return WriteJSON(w, http.StatusInternalServerError, ValidateResponse{Valid: false, Status: err.Error()})
		}
		h.Cache[eval.Dataset] = data
		records = data
	}

	_, err := dpfuncs.NewEvalQuery(eval, records[1:])

	if err != nil {
		fmt.Printf("%s", err.Error())
		return WriteJSON(w, http.StatusOK, ValidateResponse{Valid: false, Status: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, ValidateResponse{Valid: true, Status: "query is valid in GoogleDP"})
}

func (h handlers) accuracy(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusNotImplemented, "not supported in google dp engine")
}

func (h handlers) documentation(w http.ResponseWriter, r *http.Request) error {
	file, err := os.Open("./static/README.md")

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err.Error())
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/markdown")

	_, err = io.Copy(w, file)

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (h handlers) functions(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, "./static/functions.json")
	return nil
}

func (h handlers) cache(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusNotImplemented, "not implemented")
}

// HELPER FUNCTIONS

type apiFunc func(http.ResponseWriter, *http.Request) error

func wrapperHttp(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func RenderRawResponse(w http.ResponseWriter, httpStatusCode int, response []byte) error {
	w.WriteHeader(httpStatusCode)
	_, err := w.Write(response)
	return err
}

func parseJsonRequestBody[T any](r *http.Request, obj *T) error {
	var bodyCopy []byte
	var err error
	if r.Body != nil {
		bodyCopy, err = io.ReadAll(r.Body)
		if err != nil {
			return err
		}
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
	// now ok to use copy
	return json.Unmarshal(bodyCopy, obj)
}
