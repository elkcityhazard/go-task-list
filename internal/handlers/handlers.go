package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	app.UserTasks = []*models.Task{}

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

	app.SessionManager.RenewToken(r.Context())

	app.SessionManager.Put(r.Context(), "id", strconv.Itoa(createdUser.Id))

	render.RenderTemplate(w, r, "signed-up.tmpl.html", &models.TemplateData{
		Title:   "Thanks For Signing Up",
		UserMap: userMap,
	})

}

func Signup(w http.ResponseWriter, r *http.Request) {
	app.UserTasks = []*models.Task{}
	if r.Method == "GET" {
		GetSignUp(w, r)
		return
	}

	if r.Method == "POST" {
		PostSignUp(w, r)
		return
	}

}

func Login(w http.ResponseWriter, r *http.Request) {

	app.UserTasks = []*models.Task{}

	switch r.Method {
	case "GET":
		if app.SessionManager.Exists(r.Context(), "id") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		render.RenderTemplate(w, r, "login.tmpl.html", &models.TemplateData{})
	case "POST":
		err := r.ParseForm()

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		email := r.Form.Get("email")

		ptPassword := r.Form.Get("password")

		if len(email) == 0 {
			http.Error(w, "invalid entry", http.StatusLengthRequired)
			return
		}

		if len(ptPassword) == 0 {
			http.Error(w, "invalid entry", http.StatusLengthRequired)
			return
		}

		var user *models.User

		user, err = user.GetSingleUser(app, email)

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ptPassword))

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		err = app.SessionManager.RenewToken(r.Context())
		if err != nil {
			fmt.Println(err)
			return
		}

		app.SessionManager.Put(r.Context(), "id", strconv.Itoa(user.Id))

		http.Redirect(w, r, "/new-task", http.StatusSeeOther)

		return

		// render.RenderTemplate(w, r, "welcome.tmpl.html", &models.TemplateData{
		// 	Data: user,
		// })

	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	app.UserTasks = []*models.Task{}

	if !app.SessionManager.Exists(r.Context(), "id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	if len(app.UserTasks) <= 0 {
		task := models.Task{}

		taskList, err := task.FetchAllTasksForUser(app, r, app.SessionManager.GetString(r.Context(), "id"))

		if err != nil {
			err = errors.New("error fetching task list")
			log.Println(err)
			return
		}

		app.UserTasks = taskList

		if len(app.UserTasks) == 0 {
			http.Redirect(w, r, "/new-task", http.StatusSeeOther)
			return
		}

	}

	var key = strings.TrimPrefix(r.URL.Path, "/tasks/")

	if len(key) > 0 {

		checkKey, err := strconv.Atoi(key)

		if err != nil {
			fmt.Println(err)
			//http.Redirect(w, r, "/tasks/", http.StatusSeeOther)
		}

		if checkKey < 0 {
			http.Redirect(w, r, fmt.Sprintf("/tasks/%d", 0), http.StatusSeeOther)
		}

		if checkKey > len(app.UserTasks)-1 {

			maxLength := len(app.UserTasks) - 1

			http.Redirect(w, r, fmt.Sprintf("/tasks/%d", maxLength), http.StatusSeeOther)
		}

		var taskList = app.UserTasks

		for _, v := range taskList {
			fmt.Println(v)
		}

		var task = &models.Task{}

		intKey, err := strconv.Atoi(key)

		if err != nil {
			fmt.Println(err)
			return
		}

		if intKey < 0 {
			intKey = 0
		}

		if intKey >= len(app.UserTasks)-1 {
			intKey = len(app.UserTasks) - 1
		}

		task = app.UserTasks[intKey]

		if strconv.Itoa(task.UserId) != app.SessionManager.Get(r.Context(), "id") {
			nextKey, err := strconv.Atoi(key)

			if err != nil {
				http.Redirect(w, r, "/tasks/", http.StatusSeeOther)
				return
			}

			http.Redirect(w, r, fmt.Sprintf("/tasks/%d", nextKey-1), http.StatusSeeOther)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		payload := make(map[string]interface{})

		payload["task"] = app.UserTasks[intKey]
		payload["intKey"] = intKey
		payload["length"] = len(app.UserTasks)
		payload["lastIndex"] = len(app.UserTasks) - 1

		commentStmt := `
						SELECT * FROM comment where comment_id = ?;
						`

		result, err := app.DB.Query(commentStmt, task.Id)

		if err != nil {
			http.Error(w, "error processing comment fetch", http.StatusInternalServerError)
			return
		}

		defer result.Close()

		var commentList []*models.Comment

		for result.Next() {

			c := &models.Comment{}

			err := result.Scan(&c.ID, &c.UserId, &c.Title, &c.Body, &c.CreatedAt, &c.UpdatedAt, &c.CommentId)

			if err != nil {
				return
			}

			commentList = append(commentList, c)

		}

		if err = result.Err(); err != nil {
			return
		}

		payload["comments"] = commentList

		render.RenderTemplate(w, r, "single-task.tmpl.html", &models.TemplateData{
			Data: payload,
		})

		return

	}

	if !app.SessionManager.Exists(r.Context(), "id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, ok := app.SessionManager.Get(r.Context(), "id").(string)

	if !ok {
		http.Error(w, "invalid session data", http.StatusInternalServerError)
		return
	}

	var user *models.User

	query, err := user.GetSingleUserById(app, session)

	if err != nil {
		fmt.Println(err)
		return
	}

	queryID := strconv.Itoa(query.Id)

	if !app.SessionManager.Exists(r.Context(), "id") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if app.SessionManager.GetString(r.Context(), "id") != queryID {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	var payload = make(map[string]interface{})

	fmt.Println("before output", app.UserTasks)

	if len(app.UserTasks) <= 0 {
		http.Redirect(w, r, "/new-task", http.StatusSeeOther)
		return
	}

	payload["userTasks"] = app.UserTasks

	render.RenderTemplate(w, r, "task-list.tmpl.html", &models.TemplateData{
		Data: payload,
	})

}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	app.UserTasks = []*models.Task{}
	switch r.Method {
	case "GET":

		if !app.SessionManager.Exists(r.Context(), "id") {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		session, ok := app.SessionManager.Get(r.Context(), "id").(string)

		if !ok {
			return
		}

		fmt.Println(session)

		render.RenderTemplate(w, r, "new-task.tmpl.html", &models.TemplateData{})
	case "POST":

		if !app.SessionManager.Exists(r.Context(), "id") {
			http.Redirect(w, r, "/login", http.StatusNotFound)
			return
		}

		err := r.ParseForm()

		if err != nil {
			http.Redirect(w, r, "/new-task", http.StatusSeeOther)
			return
		}

		if !app.SessionManager.Exists(r.Context(), "id") {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		stmt := `
				INSERT INTO task (is_complete, title, body, user_id, comment_id, created_at, updated_at) VALUES 
				(?, ?, ?, ?,?, UTC_TIMESTAMP(), UTC_TIMESTAMP());
				`

		var task models.Task
		title := r.Form.Get("title")
		body := r.Form.Get("body")

		task.IsComplete = false
		task.Title = title
		task.Body = body
		task.CreatedAt = time.Now()
		task.UpdatedAt = time.Now()
		id, ok := app.SessionManager.Get(r.Context(), "id").(string)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		intID, err := strconv.Atoi(id)

		if err != nil {
			fmt.Println(err)
			return
		}

		task.UserId = intID

		_, err = app.DB.Exec(stmt, 0, title, body, task.UserId, 0)

		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		app.UserTasks = []*models.Task{}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return

	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	app.UserTasks = []*models.Task{}
	_ = app.SessionManager.Destroy(r.Context())

	_ = app.SessionManager.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func TaskAdmin(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		var task *models.Task

		task, err := task.FetchSingleTaskForUser(app, r, r.URL.Path[len("/admin/"):])

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		data := make(map[string]interface{})

		data["task"] = task

		fmt.Println(data["task"])

		render.RenderTemplate(w, r, "edit-task.tmpl.html", &models.TemplateData{
			Data: data,
		})
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "error parsing form", http.StatusInternalServerError)
		return
	}

	method := r.FormValue("_method")

	method = strings.ToLower(method)

	switch method {

	case "delete":
		checkSession(w, r)

		err := r.ParseForm()

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		id := r.Form.Get("id")

		var task *models.Task

		result, err := task.DeleteSingleTask(app, r, id)

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(result)

		http.Redirect(w, r, "/tasks/", http.StatusSeeOther)
		return
	}

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		err := r.ParseForm()

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		id := r.URL.Path[len("/admin/update/"):]
		status := r.Form.Get("status")
		title := r.Form.Get("title")
		body := r.Form.Get("body")

		if id == "" || title == "" || body == "" {
			http.Error(w, "sorry something went wrong", http.StatusInternalServerError)
			return
		}

		if status == "on" {
			status = strconv.Itoa(1)
		} else {
			status = strconv.Itoa(0)
		}

		fmt.Println(status)

		if !app.SessionManager.Exists(r.Context(), "id") {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		stmt := `UPDATE task 
				SET is_complete = ?, 
					title = ?, 
					body = ?, 
					updated_at = NOW()
				WHERE task_id = ?;`

		result, err := app.DB.Exec(stmt, status, title, body, id)

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(result)

		http.Redirect(w, r, "/tasks/", http.StatusSeeOther)

	}

}

func checkSession(w http.ResponseWriter, r *http.Request) {

	if !app.SessionManager.Exists(r.Context(), "id") {

		err := app.SessionManager.RenewToken(r.Context())

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusMethodNotAllowed)
	}

}

func CreateComment(w http.ResponseWriter, r *http.Request) {

	switch strings.ToUpper(r.Method) {
	case "POST":

		// handle task ID
		taskID := r.URL.Path[len("/tasks/comment/"):]

		if !app.SessionManager.Exists(r.Context(), "id") {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if taskID == "" {
			http.Error(w, "error processing taskID request", http.StatusInternalServerError)
			return
		}

		//	Handle User ID

		userID := app.SessionManager.Get(r.Context(), "id").(string)
		if userID == "" {
			http.Error(w, "error processing ID request", http.StatusInternalServerError)
			return
		}

		err := r.ParseForm()

		if err != nil {
			http.Error(w, "error processing Form request", http.StatusInternalServerError)
			return
		}

		title := r.Form.Get("title")
		taskComment := r.Form.Get("comment")

		var comment models.Comment

		comment.Title = title
		comment.Body = taskComment
		comment.CreatedAt = time.Now()
		comment.UpdatedAt = time.Now()
		comment.CommentId, _ = strconv.Atoi(taskID)
		comment.UserId, _ = strconv.Atoi(userID)

		result, err := comment.NewComment(app, comment.UserId, comment.Title, comment.Body, comment.CommentId)

		if err != nil {
			http.Error(w, "error processing new comment db request request", http.StatusInternalServerError)
			return
		}

		fmt.Println(result)

		http.Redirect(w, r, fmt.Sprintf("%s", r.Referer()), http.StatusSeeOther)

	}
}
