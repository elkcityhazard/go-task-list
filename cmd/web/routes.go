package main

import (
	"fmt"
	"net/http"

	"github.com/elkcityhazard/go-task-list/internal/config"
	"github.com/elkcityhazard/go-task-list/internal/render"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.RenderTemplate(w, r, "home.tmpl.html", config.TemplateData{Title: "Hello World"})
	}))

	mux.Handle("/signup", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":

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
