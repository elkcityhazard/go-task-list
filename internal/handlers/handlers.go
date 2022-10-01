package handlers

import (
	"github.com/elkcityhazard/go-task-list/internal/config"
	"github.com/elkcityhazard/go-task-list/internal/render"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.tmpl.html", &config.TemplateData{Title: "Welcome Screen"})
}

func GetSignUp(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signup.tmpl.html", &config.TemplateData{Title: "Sign Up For Task List"})
}
