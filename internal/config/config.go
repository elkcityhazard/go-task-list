package config

import (
	"database/sql"
	"html/template"
)

type TemplateData struct {
	Title string
}

type AppConfig struct {
	Port                string
	DSN                 string
	TemplateCache       map[string]*template.Template
	DefaultTemplateData struct{}
	DB                  *sql.DB
	IsProduction        bool
}
