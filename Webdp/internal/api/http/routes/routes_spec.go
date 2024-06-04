package routes

import (
	"net/http"
	_ "webdp/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// GetDocs godoc
// @Summary      Get OpenAPI Specification
// @Description  Returns a html of the API specification, using a Swagger UI.
// @Description  See http://localhost:8000/v2/spec/index.html
// @Tags         spec
// @Success      200  {object}  map[string]interface{}
// @Router       /v2/spec [get]
func RegisterSpec(router *mux.Router) {
	router.PathPrefix("/spec").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/v2/spec/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
}
