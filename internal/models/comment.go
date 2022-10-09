package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id        int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    int
	User      User
}

func (c *Comment) CreateComment(a *AppConfig, title string, body string, createdAt time.Time, updatedAt time.Time, id int) (sql.Result, error) {

	stmt := `INSERT INTO comment (title, body, created_at, updated_at, user_id) VALUES (?,?,?,?,?);`

	result, err := a.DB.Exec(stmt, title, body, createdAt, updatedAt, id)

	if err != nil {
		return nil, err
	}

	return result, nil

}
