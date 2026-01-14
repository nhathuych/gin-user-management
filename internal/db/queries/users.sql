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

-- name: GetUser :one
SELECT * FROM users WHERE uuid = $1 AND deleted_at IS NULL LIMIT 1;

-- name: SoftDeleteUser :one
UPDATE users
SET
  deleted_at = NOW()
WHERE
  uuid = sqlc.arg(uuid)::uuid AND
  deleted_at IS NULL
RETURNING *;

-- name: RestoreUser :one
UPDATE users
SET
  deleted_at = NULL
WHERE
  uuid = sqlc.arg(uuid)::uuid AND
  deleted_at IS NOT NULL
RETURNING *;

-- name: HardDeleteUser :one
DELETE FROM users
WHERE
  uuid = sqlc.arg(uuid)::uuid AND
  deleted_at IS NOT NULL
RETURNING *;

-- name: CountUsers :one
SELECT COUNT(*)
FROM users
WHERE
  deleted_at IS NULL AND
  (
    sqlc.narg(search)::text IS NULL OR
    sqlc.narg(search)::text = '' OR
    email ILIKE '%' || sqlc.narg(search) || '%' OR
    fullname ILIKE '%' || sqlc.narg(search) || '%'
  );

-- name: ListUsersOrderByIdASC :many
SELECT *
FROM users
WHERE
  deleted_at IS NULL AND
  (
    sqlc.narg(search)::text IS NULL OR
    sqlc.narg(search)::text = '' OR
    email ILIKE '%' || sqlc.narg(search) || '%' OR
    fullname ILIKE '%' || sqlc.narg(search) || '%'
  )
ORDER BY id ASC
LIMIT $1 OFFSET $2;

-- name: ListUsersOrderByIdDESC :many
SELECT *
FROM users
WHERE
  deleted_at IS NULL AND
  (
    sqlc.narg(search)::text IS NULL OR
    sqlc.narg(search)::text = '' OR
    email ILIKE '%' || sqlc.narg(search) || '%' OR
    fullname ILIKE '%' || sqlc.narg(search) || '%'
  )
ORDER BY id DESC
LIMIT $1 OFFSET $2;
