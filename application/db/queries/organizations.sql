-- name: CreateOrganization :one
INSERT INTO organizations (
  organization_name,
  organization_description
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations
WHERE organization_id = $1 LIMIT 1;

-- name: GetOrganizationByName :one
SELECT * FROM organizations
WHERE organization_name = $1 LIMIT 1;

-- name: UpdateOrganization :exec
UPDATE organizations
SET
  organization_name = $1,
  organization_description = $2,
  updated_at = NOW()
WHERE organization_id = $3
RETURNING *;

-- name: ListOrganizations :many
SELECT * FROM organizations
ORDER BY organization_name
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountOrganizations :one
SELECT COUNT(*) as total FROM organizations;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE organization_id = $1;
