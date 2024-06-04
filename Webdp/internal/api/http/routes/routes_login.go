package routes

import (
	"webdp/internal/api/http/handlers"

	"github.com/gorilla/mux"
)

func RegisterLogin(router *mux.Router, handler handlers.LoginHandler) {
	router.HandleFunc("/login", handlers.HandlerDecorator(handler.LoginRequestHandler)).Methods("POST")
}

func RegisterLogout(router *mux.Router, handler handlers.LoginHandler) {
	router.HandleFunc("/logout", handlers.HandlerDecorator(handler.LogoutRequestHandler)).Methods("POST")
}
