package main

import (
	"net/http"
	"tochka/database"
	"tochka/handlers"
	"tochka/middlewares"
)

func main() {
	db := database.NewDB("db.sqlite")

	http.HandleFunc("/ping", handlers.Ping) // GET POST

	api := http.NewServeMux()
	api.HandleFunc("/", handlers.APIHandler)

	http.Handle("/api/", middlewares.APIMiddleware(api, db))

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
