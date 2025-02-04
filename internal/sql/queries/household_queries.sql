-- name: CreateNewHousehold :one
INSERT INTO households (program_id, geolocation_id, name) 
VALUES ($1, $2, $3) 
RETURNING id, created_at;