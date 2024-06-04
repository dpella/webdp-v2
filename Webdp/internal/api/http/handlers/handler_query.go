package handlers

import (
	"fmt"
	"net/http"
	"strings"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/client"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/middlewares"
	"webdp/internal/api/http/response"
	"webdp/internal/api/http/services"
	"webdp/internal/api/http/utils"
)

type QueryHandler struct {
	dataset services.DatasetService
	budget  services.BudgetService
	client  client.DPClient
}

func NewQueryHandler(dataset services.DatasetService, budget services.BudgetService, cli client.DPClient) QueryHandler {
	return QueryHandler{dataset: dataset, budget: budget, client: cli}
}

// PostQueryEvaluate godoc
// @Summary      Do a query evaluation
// @Description  Request a query evaluation on a specific dataset.
// @Description  Requester must be curator or analyst.
// @Tags         queries
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 queryEvaluate 	body 	entity.QueryEvaluate true "Query Evaluation Request"
// @Param        engine   	query   string 			  false  "engine name"
// @Success      200  {object}  entity.QueryResult
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/query/evaluate [post]
// @Router       /v2/queries/evaluate [post]
func (h QueryHandler) PostQueryEvaluate(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR, entity.ANALYST}); err != nil {
		return RenderError(w, err)
	}

	user, ok := r.Context().Value(middlewares.DPContextKey{Key: middlewares.UserContextKey}).(string)
	if !ok {
		return RenderError(w, errors.ErrUnexpected)
	}
	var query entity.QueryEvaluate
	if err := utils.ParseJsonRequestBody[entity.QueryEvaluate](r, &query); err != nil {
		return RenderError(w, err)
	}

	err := query.Valid()

	if err != nil {
		return RenderError(w, err)
	}

	datainfo, err := h.dataset.GetDataset(query.Dataset)
	if err != nil {
		return RenderError(w, err)
	}

	if !h.budget.HasUserEnoughBudget(user, query.Dataset, query.Budget) {
		return RenderError(w, fmt.Errorf("%w: not have budget for making the query", errors.ErrBadRequest))
	}

	if !datainfo.Loaded {
		return RenderError(w, fmt.Errorf("%w: cannot make a query without data. dataset %d not loaded", errors.ErrBadRequest, query.Dataset))
	}

	req := entity.QueryFromClientEvaluate{
		Data:          query.Dataset,
		Budget:        query.Budget,
		Query:         query.Query,
		Schema:        datainfo.Schema,
		PrivacyNotion: datainfo.PrivacyNotion,
	}

	engine := r.URL.Query().Get("engine")

	if engine == "" {
		engine = h.client.DefaultEngine
	}

	resp, err := h.client.EvaluateQuery(engine, req)
	if err != nil {
		return RenderError(w, err)
	}

	// query was ok so we update budget
	if err := h.budget.AddConsumedBudgetToUser(user, query.Dataset, query.Budget); err != nil {
		return RenderError(w, err)
	}

	// budget is updated and we are happy
	return RenderResponse(w, response.NewSuccess(http.StatusOK, resp))

}

// PostQueryCustom godoc
// @Summary      Do a custom query (not implemented)
// @Description  Custom query on a specific dataset
// @Tags         queries
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 queryCustom body entity.QueryCustom true "Query Custom Request"
// @Failure      501  {object}  response.Error
// @Router       /v1/query/custom [post]
func (h QueryHandler) PostQueryCustom(w http.ResponseWriter, r *http.Request) error {
	return RenderError(w, errors.ErrNotImplemented)
}

// PostQueryValidate godoc
// @Summary      Validate a query
// @Description  Validate a query's syntax.
// @Description  Requester must be curator or analyst.
// @Tags         queries
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 queryEvaluate 	body 	entity.QueryEvaluate true "Query Evaluation Request"
// @Param        engine   	query   string 			  false  "engine name"
// @Success      200  {object}  client.ValidateResponse
// @Success      200  {object}  []client.ValidateResponse
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v2/queries/validate [post]
func (h QueryHandler) PostQueryValidate(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR, entity.ANALYST}); err != nil {
		return RenderError(w, err)
	}

	user, ok := r.Context().Value(middlewares.DPContextKey{Key: middlewares.UserContextKey}).(string)
	if !ok {
		return RenderError(w, errors.ErrUnexpected)
	}
	var query entity.QueryEvaluate
	if err := utils.ParseJsonRequestBody[entity.QueryEvaluate](r, &query); err != nil {
		return RenderError(w, err)
	}

	err := query.Valid()

	if err != nil {
		return RenderError(w, err)
	}

	datainfo, err := h.dataset.GetDataset(query.Dataset)
	if err != nil {
		return RenderError(w, err)
	}

	if !datainfo.Loaded {
		return RenderError(w, fmt.Errorf("%w: cannot make a query without data. dataset %d not loaded", errors.ErrBadRequest, query.Dataset))
	}

	_, err = h.budget.GetUserDatasetBudget(user, datainfo.Id)

	if err != nil {
		return RenderError(w, fmt.Errorf("user %s does not have budget allocated on dataset %d", user, query.Dataset))
	}

	req := entity.QueryFromClientEvaluate{
		Data:          query.Dataset,
		Budget:        query.Budget,
		Query:         query.Query,
		Schema:        datainfo.Schema,
		PrivacyNotion: datainfo.PrivacyNotion,
	}

	engine := r.URL.Query().Get("engine")

	var resp interface{}

	if engine == "" {
		resp, err = h.client.ValidateQueryAll(req)
	} else {
		resp, err = h.client.ValidateQuery(engine, req)
	}

	if err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, resp))
}

