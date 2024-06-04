package main

import (
	"fmt"
	"googledp/middleware"
	"googledp/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//port := os.Getenv("PORT")

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	routes.RegisterRoutes(r)

	serv := fmt.Sprintf(":%s", "8000")

	log.Fatal(http.ListenAndServe(serv, r))

}
