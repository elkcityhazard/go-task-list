package main

import (
	"fmt"
	"github.com/elkcityhazard/go-task-list/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.Home(w, r)
	}))

	mux.Handle("/signup", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			handlers.GetSignUp(w, r)
			break

		}

	}))

	mux.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	mux.Handle("/create", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	mux.Handle("/update", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	mux.Handle("/read", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	mux.Handle("/delete", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))

	mux.Handle("/view", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.URL.Path[len("/view/"):]

		fmt.Println(key)
	}))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return addDefaultHeaders(mux)
}
