package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int
	UserId    int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	CommentId int
}

func (c *Comment) NewComment(a *AppConfig, userId int, title string, body string, taskId int) (sql.Result, error) {
	c.UserId = userId
	c.Title = title
	c.Body = body
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.CommentId = taskId

	stmt := `
	INSERT INTO comment (user_id, title, body, comment_id) VALUES (?, ?, ?, ?);
`

	result, err := a.DB.Exec(stmt, c.UserId, c.Title, c.Body, c.CommentId)

	if err != nil {
		return nil, err
	}

	a.UserTasks = []*Task{}

	return result, nil
}
