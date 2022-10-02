package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id        int
	Email     string
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	IsAdmin   int
	Tasks     []Task
	DB        *sql.DB
}

type Task struct {
	Id         int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsComplete bool
}

func (u *User) Insert(a *AppConfig) (sql.Result, error) {
	db := a.DB

	stmt := `
		INSERT INTO user (email, password, created_at, updated_at, is_admin) VALUES (?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP(), ?)
	`

	result, err := db.Exec(stmt, u.Email, u.Password, 0)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *User) GetSingleUser(a *AppConfig, s string) (*User, error) {

	stmt := `
		SELECT user_id, email, password, created_at, updated_at, is_admin FROM user WHERE email = ?;
	`

	row := a.DB.QueryRow(stmt, s)

	user := &User{}

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.IsAdmin)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("error fetching record")
		} else {
			return nil, err
		}

	}
	return user, nil
}

func (u *User) GetAllUsers(a *AppConfig) ([]*User, error) {

	stmt := `
		SELECT user_id, email, password, is_admin, created_at, updated_at FROM user;
	`

	rows, err := a.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		u := &User{}

		err := rows.Scan(&u.Id, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.IsAdmin)

		if err != nil {
			return nil, err
		}

		users := append(users, u)

		fmt.Println(users)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
