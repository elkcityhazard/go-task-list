package models

import (
	"database/sql"
	"html/template"

	"github.com/alexedwards/scs/v2"
)

type TemplateData struct {
	SiteTitle      string
	Title          string
	MainNavigation []NavItem
	Data           interface{}
	StringMap      map[string]string
	UserMap        map[string]User
}

type AppConfig struct {
	Port           string
	DSN            string
	TemplateCache  map[string]*template.Template
	TemplateData   TemplateData
	MainNavigation []NavItem
	DB             *sql.DB
	IsProduction   bool
	SessionManager *scs.SessionManager
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
