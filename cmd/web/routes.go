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

	mux.Handle("/login", addSessionManager(http.HandlerFunc(handlers.Login)))

	mux.Handle("/new-task", addSessionManager(http.HandlerFunc(handlers.CreateTask)))

	mux.Handle("/tasks/", addSessionManager(http.HandlerFunc(handlers.GetAllTasks)))

	mux.Handle("/tasks/comment/", addSessionManager(http.HandlerFunc(handlers.CreateComment)))

	mux.Handle("/admin/", addSessionManager(http.HandlerFunc(handlers.TaskAdmin)))
	mux.Handle("/admin/update/", addSessionManager(http.HandlerFunc(handlers.UpdateTask)))

	mux.Handle("/logout", addSessionManager(http.HandlerFunc(handlers.Logout)))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return pathToLowerCase(addDefaultHeaders(mux))

}
