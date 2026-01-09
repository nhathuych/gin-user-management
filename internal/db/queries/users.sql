-- name: CreateUser :one
INSERT INTO users(email, password, fullname, age, status, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  password = COALESCE(sqlc.narg(password), password),
  fullname = COALESCE(sqlc.narg(fullname), fullname),
  age = COALESCE(sqlc.narg(age), age),
  status = COALESCE(sqlc.narg(status), status),
  role = COALESCE(sqlc.narg(role), role)
WHERE
  uuid = sqlc.arg(uuid)::uuid AND
  deleted_at IS NULL
RETURNING *;
