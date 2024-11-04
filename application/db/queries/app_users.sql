-- name: CreateUser :one
INSERT INTO app_users (
  organization_id,
  username,
  display_name,
  email,
  password
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListUsersByOrg :many
SELECT u.* 
FROM app_users u
WHERE u.organization_id = $1
AND (u.display_name ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%' OR u.username ILIKE '%' || $2 || '%')
ORDER BY u.display_name
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountUsersByOrg :one
SELECT COUNT(*) as total 
FROM app_users u
WHERE u.organization_id = $1
AND (u.display_name ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%' OR u.username ILIKE '%' || $2 || '%');

-- name: GetUserByName :one
SELECT u.* 
FROM app_users u
WHERE u.username = $1
LIMIT 1;

-- name: GetUserByOrg :one
SELECT u.*
FROM app_users u
WHERE u.user_id = $1 AND u.organization_id = $2
LIMIT 1;

-- name: UpdateUser :exec
UPDATE app_users
SET
  display_name = $1,
  email = $2,
  updated_at = NOW()
WHERE user_id = $3 AND organization_id = $4
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM app_users
WHERE user_id = $1 AND organization_id = $2;

-- name: GetUserByID :one
SELECT u.*
FROM app_users u
WHERE u.user_id = $1
LIMIT 1;