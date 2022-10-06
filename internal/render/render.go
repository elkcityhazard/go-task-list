package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/elkcityhazard/go-task-list/internal/models"
)

var funcMap = template.FuncMap{
	"humanDate": humanDate,
	"prevTask":  prevTask,
	"nextTask":  nextTask,
}

func humanDate(t time.Time) string {

	return t.Format("02 Jan 2006 at 15:04")
}

func prevTask(id int) int {

	length := len(app.UserTasks)

	if id <= 0 {
		return 0
	}

	if id > length-1 {
		return id + 1
	}

	return id - 1
}

func nextTask(id int) int {

	length := len(app.UserTasks)

	if id > length-1 {
		return length - 1
	}

	return id + 1
}

var app *models.AppConfig

//	NewRenderer passes the current appconfig to the render package so it can
//	have access to the app config struct and whatever the current context is

func NewRenderer(a *models.AppConfig) {
	app = a
}

func AddDefaultTemplateData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.SiteTitle = "Go Task List"
	td.MainNavigation = app.InitializeMenu()

	if app.SessionManager.Exists(r.Context(), "id") {
		td.IsAuthenticated = 1
	}

	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.IsProduction {
		tc = app.TemplateCache
	}

	if !app.IsProduction {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatalln("error creating template cache")
		return
	}

	td = AddDefaultTemplateData(td, r)

	buf := new(bytes.Buffer)

	err := t.Execute(buf, td)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		fmt.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/pages/*.tmpl.html")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(page)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		matches, err := filepath.Glob("./templates/layouts/*.tmpl.html")

		if err != nil {
			fmt.Println(err)
		}

		if len(matches) > 0 {
			tmpl, err = tmpl.ParseGlob("./templates/layouts/*.tmpl.html")

			if err != nil {
				fmt.Println(err)
				return nil, err
			}

		}

		partials, err := filepath.Glob("./templates/partials/*.tmpl.html")

		if err != nil {
			return nil, err
		}

		if len(partials) > 0 {
			tmpl, err = tmpl.ParseGlob("./templates/partials/*.tmpl.html")

			if err != nil {
				return nil, err
			}
		}
		myCache[name] = tmpl
	}

	return myCache, nil

}
