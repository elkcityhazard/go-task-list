package config

import (
	"database/sql"
	"html/template"
)

type AppConfig struct {
	Port          string
	DSN           string
	TemplateCache map[string]*template.Template
	DB            *sql.DB
	IsProduction  bool
}
