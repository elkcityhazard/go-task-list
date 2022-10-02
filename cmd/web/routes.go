package main

import (
	"github.com/go-chi/chi"
	"net/http"

	"github.com/elkcityhazard/go-task-list/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(addSessionManager)
	mux.Use(addDefaultHeaders)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Get("/", handlers.Home)

	mux.Get("/signup", handlers.Signup)

	mux.Post("/signup", handlers.Signup)

	mux.Get("/tasks", handlers.GetAllTasks)

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}
