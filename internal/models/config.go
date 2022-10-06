package models

import (
	"database/sql"
	"html/template"

	"github.com/alexedwards/scs/v2"
)

type TemplateData struct {
	SiteTitle       string
	Title           string
	MainNavigation  []NavItem
	Data            map[string]interface{}
	StringMap       map[string]string
	UserMap         map[string]User
	IsAuthenticated int
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
	UserTasks      []*Task
}

type NavItem struct {
	URL      string
	Name     string
	Weight   int
	LoggedIn bool
}

func (a *AppConfig) InitializeMenu() []NavItem {
	a.MainNavigation = []NavItem{
		{
			URL:    "/",
			Name:   "Home",
			Weight: 1,
		},
		{
			URL:      "/signup",
			Name:     "Signup",
			Weight:   2,
			LoggedIn: false,
		},
		{
			URL:      "/new-task",
			Name:     "New Task",
			Weight:   3,
			LoggedIn: true,
		},
	}

	return a.MainNavigation
}
