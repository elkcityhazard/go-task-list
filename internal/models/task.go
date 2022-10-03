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
	CommentId  sql.NullInt32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TaskList []*Task

func (t *Task) FetchAllTasksForUser(a *AppConfig, r *http.Request, uid string) (TaskList, error) {

	session, ok := a.SessionManager.Get(r.Context(), "id").(string)

	fmt.Println(session)

	if !ok {
		err := errors.New("Error")
		return nil, err
	}

	stmt := `
		SELECT * FROM task where user_id = ?;	
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

		fmt.Println(t)

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
