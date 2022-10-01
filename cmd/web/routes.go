package main

import (
	"fmt"
	"net/http"
)

func routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello From Task List")
	}))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return addDefaultHeaders(mux)
}
