-- name: CreateNewHousehold :one
INSERT INTO households (program_id, geolocation_id, name) 
VALUES ($1, $2, $3) 
RETURNING id, created_at;



-- name: CreateNewHouseholdHead :one
INSERT INTO household_heads (household_id, name, national_id, phone_number, age)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at;