// PostQueryAccuracy godoc
// @Summary      Check a query's accuracy
// @Description  Request query accuracy on a specific dataset.
// @Description  Requester must be curator or analyst.
// @Tags         queries
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param		 queryAccuracy body entity.QueryAccuracy true "Query Accuracy Request"
// @Param        engine   	query   string 			  false  "engine name"
// @Success      200  {object}  []float64
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v1/query/accuracy [post]
// @Router       /v2/queries/accuracy [post]
func (h QueryHandler) PostQueryAccuracy(w http.ResponseWriter, r *http.Request) error {
	if err := middlewares.ValidateRoles(r, []string{entity.CURATOR, entity.ANALYST}); err != nil {
		return RenderError(w, err)
	}

	user, ok := r.Context().Value(middlewares.DPContextKey{Key: middlewares.UserContextKey}).(string)
	if !ok {
		return RenderError(w, errors.ErrUnexpected)
	}

	var query entity.QueryAccuracy
	if err := utils.ParseJsonRequestBody[entity.QueryAccuracy](r, &query); err != nil {
		return RenderError(w, err)
	}

	err := query.Valid()

	if err != nil {
		return RenderError(w, err)
	}

	datainfo, err := h.dataset.GetDataset(query.Dataset)
	if err != nil {
		return RenderError(w, err)
	}

	if _, err = h.budget.GetUserDatasetBudget(user, datainfo.Id); err != nil {
		return RenderError(w, err)
	}

	if !datainfo.Loaded {
		return RenderError(w, fmt.Errorf("%w: cannot make a query without data. dataset %d not loaded", errors.ErrBadRequest, query.Dataset))
	}

	req := entity.QueryFromClientAccuracy{
		Data:          query.Dataset,
		Budget:        query.Budget,
		Query:         query.Query,
		Schema:        datainfo.Schema,
		PrivacyNotion: datainfo.PrivacyNotion,
		Confidence:    query.Confidence,
	}

	engine := r.URL.Query().Get("engine")

	if engine == "" {
		engine = h.client.DefaultEngine
	}

	res, err := h.client.GetQueryAccuracy(engine, req)

	if err != nil {
		return RenderError(w, err)
	}

	return RenderJsonResponse(w, 200, res)
}

// GetQueryEngines godoc
// @Summary      List available engines
// @Description  Returns a list of available engines.
// @Tags         queries
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  []string
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v2/queries/engines [get]
func (h QueryHandler) GetQueryEngines(w http.ResponseWriter, r *http.Request) error {
	engines := h.client.GetAvailableDPEngines()
	return RenderResponse(w, response.NewSuccess(http.StatusOK, engines))
}

// GetQueryFunctions godoc
// @Summary      List engine functions
// @Description  Returns a json with supported functions for each engine
// @Description  or single engine if specified in query param
// @Tags         queries
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      json
// @Param        engine   	query   string 			  false  "engine name"
// @Success      200  {object}  response.AllFunctions
// @Success      200  {object}  response.EngineFunctions
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v2/queries/functions [get]
func (h QueryHandler) GetQueryFunctions(w http.ResponseWriter, r *http.Request) error {
	engine := r.URL.Query().Get("engine")

	if engine != "" {
		if h.client.IsAvailable(engine) {
			res, err := h.client.GetSingleEngineFunctions(engine)
			if err != nil {
				return RenderError(w, err)
			}
			return RenderJsonResponse(w, 200, res, CONTENT_TYPE, APP_JSON)
		}
		return RenderError(w, fmt.Errorf("%w: supplied engine not available", errors.ErrBadRequest))
	}

	res, err := h.client.GetAllEngineFunctions()
	if err != nil {
		return RenderError(w, err)
	}

	return RenderResponse(w, response.NewSuccess(http.StatusOK, res))

}

// GetQueryDocsEngine godoc
// @Summary      Get engine query documentation
// @Description  Returns a markdown file with features for each engine
// @Description  or single engine if specified in query param
// @Tags         queries
// @Security 	 BearerTokenAuth
// @Accept       json
// @Produce      text/markdown
// @Param        engine   	query   string 			  false  "engine name"
// @Success      200  {object}  []byte
// @Failure      400  {object}  response.Error
// @Failure      403  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /v2/queries/docs [get]
func (h QueryHandler) GetQueryDocs(w http.ResponseWriter, r *http.Request) error {
	engine := r.URL.Query().Get("engine")

	if engine != "" {
		if h.client.IsAvailable(engine) {
			res, err := h.client.GetDocumentation(engine)
			if err != nil {
				return RenderError(w, err)
			}
			return RenderMDResponse(w, 200, res)
		}
		return RenderError(w, fmt.Errorf("%w: supplied engine not available", errors.ErrBadRequest))
	}

	allEngines := h.client.GetAvailableDPEngines()

	var combinedDocs strings.Builder

	for _, eng := range allEngines {
		res, _ := h.client.GetDocumentation(eng)

		combinedDocs.WriteString(string(res))
		combinedDocs.WriteString("\n")
		combinedDocs.WriteString("***")
		combinedDocs.WriteString("\n")
	}

	return RenderMDResponse(w, http.StatusOK, []byte(combinedDocs.String()))
}
