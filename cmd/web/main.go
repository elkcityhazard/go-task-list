package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elkcityhazard/go-task-list/internal/config"
	"github.com/elkcityhazard/go-task-list/internal/render"
	_ "github.com/go-sql-driver/mysql"
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
	flag.StringVar(&app.DSN, "DSN", "", "SQL Data Source Name")
	flag.BoolVar(&app.IsProduction, "Environment", false, "Set the environment [production | development ]")
	flag.Parse()

	db, err := sql.Open("mysql", app.DSN)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	app.DB = db
	app.InitializeMenu()

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
