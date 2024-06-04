package routes

import (
	"webdp/internal/api/http/handlers"

	"github.com/gorilla/mux"
)

func RegisterBudgetsV1(router *mux.Router, handler handlers.BudgetHandler) {
	budget := router.PathPrefix("/budget").Subrouter()
	budget.HandleFunc("/user/{userHandle}", handlers.HandlerDecorator(handler.GetUserBudgets)).Methods("GET")
	budget.HandleFunc("/dataset/{datasetId}", handlers.HandlerDecorator(handler.GetDatasetBudget)).Methods("GET")
	budget.HandleFunc("/allocation/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.GetUserDatasetBudget)).Methods("GET")
	budget.HandleFunc("/allocation/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.PostUserDatasetBudget)).Methods("POST")
	budget.HandleFunc("/allocation/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.PatchUserDatasetBudget)).Methods("PATCH")
	budget.HandleFunc("/allocation/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.DeleteUserDatasetBudget)).Methods("DELETE")
}

func RegisterBudgetsV2(router *mux.Router, handler handlers.BudgetHandler) {
	budget := router.PathPrefix("/budgets").Subrouter()
	budget.HandleFunc("/users/{userHandle}", handlers.HandlerDecorator(handler.GetUserBudgets)).Methods("GET")
	budget.HandleFunc("/datasets/{datasetId}", handlers.HandlerDecorator(handler.GetDatasetBudget)).Methods("GET")
	budget.HandleFunc("/allocations/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.GetUserDatasetBudget)).Methods("GET")
	budget.HandleFunc("/allocations/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.PostUserDatasetBudget)).Methods("POST")
	budget.HandleFunc("/allocations/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.PatchUserDatasetBudget)).Methods("PATCH")
	budget.HandleFunc("/allocations/{userHandle}/{datasetId}", handlers.HandlerDecorator(handler.DeleteUserDatasetBudget)).Methods("DELETE")
}
