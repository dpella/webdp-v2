package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"webdp/internal/api/http/client"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/middlewares"
	"webdp/internal/api/http/response"
	"webdp/internal/api/http/services"
	"webdp/internal/api/http/utils"

	"github.com/gorilla/mux"
)

type DatasetHandler struct {
	datasetService services.DatasetService
	userService    services.UserService
	budgetService  services.BudgetService
	dpClient       client.DPClient
}

type DataInfoList = []entity.DatasetInfo

func NewDatasetHandler(ds services.DatasetService, us services.UserService, bs services.BudgetService, cli client.DPClient) DatasetHandler {
	return DatasetHandler{datasetService: ds, userService: us, budgetService: bs, dpClient: cli}
}

/*
Gets all datasets which requester has access to.
Requester needs role admin or curator, or needs granted access via budget allocation.
* TODO: get only those to which the requester has granted access!
*/
// GetDatasets godoc
// @Summary      Gets all datasets which requester has access to.
// @Description  Requester needs role admin or curator, or needs granted access via budget allocation.
// @Tags         datasets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  []entity.DatasetInfo
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/datasets [get]
// @Router       /v2/datasets [get]
func (h DatasetHandler) GetDatasets(w http.ResponseWriter, r *http.Request) error {
	datasets, err := h.datasetService.GetAllDatasets()
	if err != nil {
		return RenderError(w, err)
	}

	if err := middlewares.ValidateRoles(r, []string{entity.ADMIN, entity.CURATOR}); err != nil {
		err := middlewares.FilterGrantedAccess(r, &h.budgetService, &datasets)
		if err != nil || len(datasets) == 0 {
			return RenderError(w, err)
		}
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, datasets))
}

/*
Gets a dataset.
Requester needs role admin or curator, or needs granted access via budget allocation.
*/
// GetDataset godoc
// @Summary      Gets all datasets which requester has access to.
// @Description  Requester needs role admin or curator, or needs granted access via budget allocation.
// @Tags         datasets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        datasetId   	path   int 			  true  "Dataset Id"
// @Success      200  {object}  entity.DatasetInfo
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/dataset/{datasetId} [get]
// @Router       /v2/datasets/{datasetId} [get]
func (h DatasetHandler) GetDataset(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.ADMIN, entity.CURATOR}); err != nil {
		if err := middlewares.ValidateGrantedAccess(r, &h.budgetService); err != nil {
			return RenderError(w, err)
		}
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}
	res, err := h.datasetService.GetDataset(id)
	if err != nil {
		return RenderError(w, err)
	}
	return RenderResponse(w, response.NewSuccess(http.StatusOK, res))
}

/*
Creates a dataset.
Requester needs role curator. New owner of dataset needs role curator.
*/
// PostDataset godoc
// @Summary      Creates a dataset.
// @Description  Requester needs role curator. New owner of dataset needs role curator.
// @Tags         datasets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 requestBody	body   entity.DatasetCreate  true  "request body"
// @Success      201  {object}  response.Id
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/datasets [post]
// @Router       /v2/datasets [post]
func (h DatasetHandler) PostDataset(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR}); err != nil {
		return RenderError(w, err)
	}
	if err := middlewares.ValidateNewOwnerRole(r, &h.userService, entity.CURATOR); err != nil {
		return RenderError(w, err)
	}

	var createDataset entity.DatasetCreate
	if err := utils.ParseJsonRequestBody[entity.DatasetCreate](r, &createDataset); err != nil {
		return RenderError(w, err)
	}

	res, err := h.datasetService.CreateDataset(createDataset)
	if err != nil {
		return RenderError(w, err)
	}

	resp := response.Id{Id: res}
	return RenderResponse(w, response.NewSuccess(http.StatusCreated, resp))
}

