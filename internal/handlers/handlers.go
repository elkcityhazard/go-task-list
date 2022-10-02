package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/elkcityhazard/go-task-list/internal/models"
	"github.com/elkcityhazard/go-task-list/internal/render"
	"golang.org/x/crypto/bcrypt"
)

const (
	IsTrue  = 1
	IsFalse = 0
)

var app *models.AppConfig

func NewHandlers(a *models.AppConfig) {
	app = a
}

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.tmpl.html", &models.TemplateData{Title: "Welcome Screen"})
}

func GetSignUp(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signup.tmpl.html", &models.TemplateData{Title: "Sign Up For Task List"})
}

func PostSignUp(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "invalid request type", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	email := r.PostForm.Get("email")

	unsafePassword := r.PostForm.Get("password")

	if email == "" || unsafePassword == "" {
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(unsafePassword), 10)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	var user models.User

	user.Email = email
	user.Password = password
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsAdmin = IsFalse

	res, err := user.Insert(app)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	fmt.Println(res.RowsAffected())

	var userMap = make(map[string]models.User)

	userMap["current"] = user

	createdUser, err := user.GetSingleUser(app, email)

	fmt.Println(createdUser)

	if err != nil {
		fmt.Println("error thrown from createdUser: ", err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	fmt.Println(strconv.Itoa(createdUser.Id))

	app.SessionManager.Put(r.Context(), "id", strconv.Itoa(createdUser.Id))

	render.RenderTemplate(w, r, "signed-up.tmpl.html", &models.TemplateData{
		Title:   "Thanks For Signing Up",
		UserMap: userMap,
	})

}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetSignUp(w, r)
		return
	}

	if r.Method == "POST" {
		PostSignUp(w, r)
		return
	}

}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	session, ok := app.SessionManager.Get(r.Context(), "id").(string)

	if !ok {
		fmt.Println(session)
		http.Error(w, "invalid session data", http.StatusInternalServerError)
		return
	}

	var user *models.User

	query, err := user.GetSingleUserById(app, session)

	if err != nil {
		fmt.Println(err)
		return
	}

	var task *models.Task

	taskList, err := task.FetchAllTasksForUser(app, r, fmt.Sprintf("%d", query.Id))

	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	render.RenderTemplate(w, r, "task-list.tmpl.html", &models.TemplateData{
		Data: taskList,
	})

}
