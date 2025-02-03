-- name: CreateNewProgram :one
INSERT INTO programs (name, category, description) 
VALUES ($1, $2, $3) 
RETURNING id, created_at, updated_at;