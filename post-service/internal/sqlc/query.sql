-- name: GetAll :many
SELECT * FROM posts;

-- name: CreatePost :exec
INSERT INTO posts (owner_id,body,media_id) VALUES ($1,$2,$3);

-- name: GetPostByUsername :one
SELECT * FROM posts WHERE id = $1;
