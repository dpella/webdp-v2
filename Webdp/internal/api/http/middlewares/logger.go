package middlewares

import (
	"fmt"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := fmt.Sprintf("METHOD: %s | REQUEST_URI: %s", r.Method, r.URL.RequestURI())
		fmt.Println(s)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
