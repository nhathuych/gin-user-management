-- name: CreateUser :one
INSERT INTO users(email, password, fullname, age, status, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;
