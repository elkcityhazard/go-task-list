package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/elkcityhazard/go-task-list/internal/config"
	"github.com/elkcityhazard/go-task-list/internal/render"
)

var app config.AppConfig

func main() {

	app.IsProduction = false

	td, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatalln(err)
	}

	app.TemplateCache = td

	flag.StringVar(&app.Port, "port", ":8080", "the port the application listens on")
	isProduction := flag.String("Environment", "false", "Set the environment [production | development ]")
	flag.Parse()

	app.IsProduction, _ = strconv.ParseBool(*isProduction)

	srv := &http.Server{
		Addr:    app.Port,
		Handler: routes(),
	}

	render.NewRenderer(&app)

	fmt.Printf("starting server on %s", srv.Addr)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

}
