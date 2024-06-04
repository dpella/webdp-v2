package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	errors "webdp/internal/api/http"
)

// assumes that a request is a json
// copies the request body and unmarshals the json contents
func ParseJsonRequestBody[T any](r *http.Request, obj *T) error {
	var bodyCopy []byte
	var err error
	if r.Body != nil {
		bodyCopy, err = io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("%w: parse: %s", errors.ErrBadFormatting, err.Error())
		}
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
	// now ok to use copy
	err = json.Unmarshal(bodyCopy, obj)
	if err != nil {
		return fmt.Errorf("%w: unmarshal: %s", errors.ErrBadFormatting, err.Error())
	}
	return err
}
