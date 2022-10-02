package main

import (
	"net/http"

	"github.com/elkcityhazard/go-task-list/internal/handlers"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/", addSessionManager(http.HandlerFunc(handlers.Home)))

	mux.Handle("/signup", addSessionManager(http.HandlerFunc(handlers.Signup)))

	mux.Handle("/tasks", addSessionManager(http.HandlerFunc(handlers.GetAllTasks)))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return addDefaultHeaders(mux)

}
