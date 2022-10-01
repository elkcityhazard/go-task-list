package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"personal.github.com/elkcityhazard/go-task-list/internal/config"
	"strconv"
)

var app config.AppConfig

func main() {

	app.IsProduction = false

	flag.StringVar(&app.Port, "port", ":8080", "the port the application listens on")
	isProduction := flag.String("Environment", "false", "Set the environment [production | development ]")
	flag.Parse()

	app.IsProduction, _ = strconv.ParseBool(*isProduction)

	srv := &http.Server{
		Addr:    app.Port,
		Handler: routes(),
	}

	fmt.Printf("starting server on %s", srv.Addr)

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

}
