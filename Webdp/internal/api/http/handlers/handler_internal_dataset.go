package handlers

import (
	"net/http"
	"strconv"
	"webdp/internal/api/http/services"

	"github.com/gorilla/mux"
)

/*
	This struct is for handling internal requests on datasets
	Do not expose this handler outside the private network
*/

type InternalDatasetHandler struct {
	dataService services.InternalDatasetService
}

func NewInternalDatasetHandler(service services.InternalDatasetService) *InternalDatasetHandler {
	return &InternalDatasetHandler{dataService: service}
}

func (dh InternalDatasetHandler) GetDataset(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	res, err := dh.dataService.GetTable(id)

	if err != nil {
		return RenderError(w, err)
	}

	return RenderCSVResponse(w, http.StatusOK, res, "Content-Type", "text/csv")

}
