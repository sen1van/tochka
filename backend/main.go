package main

import (
	"net/http"
	"tochka/database"
	"tochka/handlers/api"
	"tochka/middlewares"
)

func main() {
	db := database.NewDB("db.sqlite")
	defer db.Close()

	mux := http.NewServeMux()

	apiHandler := api.GetApiHandler()
	apiHandler = middlewares.AuthMiddleware(apiHandler)
	apiHandler = middlewares.APIMiddleware(apiHandler, db)
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	err := http.ListenAndServe(":8080", middlewares.LoggingMiddleware(mux))
	if err != nil {
		panic(err)
	}
}
