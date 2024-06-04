package routes

import (
	"webdp/internal/api/http/handlers"

	"github.com/gorilla/mux"
)

func RegisterQueriesV1(router *mux.Router, handler handlers.QueryHandler) {
	query := router.PathPrefix("/query").Subrouter()
	query.HandleFunc("/evaluate", handlers.HandlerDecorator(handler.PostQueryEvaluate)).Methods("POST")
	query.HandleFunc("/accuracy", handlers.HandlerDecorator(handler.PostQueryAccuracy)).Methods("POST")
	query.HandleFunc("/custom", handlers.HandlerDecorator(handler.PostQueryCustom)).Methods("POST")
}

func RegisterQueriesV2(router *mux.Router, handler handlers.QueryHandler) {
	queries := router.PathPrefix("/queries").Subrouter()
	queries.HandleFunc("/evaluate", handlers.HandlerDecorator(handler.PostQueryEvaluate)).Methods("POST")
	queries.HandleFunc("/accuracy", handlers.HandlerDecorator(handler.PostQueryAccuracy)).Methods("POST")
	queries.HandleFunc("/validate", handlers.HandlerDecorator(handler.PostQueryValidate)).Methods("POST")
	queries.HandleFunc("/functions", handlers.HandlerDecorator(handler.GetQueryFunctions)).Methods("GET")
	queries.HandleFunc("/docs", handlers.HandlerDecorator(handler.GetQueryDocs)).Methods("GET")
	queries.HandleFunc("/engines", handlers.HandlerDecorator(handler.GetQueryEngines)).Methods("GET")
}
