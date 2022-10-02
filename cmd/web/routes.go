package main

import (
	"fmt"
	"net/http"

	"github.com/elkcityhazard/go-task-list/internal/handlers"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.HandleFunc("/", handlers.Home)

	mux.HandleFunc("/signup", handlers.Signup)

	mux.Handle("/view", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.URL.Path[len("/view/"):]

		fmt.Println(key)
	}))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	addSessionManager(mux)

	return addDefaultHeaders(mux)
}
