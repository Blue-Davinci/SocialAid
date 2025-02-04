-- name: CreateNewProgram :one
INSERT INTO programs (name, category, description) 
VALUES ($1, $2, $3) 
RETURNING id, created_at, updated_at;

-- name: GetProgramById :one
SELECT 
    id,
    name,
    category,
    description,
    created_at,
    updated_at
FROM programs
WHERE id = $1;

-- name: UpdateProgramById :one
UPDATE programs
SET 
    name = $2,
    category = $3,
    description = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING updated_at;