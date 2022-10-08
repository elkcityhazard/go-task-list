package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	Id         int
	IsComplete bool
	Title      string
	Body       string
	UserId     int
	CommentId  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TaskList []*Task

func (t *Task) FetchAllTasksForUser(a *AppConfig, r *http.Request, uid string) (TaskList, error) {

	session, ok := a.SessionManager.Get(r.Context(), "id").(string)

	fmt.Println("Session ID: ", session)

	if !ok {
		err := errors.New("Error")
		return nil, err
	}

	stmt := `
		SELECT * FROM task where user_id = ? ORDER BY created_at ASC;	
	`

	rows, err := a.DB.Query(stmt, uid)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	var currentTasks = []*Task{}

	for rows.Next() {

		t := &Task{}

		err := rows.Scan(&t.Id, &t.IsComplete, &t.Title, &t.Body, &t.UserId, &t.CommentId, &t.CreatedAt, &t.UpdatedAt)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		currentTasks = append(currentTasks, t)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return currentTasks, nil

}

func (t *Task) FetchSingleTaskForUser(a *AppConfig, r *http.Request, uid string) (*Task, error) {

	session, ok := a.SessionManager.Get(r.Context(), "id").(string)

	fmt.Println(session)

	if !ok {
		err := errors.New("Error")
		return nil, err
	}

	stmt := `
	SELECT * FROM task where task_id = ?;
	`

	row := a.DB.QueryRow(stmt, uid)

	var task = &Task{}

	err := row.Scan(&task.Id, &task.IsComplete, &task.Title, &task.Body, &task.UserId, &task.CommentId, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		} else {
			return nil, err
		}
	}

	return task, nil
}

func (t *Task) DeleteSingleTask(a *AppConfig, r *http.Request, uid string) (sql.Result, error) {
	if !a.SessionManager.Exists(r.Context(), "id") {
		err := errors.New("invalid operation")

		return nil, err
	}

	stmt := `DELETE FROM task WHERE task_id = ?;`

	result, err := a.DB.Exec(stmt, uid)

	if err != nil {
		return nil, err
	}

	return result, nil
}
