package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"webdp/internal/api/http/client"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/handlers"
	"webdp/internal/api/http/middlewares"
	"webdp/internal/api/http/repo/postgres"
	"webdp/internal/api/http/routes"
	"webdp/internal/api/http/services"
	"webdp/internal/config"
	"webdp/internal/config/dbconnection"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// api versions
const (
	VERSION_1 = "/v1"
	VERSION_2 = "/v2"
)

type repo struct {
	users    postgres.UserPostgres
	tokens   postgres.TokenPostgres
	datasets postgres.DatasetPostgres
	budgets  postgres.BudgetPostgres
}

type service struct {
	users      services.UserService
	mockTokens services.TokenService
	realTokens services.TokenService
	datasets   services.DatasetService
	budgets    services.BudgetService
}

type handler struct {
	users    handlers.UserHandler
	login    handlers.LoginHandler
	datasets handlers.DatasetHandler
	budgets  handlers.BudgetHandler
	queries  handlers.QueryHandler
}

// @title Webdp API - Reworked
// @version 2.0
// @description Welcome to the official OpenAPI documentation for WebDP, our versatile API designed to provide transparent interoperability with a range of differentially private frameworks.

// @securityDefinitions.basic BearerTokenAuth

// @license.name Mozilla Public License version 2.0
// @license.url https://www.mozilla.org/en-US/MPL/2.0/

// @host localhost:8000
// @BasePath /
func main() {
	// collect environment variables
	env, err := config.Getenvars()
	if err != nil {
		panic(err)
	}
	engines, err := config.GetDpEngines(env.Config_path)
	if err != nil {
		panic(err)
	}

	// initiate database and client
	pg, err := dbconnection.ConnectPostgresDB(env.Db_name, env.Db_user, env.Db_pw, env.Db_host, env.Dp_port)
	if err != nil {
		panic(err)
	}
	durl := fmt.Sprintf("http://webdp-api:%s/datasets", env.Port_int)
	client := client.NewDPClient(*engines, durl, nil)
	repos, services, handlers := initRSH(pg, client, env.Auth_key)

	// external routes
	router := mux.NewRouter()
	router.Use(middlewares.Logger)
	registerExRoutes(router, VERSION_1, &repos.tokens, handlers)
	registerExRoutes(router, VERSION_2, &repos.tokens, handlers)

	// internal routes
	internalRouter := mux.NewRouter()
	registerInRoutes(internalRouter, &repos.datasets)

	// start internal server
	intServ := fmt.Sprintf(":%s", env.Port_int)
	go http.ListenAndServe(intServ, internalRouter)

	// create root user
	go func() {
		defer fmt.Printf("\n==============\nLogin with username: %s, password: %s\n==============\n", "root", env.Root_pw)
		createRootUser(pg, services.users, env.Root_pw)
	}()

	// start server
	serv := fmt.Sprintf(":%s", env.Port_ext)
	log.Fatal(http.ListenAndServe(serv, router))
}

/*
Initiates repos, services and handlers
*/
func initRSH(db *sql.DB, client *client.DPClient, skey string) (*repo, *service, *handler) {

	// repos
	repo := &repo{
		users:    postgres.NewUserPostgres(db),
		tokens:   postgres.NewTokenPostgres(db),
		datasets: postgres.NewDatasetPostgres(db),
		budgets:  postgres.NewBudgetPostgres(db),
	}

	// services
	service := &service{
		users:      services.NewUserService(repo.users),
		mockTokens: services.NewTokenService(repo.tokens, []byte(""), jwt.SigningMethodNone),
		realTokens: services.NewTokenService(repo.tokens, []byte(skey), jwt.SigningMethodHS256),
		datasets:   services.NewDatasetService(repo.datasets, repo.budgets),
		budgets:    services.NewBudgetService(repo.budgets),
	}

	// handlers
	handler := &handler{
		users:    handlers.NewUserHandler(service.users, service.mockTokens),
		datasets: handlers.NewDatasetHandler(service.datasets, service.users, service.budgets, *client),
		login:    handlers.NewLoginHandler(service.users, service.realTokens),
		budgets:  handlers.NewBudgetHandler(service.budgets, service.datasets),
		queries:  handlers.NewQueryHandler(service.datasets, service.budgets, *client),
	}

	return repo, service, handler
}

/*
Registers external routes
*/
func registerExRoutes(r *mux.Router, version string, tokenrepo *postgres.TokenPostgres, handler *handler) {
	router := r.PathPrefix(version).Subrouter()
	notoken := router.PathPrefix("").Subrouter()
	token := router.PathPrefix("").Subrouter()
	token.Use(middlewares.GetTokenAuthentication(*tokenrepo))

	routes.RegisterLogout(token, handler.login)
	routes.RegisterLogin(notoken, handler.login)

	if version == VERSION_1 {
		routes.RegisterUserV1(token, handler.users)
		routes.RegisterDatasetsV1(token, handler.datasets)
		routes.RegisterBudgetsV1(token, handler.budgets)
		routes.RegisterQueriesV1(token, handler.queries)
	} else if version == VERSION_2 {
		routes.RegisterUserV2(token, handler.users)
		routes.RegisterDatasetsV2(token, handler.datasets)
		routes.RegisterBudgetsV2(token, handler.budgets)
		routes.RegisterQueriesV2(token, handler.queries)
		routes.RegisterSpec(router) // no auth for this one
	}
}

/*
Registers internal routes.
*/
func registerInRoutes(ir *mux.Router, datasetrepo *postgres.DatasetPostgres) {
	serv := services.NewInternalDatasetService(*datasetrepo)
	hand := handlers.NewInternalDatasetHandler(*serv)
	routes.RegisterInternalDatasets(ir, *hand)
}

/*
Creates new root user.
First deletes any old root user.
*/
func createRootUser(pg *sql.DB, s services.UserService, pw string) {
	for {
		if err := pg.Ping(); err != nil {
			fmt.Println("Waiting for database connection...")
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("Database connection established.")
			break
		}
	}
	s.DeleteUser("root")
	root := entity.UserPost{
		Handle: "root",
		PWD:    pw,
		Name:   "root",
		Roles:  []string{entity.ADMIN, entity.CURATOR, entity.ANALYST},
	}
	s.CreateUser(root)
}
