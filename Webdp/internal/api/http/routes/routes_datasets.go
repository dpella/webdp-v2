package routes

import (
	"webdp/internal/api/http/handlers"

	"github.com/gorilla/mux"
)

func RegisterDatasetsV1(router *mux.Router, handler handlers.DatasetHandler) {
	datasets := router.PathPrefix("/datasets").Subrouter()
	dataset := router.PathPrefix("/dataset").Subrouter()

	datasets.HandleFunc("", handlers.HandlerDecorator(handler.GetDatasets)).Methods("GET")
	datasets.HandleFunc("", handlers.HandlerDecorator(handler.PostDataset)).Methods("POST")
	dataset.HandleFunc("/{datasetId}", handlers.HandlerDecorator(handler.GetDataset)).Methods("GET")
	dataset.HandleFunc("/{datasetId}", handlers.HandlerDecorator(handler.PatchDataset)).Methods("PATCH")
	dataset.HandleFunc("/{datasetId}", handlers.HandlerDecorator(handler.DeleteDataset)).Methods("DELETE")

	// upload the csv to the db
	dataset.HandleFunc("/{datasetId}/upload", handlers.HandlerDecorator(handler.UploadData)).Methods("POST")
}

func RegisterDatasetsV2(router *mux.Router, handler handlers.DatasetHandler) {
	datasets := router.PathPrefix("/datasets").Subrouter()

	datasets.HandleFunc("", handlers.HandlerDecorator(handler.GetDatasets)).Methods("GET")
	datasets.HandleFunc("", handlers.HandlerDecorator(handler.PostDataset)).Methods("POST")
	datasets.HandleFunc("/{datasetId}", handlers.HandlerDecorator(handler.GetDataset)).Methods("GET")
	datasets.HandleFunc("/{datasetId}", handlers.HandlerDecorator(handler.PatchDataset)).Methods("PATCH")
	datasets.HandleFunc("/{datasetId}", handlers.HandlerDecorator(handler.DeleteDataset)).Methods("DELETE")

	// upload the csv to the db
	datasets.HandleFunc("/{datasetId}/upload", handlers.HandlerDecorator(handler.UploadData)).Methods("POST")
}

func RegisterInternalDatasets(router *mux.Router, handler handlers.InternalDatasetHandler) {
	datasets := router.PathPrefix("/datasets").Subrouter()
	datasets.HandleFunc("/{datasetId}", handlers.HandlerDecorator(handler.GetDataset)).Methods("GET")
}
