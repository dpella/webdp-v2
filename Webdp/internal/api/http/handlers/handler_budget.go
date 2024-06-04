package handlers

import (
	"net/http"
	"strconv"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/middlewares"
	"webdp/internal/api/http/response"
	"webdp/internal/api/http/services"
	"webdp/internal/api/http/utils"

	"github.com/gorilla/mux"
)

type BudgetHandler struct {
	budgetService  services.BudgetService
	datasetService services.DatasetService
}

func NewBudgetHandler(bs services.BudgetService, ds services.DatasetService) BudgetHandler {
	return BudgetHandler{budgetService: bs, datasetService: ds}
}

/*
Gets budgets for user.
Request parameters: User handle.
Response: List of user budgets (dataset id, allocated, consumed).
Requester can get own budgets. For others, requester needs curator role.
*/
// GetBudgets godoc
// @Summary      Get budgets for user
// @Description  Gets all budgets that are allocated to the user
// @Tags         budgets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        userHandle   path      string  true  "User Handle"
// @Success      200  {object}  []entity.UserBudgetsResponse
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/budget/user/{userHandle} [get]
// @Router       /v2/budgets/users/{userHandle} [get]
func (h BudgetHandler) GetUserBudgets(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR}); err != nil {
		if err := middlewares.ValidateSelfRequest(r); err != nil {
			return RenderError(w, err)
		}
	}

	vars := mux.Vars(r)
	budgets, err := h.budgetService.GetUserBudgets(vars["userHandle"])
	if err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, budgets))
}

/*
Gets dataset budget for a dataset.
Request parameters: Dataset id.
Response: Dataset budget allocations.
Requester needs curator role, or granted access to a dataset.
*/
// GetDatasetBudgets godoc
// @Summary      Gets dataset budget for a dataset
// @Description  Gets the dataset budget allocation
// @Tags         budgets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        datasetId   path      int  true  "Dataset Id"
// @Success      200  {object}  entity.DatasetBudgetAllocationResponse
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/budget/dataset/{datasetId} [get]
// @Router       /v2/budgets/datasets/{datasetId} [get]
func (h BudgetHandler) GetDatasetBudget(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR}); err != nil {
		if err := middlewares.ValidateGrantedAccess(r, &h.budgetService); err != nil {
			return RenderError(w, err)
		}
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	dataset, err := h.budgetService.GetDatasetBudget(id)
	if err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, dataset))

}

/*
Gets budget for user and dataset.
Request parameters: User handle, dataset id
Response: Budget for user and dataset.
Requester can get own budgets. For others, requester needs curator role.
*/
// GetUserDatasetBudget godoc
// @Summary      Gets user budget on a dataset
// @Description  Gets the specified users budget on a specific dataset
// @Tags         budgets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        userHandle   path      string  true  "User Handle"
// @Param        datasetId   path      int  true  "Dataset Id"
// @Success      200  {object}  entity.Budget
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/budget/allocation/{userHandle}/{datasetId} [get]
// @Router       /v2/budgets/allocations/{userHandle}/{datasetId} [get]
func (h BudgetHandler) GetUserDatasetBudget(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR}); err != nil {
		if err := middlewares.ValidateSelfRequest(r); err != nil {
			return RenderError(w, err)
		}
	}

	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	budget, err := h.budgetService.GetUserDatasetBudget(vars["userHandle"], id)
	if err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, budget))
}

/*
Creates user budget for dataset.
Request parameters: User handle, dataset id
Request body: New allocation.
Requester needs curator or analyst role, and needs to be the owner of the dataset.
*/
// PostUserDatasetBudget godoc
// @Summary      Adds a user budget on a dataset
// @Description  Adds a user budget on a dataset
// @Tags         budgets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        userHandle   	path   string  		  true  "User Handle"
// @Param        datasetId   	path   int 			  true  "Dataset Id"
// @Param		 requestBody	body   entity.Budget  true  "request body"
// @Success      201
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/budget/allocation/{userHandle}/{datasetId} [post]
// @Router       /v2/budgets/allocations/{userHandle}/{datasetId} [post]
func (h BudgetHandler) PostUserDatasetBudget(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR, entity.ANALYST}); err != nil {
		return RenderError(w, err)
	}
	if err := middlewares.ValidateOwnership(r, &h.datasetService); err != nil {
		return RenderError(w, err)
	}

	var addBudget entity.Budget
	if err := utils.ParseJsonRequestBody[entity.Budget](r, &addBudget); err != nil {
		return RenderError(w, err)
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	if err := h.budgetService.PostUserDatasetBudget(vars["userHandle"], id, addBudget); err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.EmptyResponse(http.StatusCreated))
}

/*
Updates budget for user and dataset.
Request parameters: User handle, dataset id.
Request body: New allocation.
Requester needs to be the owner of the dataset.
// Note: user needs to have allocated budget (otherwise use post?). TODO: Assess if this needs intervention.
*/
// PatchUserDatasetBudget godoc
// @Summary      Update a user budget on a dataset
// @Description  Update a user budget on a dataset
// @Tags         budgets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        userHandle   path      string  true  "User Handle"
// @Param        datasetId   path      int  true  "Dataset Id"
// @Param		 requestBody	body   entity.Budget  true  "request body"
// @Success      204
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/budget/allocation/{userHandle}/{datasetId} [patch]
// @Router       /v2/budgets/allocations/{userHandle}/{datasetId} [patch]
func (h BudgetHandler) PatchUserDatasetBudget(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateOwnership(r, &h.datasetService); err != nil {
		return RenderError(w, err)
	}

	var patchBudget entity.Budget
	if err := utils.ParseJsonRequestBody[entity.Budget](r, &patchBudget); err != nil {
		return RenderError(w, err)
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	if err := h.budgetService.PatchUserDatasetBudget(vars["userHandle"], id, patchBudget); err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NoContent())
}

/*
Deletes budget for user and dataset.
Request parameters: User handle, dataset id.
Requester needs to be the owner of the dataset.
*/
// DeleteUserDatasetBudget godoc
// @Summary      Deletes budget for user and dataset.
// @Description  Deletes budget for user and dataset.
// @Tags         budgets
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        userHandle   path      string  true  "User Handle"
// @Param        datasetId   path      int  true  "Dataset Id"
// @Success      204
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/budget/allocation/{userHandle}/{datasetId} [delete]
// @Router       /v2/budgets/allocations/{userHandle}/{datasetId} [delete]
func (h BudgetHandler) DeleteUserDatasetBudget(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateOwnership(r, &h.datasetService); err != nil {
		return RenderError(w, err)
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["datasetId"], 10, 64)
	if err != nil {
		return RenderError(w, err)
	}

	if err := h.budgetService.DeleteUserDatasetBudget(vars["userHandle"], id); err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NoContent())
}