/*
Update a dataset.
Requester needs to be curator and owner of the dataset. New owner can have any role.
  - Allowing analysts to own datasets will disallow further patches and gets; only
  	dataset deletion, data uploads and budget allocation on the dataset is allowed
	from that point on.
  - Allowing admins to own datasets will disallow anything but deletions and data
    uploads. An admin can patch a dataset budget iff a curator or analyst created
    a budget for that user before handing over the ownership of the dataset.
*/
// PatchDataset godoc
// @Summary      Update a dataset.
// @Description  Update name, owner or total budget of a dataset.
// @Description  Requester needs to be curator and owner of the dataset.
// @Tags         datasets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        datasetId   	path   int 			  true  "Dataset Id"
// @Param		 requestBody	body   entity.DatasetCreate  true  "request body"
// @Success      204
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/dataset/{datasetId} [patch]
// @Router       /v2/datasets/{datasetId} [patch]
func (h DatasetHandler) PatchDataset(w http.ResponseWriter, r *http.Request) error {
	// Checks are a bit different in the proof-of-concept.
	// Compare: "is curator, is owner" vs "is owner, owner is curator"
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR}); err != nil {
		return RenderError(w, err)
	}
	if err := middlewares.ValidateOwnership(r, &h.datasetService); err != nil {
		return RenderError(w, err)
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	var patch entity.DatasetPatch
	if err := utils.ParseJsonRequestBody[entity.DatasetPatch](r, &patch); err != nil {
		return RenderError(w, err)
	}

	if err := h.datasetService.UpdateDataset(id, patch); err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NoContent())
}

/*
Delete a dataset. Requester needs to be the owner of the dataset.
*/
// DeleteDataset godoc
// @Summary      Delete a dataset.
// @Description  Requester needs to be the owner of the dataset.
// @Tags         datasets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        datasetId   	path   int 			  true  "Dataset Id"
// @Success      204
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/dataset/{datasetId} [delete]
// @Router       /v2/datasets/{datasetId} [delete]
func (h DatasetHandler) DeleteDataset(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateOwnership(r, &h.datasetService); err != nil {
		return RenderError(w, err)
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	if err := h.datasetService.DeleteDataset(id); err != nil {
		return RenderError(w, err)
	}

	h.dpClient.RemoveDatasetFromEngineCache(id)

	return RenderResponse(w, response.NoContent())
}

/*
Upload a dataset. Requester needs to be the owner of the dataset.
*/
// UploadDataset godoc
// @Summary      Upload a dataset.
// @Description  Requester needs to be the owner of the dataset.
// @Tags         datasets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 csvData 	body	string	true  "CSV Data"
// @Param        datasetId  path    int 	true  "Dataset Id"
// @Success      204
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/dataset/{datasetId}/upload [post]
// @Router       /v2/datasets/{datasetId}/upload [post]
func (h DatasetHandler) UploadData(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateOwnership(r, &h.datasetService); err != nil {
		return RenderError(w, err)
	}

	vars := mux.Vars(r)
	x, err := io.ReadAll(r.Body)
	if err != nil {
		return RenderError(w, err)
	}

	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	dataset, err := h.datasetService.GetDataset(id)

	if err != nil {
		return RenderError(w, err)
	}

	reader := csv.NewReader(bytes.NewReader(x))

	records, err := reader.ReadAll()

	if err != nil {
		return RenderError(w, err)
	}

	if !validateUploadedDataset(dataset, records) {
		return RenderError(w, fmt.Errorf("uploaded data does not fit the dataset schema"))
	}

	if err := h.datasetService.UploadData(id, x); err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NoContent())
}

func validateUploadedDataset(dataset entity.DatasetInfo, uploadedData [][]string) bool {
	colNames := getColNamesFromDataset(dataset)

	if len(uploadedData) < 1 {
		return false
	}
	csvColumns := uploadedData[0]

	return compareSchemaWithUploaded(colNames, csvColumns)
}

func getColNamesFromDataset(dataset entity.DatasetInfo) []string {
	var cols []string
	for _, colSchema := range dataset.Schema {
		cols = append(cols, colSchema.Name)
	}

	return cols
}

func compareSchemaWithUploaded(schemaCols []string, uploadedCols []string) bool {
	if len(schemaCols) != len(uploadedCols) {
		return false
	}

	c1 := make([]string, len(schemaCols))
	c2 := make([]string, len(uploadedCols))
	copy(c1, schemaCols)
	copy(c2, uploadedCols)

	sort.Strings(c1)
	sort.Strings(c2)

	for i := range c1 {
		if c1[i] != c2[i] {
			return false
		}
	}
	return true
}
