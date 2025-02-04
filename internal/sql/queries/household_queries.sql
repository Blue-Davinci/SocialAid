-- name: CreateNewHousehold :one
INSERT INTO households (program_id, geolocation_id, name) 
VALUES ($1, $2, $3) 
RETURNING id, created_at;

-- name: CreateNewHouseholdHead :one
INSERT INTO household_heads (household_id, name, national_id, phone_number, age)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at;

-- name: CreateNewHouseholdMember :one
INSERT INTO household_members (household_id, name, age, relation)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetHouseholdHeadByHouseholdId :one
SELECT
    id,
    household_id,
    name,
    national_id,
    phone_number,
    age,
    created_at,
    updated_at
FROM household_heads
WHERE household_id = $1;

-- name: GetHouseHoldInformation :many
SELECT 
    h.id AS household_id, 
    h.program_id, 
    p.name AS program_name, 
    h.geolocation_id, 
    g.county, 
    g.sub_county, 
    hh.id AS household_head_id, 
    hh.name AS household_head_name, 
    hh.phone_number,
    COUNT(hm.id) AS household_member_count
FROM households h
JOIN programs p ON h.program_id = p.id
JOIN geolocations g ON h.geolocation_id = g.id
JOIN household_heads hh ON hh.household_id = h.id
LEFT JOIN household_members hm ON hm.household_id = h.id
WHERE h.id = $1
GROUP BY h.id, p.id, g.id, hh.id;
