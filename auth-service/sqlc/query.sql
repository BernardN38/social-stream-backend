-- name: GetAll :many
SELECT * FROM users;

-- name: CreateUser :exec
INSERT INTO USERS (username,email,encoded_password) VALUES ($1,$2,$3) ;

-- name: GetUser :one
SELECT username,email FROM users WHERE id = $1;

-- name: GetUserPassword :one
SELECT id,encoded_password FROM users WHERE email = $1;