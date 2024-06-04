package routes

import (
	"webdp/internal/api/http/handlers"

	"github.com/gorilla/mux"
)

func RegisterUserV1(router *mux.Router, handler handlers.UserHandler) {
	users := router.PathPrefix("/users").Subrouter()
	users.HandleFunc("", handlers.HandlerDecorator(handler.GetUsers)).Methods("GET")
	users.HandleFunc("", handlers.HandlerDecorator(handler.PostUsers)).Methods("POST")

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{userHandle}", handlers.HandlerDecorator(handler.GetUser)).Methods("GET")
	user.HandleFunc("/{userHandle}", handlers.HandlerDecorator(handler.PatchUser)).Methods("PATCH")
	user.HandleFunc("/{userHandle}", handlers.HandlerDecorator(handler.DeleteUser)).Methods("DELETE")
}

func RegisterUserV2(router *mux.Router, handler handlers.UserHandler) {
	users := router.PathPrefix("/users").Subrouter()

	users.HandleFunc("", handlers.HandlerDecorator(handler.GetUsers)).Methods("GET")
	users.HandleFunc("", handlers.HandlerDecorator(handler.PostUsers)).Methods("POST")
	users.HandleFunc("/{userHandle}", handlers.HandlerDecorator(handler.GetUser)).Methods("GET")
	users.HandleFunc("/{userHandle}", handlers.HandlerDecorator(handler.PatchUser)).Methods("PATCH")
	users.HandleFunc("/{userHandle}", handlers.HandlerDecorator(handler.DeleteUser)).Methods("DELETE")
}
