-- name: CreateLoginHistory :exec
INSERT INTO login_histories (
  user_id,
  username,
  login_time,
  login_status
) VALUES (
  $1, $2, $3, $4
);

-- name: GetLastLogin :one
SELECT login_time
FROM login_histories
WHERE username = $1
ORDER BY login_time DESC
LIMIT 1;