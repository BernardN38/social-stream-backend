-- name: GetAll :many
SELECT * FROM users;

-- name: CreateUser :exec
INSERT INTO users (username,email,first_name,last_name,dob) VALUES ($1,$2,$3,$4,$5);

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;
