package main

import (
	"errors"
	"net/http"
	"time"
	"tochka/database"
	"tochka/handlers/api"
	"tochka/middlewares"
)

const (
	port              = ":8080"
	readTimeout       = 5 * time.Second
	writeTimeout      = 10 * time.Second
	readHeaderTimeout = 3 * time.Second
	idleTimeout       = 30 * time.Second
	maxHeaderBytes    = 1 << 20
)

func main() {
	db := database.NewDB("db.sqlite")
	defer db.Close()

	mux := http.NewServeMux()

	apiHandler := api.GetApiHandler()
	apiHandler = middlewares.AuthMiddleware(apiHandler)
	apiHandler = middlewares.APIMiddleware(apiHandler, db)
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	middlewares.LoggingMiddleware(mux)
	server := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
		TLSConfig:         nil,
	}

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
