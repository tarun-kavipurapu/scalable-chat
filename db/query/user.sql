-- name: CreateUser :one
INSERT INTO
    users (email, password, username)
VALUES 
 ($1,$2,$3) RETURNING *;

-- name: GetUserById :one

SELECT * FROM users where id = $1;

-- name: GetUserByEmail :one

SELECT * FROM users where email = $1;