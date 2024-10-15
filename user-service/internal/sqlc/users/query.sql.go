// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package users_sql

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (username,email,first_name,last_name,dob) VALUES ($1,$2,$3,$4,$5)
`

type CreateUserParams struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Dob       string `json:"dob"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.Username,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.Dob,
	)
	return err
}

const getAll = `-- name: GetAll :many
SELECT id, username, email, first_name, last_name, dob FROM users
`

func (q *Queries) GetAll(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.FirstName,
			&i.LastName,
			&i.Dob,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, first_name, last_name, dob FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Dob,
	)
	return i, err
}
