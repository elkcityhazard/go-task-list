package config

import (
	"database/sql"
	"html/template"
)

type TemplateData struct {
	SiteTitle      string
	Title          string
	MainNavigation []NavItem
}

type AppConfig struct {
	Port           string
	DSN            string
	TemplateCache  map[string]*template.Template
	TemplateData   TemplateData
	MainNavigation []NavItem
	DB             *sql.DB
	IsProduction   bool
}

type NavItem struct {
	URL    string
	Name   string
	Weight int
}

func (a *AppConfig) InitializeMenu() []NavItem {
	a.MainNavigation = []NavItem{
		{
			URL:    "/",
			Name:   "Home",
			Weight: 1,
		},
		{
			URL:    "/signup",
			Name:   "Signup",
			Weight: 2,
		},
		{
			URL:    "/login",
			Name:   "Login",
			Weight: 3,
		},
	}

	return a.MainNavigation
}
