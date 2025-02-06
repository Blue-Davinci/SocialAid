-- name: GetForApiKey :one
SELECT 
    id,
    email,
    name,
    api_key,
    created_at,
    updated_at
FROM users
WHERE api_key = $1;